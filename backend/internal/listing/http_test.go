package listing_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/listing"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

type fixture struct {
	race    sqlcgen.Race
	listing sqlcgen.Listing
}

func seed(t *testing.T, pool *pgxpool.Pool, raceStatus string) fixture {
	t.Helper()
	ctx := context.Background()
	q := sqlcgen.New(pool)

	userID := ids.New()
	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, display_name)
		 VALUES ($1, $2, 'x', 'Fixture Seller')`, userID, userID.String()+"@test.local")
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	raceID := ids.New()
	src := "https://example.org/source"
	race, err := q.CreateRace(ctx, sqlcgen.CreateRaceParams{
		ID: raceID, Slug: "t-" + raceID.String(), Name: "Fixture Race",
		Sport: "running", EventDate: time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC),
		City: "Testville", Country: "ZZ", TransferPolicy: "platform_sale",
		PolicySourceUrl: &src, Status: raceStatus,
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}

	price, orig := int32(5000), int32(6000)
	l, err := q.CreateListing(ctx, sqlcgen.CreateListingParams{
		ID: ids.New(), RaceID: race.ID, SellerID: userID,
		PriceCents: &price, Currency: "EUR", OriginalPriceCents: &orig,
		ExpiresAt: race.EventDate,
	})
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id = $1`, l.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})
	return fixture{race: race, listing: l}
}

func handler(pool *pgxpool.Pool) http.Handler {
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool, listing.Routes(sqlcgen.New(pool)))
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestListingsByRace(t *testing.T) {
	pool := testdb.Pool(t)
	fx := seed(t, pool, "published")

	rec := get(t, handler(pool), "/api/v1/races/"+fx.race.Slug+"/listings")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Items []listing.Summary `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(body.Items))
	}
	got := body.Items[0]
	if got.ID != fx.listing.ID || got.SellerName != "Fixture Seller" || got.PriceCents == nil || *got.PriceCents != 5000 {
		t.Errorf("unexpected item: %+v", got)
	}
}

func TestListingDetailCarriesRaceContext(t *testing.T) {
	pool := testdb.Pool(t)
	fx := seed(t, pool, "published")

	rec := get(t, handler(pool), "/api/v1/listings/"+fx.listing.ID.String())
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var body listing.Detail
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if body.Race.Slug != fx.race.Slug || body.Race.TransferPolicy != "platform_sale" {
		t.Errorf("race context: %+v", body.Race)
	}
}

func TestListingOnDraftRaceIsHidden(t *testing.T) {
	pool := testdb.Pool(t)
	fx := seed(t, pool, "draft")

	if rec := get(t, handler(pool), "/api/v1/listings/"+fx.listing.ID.String()); rec.Code != http.StatusNotFound {
		t.Errorf("detail on draft race: status = %d, want 404", rec.Code)
	}
	if rec := get(t, handler(pool), "/api/v1/races/"+fx.race.Slug+"/listings"); rec.Code != http.StatusNotFound {
		t.Errorf("list on draft race: status = %d, want 404", rec.Code)
	}
}

func TestListingNotFound(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	if rec := get(t, h, "/api/v1/listings/not-a-uuid"); rec.Code != http.StatusNotFound {
		t.Errorf("bad id: status = %d, want 404", rec.Code)
	}
	if rec := get(t, h, "/api/v1/listings/"+ids.New().String()); rec.Code != http.StatusNotFound {
		t.Errorf("unknown id: status = %d, want 404", rec.Code)
	}
}
