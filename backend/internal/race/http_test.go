package race_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/race"
)

// fixtures get unique slugs and are deleted on cleanup, so tests are safe
// against any database state (seeded or empty).
func seedRace(t *testing.T, pool *pgxpool.Pool, policy, status, country string) sqlcgen.Race {
	t.Helper()
	ctx := context.Background()
	id := ids.New()
	src := "https://example.org/source"
	row, err := sqlcgen.New(pool).CreateRace(ctx, sqlcgen.CreateRaceParams{
		ID: id, Slug: "t-" + id.String(), Name: "Test Race " + id.String()[:8],
		Sport: "running", EventDate: mustDate("2027-06-01"), City: "Testville",
		Country: country, TransferPolicy: policy, PolicySourceUrl: &src,
		OfficialTransferUrl: &src, Status: status,
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, row.ID)
	})
	return row
}

func mustDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return d
}

func handler(pool *pgxpool.Pool) http.Handler {
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool, race.Routes(sqlcgen.New(pool)))
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestGetRaceBySlug(t *testing.T) {
	pool := testdb.Pool(t)
	r := seedRace(t, pool, "official_only", "published", "ES")

	rec := get(t, handler(pool), "/api/v1/races/"+r.Slug)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var body race.Detail
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if body.Slug != r.Slug || body.TransferPolicy != "official_only" {
		t.Errorf("unexpected body: %+v", body)
	}
	if body.OfficialTransferURL == nil {
		t.Error("official_transfer_url missing for official_only race")
	}
	if cc := rec.Header().Get("Cache-Control"); cc == "" {
		t.Error("missing Cache-Control on anonymous catalog response")
	}
}

func TestDraftRaceIs404(t *testing.T) {
	pool := testdb.Pool(t)
	r := seedRace(t, pool, "unknown", "draft", "ES")

	if rec := get(t, handler(pool), "/api/v1/races/"+r.Slug); rec.Code != http.StatusNotFound {
		t.Fatalf("draft race: status = %d, want 404", rec.Code)
	}
	if rec := get(t, handler(pool), "/api/v1/races/no-such-race"); rec.Code != http.StatusNotFound {
		t.Fatalf("missing race: status = %d, want 404", rec.Code)
	}
}

func TestListRacesFiltersAndPaginates(t *testing.T) {
	pool := testdb.Pool(t)
	// Unique country code avoids collisions with seeded data.
	a := seedRace(t, pool, "platform_sale", "published", "ZZ")
	b := seedRace(t, pool, "connect_only", "published", "ZZ")
	_ = seedRace(t, pool, "unknown", "draft", "ZZ") // must not appear

	h := handler(pool)
	rec := get(t, h, "/api/v1/races?country=zz&limit=1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var page1 struct {
		Items      []race.Summary `json:"items"`
		NextCursor *string        `json:"next_cursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page1); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(page1.Items) != 1 || page1.NextCursor == nil {
		t.Fatalf("page1: items = %d, cursor = %v", len(page1.Items), page1.NextCursor)
	}

	rec = get(t, h, fmt.Sprintf("/api/v1/races?country=zz&limit=10&cursor=%s", *page1.NextCursor))
	var page2 struct {
		Items []race.Summary `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page2); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(page2.Items) != 1 {
		t.Fatalf("page2: items = %d, want exactly the second published race", len(page2.Items))
	}
	got := map[uuid.UUID]bool{page1.Items[0].ID: true, page2.Items[0].ID: true}
	if !got[a.ID] || !got[b.ID] {
		t.Errorf("pages = %v, want both published ZZ races exactly once", got)
	}

	// Policy filter.
	rec = get(t, h, "/api/v1/races?country=zz&policy=platform_sale")
	var filtered struct {
		Items []race.Summary `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(filtered.Items) != 1 || filtered.Items[0].ID != a.ID {
		t.Errorf("policy filter returned %d items", len(filtered.Items))
	}
}

func TestListRacesRejectsBadParams(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	for _, path := range []string{
		"/api/v1/races?country=ESP",
		"/api/v1/races?sport=swimming-with-sharks",
		"/api/v1/races?policy=free_for_all",
		"/api/v1/races?date_from=tomorrow",
		"/api/v1/races?limit=0",
		"/api/v1/races?limit=101",
		"/api/v1/races?cursor=garbage",
	} {
		if rec := get(t, h, path); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", path, rec.Code)
		}
	}
}
