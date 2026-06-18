package listing_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/listing"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

type noopMailer struct{}

func (noopMailer) SendVerification(_, _ string) error  { return nil }
func (noopMailer) SendPasswordReset(_, _ string) error { return nil }

// authedHandler mounts the listing routes behind ResolveUser, plus auth.Routes
// so a test can register and obtain a real session cookie.
func authedHandler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)},
		listing.Routes(q), auth.Routes(pool, noopMailer{}, "http://test.local"))
}

// registerSeller creates a throwaway account and (optionally) marks it
// email-verified - the gate the create-listing endpoint enforces.
func registerSeller(t *testing.T, h http.Handler, pool *pgxpool.Pool, verified bool) (token string, id uuid.UUID) {
	t.Helper()
	email := "s-" + ids.New().String() + "@test.local"
	body := `{"email":"` + email + `","password":"correct horse battery staple","display_name":"Seller One"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Token string       `json:"token"`
		User  auth.Account `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("register: bad JSON: %v", err)
	}
	if verified {
		if err := sqlcgen.New(pool).MarkEmailVerified(context.Background(), resp.User.ID); err != nil {
			t.Fatalf("verify seller: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, resp.User.ID)
	})
	return resp.Token, resp.User.ID
}

// seedRace inserts a race for listing tests and cleans up its listings + itself.
func seedRace(t *testing.T, pool *pgxpool.Pool, status string, eventDate time.Time) sqlcgen.Race {
	t.Helper()
	src := "https://example.org/source"
	id := ids.New()
	race, err := sqlcgen.New(pool).CreateRace(context.Background(), sqlcgen.CreateRaceParams{
		ID: id, Slug: "tr-" + id.String(), Name: "Mut Race", Sport: "running",
		EventDate: eventDate, City: "Testville", Country: "ZX",
		TransferPolicy: "platform_sale", PolicySourceUrl: &src, Status: status,
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE race_id = $1`, race.ID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM races WHERE id = $1`, race.ID)
	})
	return race
}

func doJSON(t *testing.T, h http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

// createListing posts a valid listing and returns its id.
func createListing(t *testing.T, h http.Handler, raceID uuid.UUID, token string) string {
	t.Helper()
	body := `{"race_id":"` + raceID.String() + `","price_cents":5000,"original_price_cents":6000}`
	rec := doJSON(t, h, http.MethodPost, "/api/v1/listings", body, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, body = %s", rec.Code, rec.Body)
	}
	var got listing.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("create: bad JSON: %v", err)
	}
	return got.ID.String()
}

func future() time.Time { return time.Now().UTC().AddDate(0, 1, 0).Truncate(24 * time.Hour) }

// catalogHas reports whether the public per-race listing endpoint includes id.
func catalogHas(t *testing.T, h http.Handler, slug, id string) bool {
	t.Helper()
	rec := get(t, h, "/api/v1/races/"+slug+"/listings")
	if rec.Code != http.StatusOK {
		t.Fatalf("listByRace: status = %d, body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("listByRace: bad JSON: %v", err)
	}
	for _, it := range body.Items {
		if it.ID == id {
			return true
		}
	}
	return false
}

func TestCreateListing(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	token, sellerID := registerSeller(t, h, pool, true)
	race := seedRace(t, pool, "published", future())

	body := `{"race_id":"` + race.ID.String() + `","price_cents":5000,"original_price_cents":6000,"description":"bib for sale"}`
	rec := doJSON(t, h, http.MethodPost, "/api/v1/listings", body, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, body = %s", rec.Code, rec.Body)
	}
	var got listing.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if got.Status != "active" || got.PriceCents == nil || *got.PriceCents != 5000 || got.SellerName != "Seller One" {
		t.Errorf("unexpected listing: %+v", got)
	}
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM listings WHERE id = $1 AND seller_id = $2 AND status = 'active'`,
		got.ID, sellerID).Scan(&n); err != nil || n != 1 {
		t.Fatalf("stored listing: n = %d, err = %v", n, err)
	}
}

func TestCreateListingGates(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	token, _ := registerSeller(t, h, pool, true)
	unverified, _ := registerSeller(t, h, pool, false)
	pub := seedRace(t, pool, "published", future())

	body := `{"race_id":"` + pub.ID.String() + `"}`
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings", body, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no session: status = %d, want 401", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings", body, unverified); rec.Code != http.StatusForbidden {
		t.Errorf("unverified: status = %d, want 403", rec.Code)
	}

	draft := seedRace(t, pool, "draft", future())
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings",
		`{"race_id":"`+draft.ID.String()+`"}`, token); rec.Code != http.StatusNotFound {
		t.Errorf("draft race: status = %d, want 404", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings",
		`{"race_id":"`+ids.New().String()+`"}`, token); rec.Code != http.StatusNotFound {
		t.Errorf("unknown race: status = %d, want 404", rec.Code)
	}

	past := seedRace(t, pool, "published", time.Now().UTC().AddDate(0, 0, -2).Truncate(24*time.Hour))
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings",
		`{"race_id":"`+past.ID.String()+`"}`, token); rec.Code != http.StatusBadRequest {
		t.Errorf("past race: status = %d, want 400", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings",
		`{"race_id":"`+pub.ID.String()+`","price_cents":7000,"original_price_cents":6000}`, token); rec.Code != http.StatusBadRequest {
		t.Errorf("over-cap price: status = %d, want 400", rec.Code)
	}
}

func TestUpdateListing(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	token, _ := registerSeller(t, h, pool, true)
	other, _ := registerSeller(t, h, pool, true)
	race := seedRace(t, pool, "published", future())
	id := createListing(t, h, race.ID, token)

	if rec := doJSON(t, h, http.MethodPatch, "/api/v1/listings/"+id,
		`{"price_cents":4000}`, other); rec.Code != http.StatusForbidden {
		t.Errorf("non-owner edit: status = %d, want 403", rec.Code)
	}

	rec := doJSON(t, h, http.MethodPatch, "/api/v1/listings/"+id,
		`{"price_cents":4000,"original_price_cents":6000,"description":"reduced"}`, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("edit: status = %d, body = %s", rec.Code, rec.Body)
	}
	var got listing.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if got.PriceCents == nil || *got.PriceCents != 4000 {
		t.Errorf("price not updated: %+v", got)
	}

	if rec := doJSON(t, h, http.MethodPatch, "/api/v1/listings/"+id,
		`{"price_cents":9000,"original_price_cents":6000}`, token); rec.Code != http.StatusBadRequest {
		t.Errorf("over-cap edit: status = %d, want 400", rec.Code)
	}

	// A non-active listing cannot be edited (409). Cancel, then try to edit.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+id+"/cancel", "", token); rec.Code != http.StatusOK {
		t.Fatalf("cancel for setup: status = %d, body = %s", rec.Code, rec.Body)
	}
	if rec := doJSON(t, h, http.MethodPatch, "/api/v1/listings/"+id,
		`{"price_cents":3000}`, token); rec.Code != http.StatusConflict {
		t.Errorf("edit cancelled listing: status = %d, want 409", rec.Code)
	}
}

func TestCancelListing(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	token, _ := registerSeller(t, h, pool, true)
	other, _ := registerSeller(t, h, pool, true)
	race := seedRace(t, pool, "published", future())
	id := createListing(t, h, race.ID, token)

	// An active listing shows in the public catalog.
	if !catalogHas(t, h, race.Slug, id) {
		t.Error("listing missing from catalog while active")
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+id+"/cancel", "", other); rec.Code != http.StatusForbidden {
		t.Errorf("non-owner cancel: status = %d, want 403", rec.Code)
	}

	rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+id+"/cancel", "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("cancel: status = %d, body = %s", rec.Code, rec.Body)
	}
	var got listing.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if got.Status != "cancelled" {
		t.Errorf("status = %q, want cancelled", got.Status)
	}

	// And it disappears from the public catalog.
	if catalogHas(t, h, race.Slug, id) {
		t.Error("cancelled listing still in catalog")
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+id+"/cancel", "", token); rec.Code != http.StatusConflict {
		t.Errorf("re-cancel: status = %d, want 409", rec.Code)
	}
}

func TestListMine(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	token, _ := registerSeller(t, h, pool, true)
	race := seedRace(t, pool, "published", future())
	id := createListing(t, h, race.ID, token)

	rec := doJSON(t, h, http.MethodGet, "/api/v1/me/listings", "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("list mine: status = %d, body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Items []struct {
			ID       string `json:"id"`
			RaceName string `json:"race_name"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(body.Items) != 1 || body.Items[0].ID != id || body.Items[0].RaceName != "Mut Race" {
		t.Errorf("mine: %+v", body.Items)
	}

	if rec := doJSON(t, h, http.MethodGet, "/api/v1/me/listings", "", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("unauth list mine: status = %d, want 401", rec.Code)
	}
}

// TestListingDetailIsOwn proves the listing detail's is_own_listing is computed
// per viewer from the session: true only for the seller, false for others and
// anonymous (so no stable seller id is exposed publicly).
func TestListingDetailIsOwn(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerSeller(t, h, pool, true)
	otherTok, _ := registerSeller(t, h, pool, true)
	race := seedRace(t, pool, "published", future())
	id := createListing(t, h, race.ID, sellerTok)

	isOwn := func(token string) bool {
		rec := doJSON(t, h, http.MethodGet, "/api/v1/listings/"+id, "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("get listing: status = %d, body = %s", rec.Code, rec.Body)
		}
		var d listing.Detail
		if err := json.Unmarshal(rec.Body.Bytes(), &d); err != nil {
			t.Fatalf("bad JSON: %v", err)
		}
		return d.IsOwnListing
	}

	if !isOwn(sellerTok) {
		t.Error("seller viewing own listing: is_own_listing = false, want true")
	}
	if isOwn(otherTok) {
		t.Error("other user: is_own_listing = true, want false")
	}
	if isOwn("") {
		t.Error("anonymous: is_own_listing = true, want false")
	}
}
