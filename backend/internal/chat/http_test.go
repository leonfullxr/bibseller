package chat_test

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
	"github.com/leonfullxr/bibseller/backend/internal/chat"
	"github.com/leonfullxr/bibseller/backend/internal/listing"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

type noopMailer struct{}

func (noopMailer) SendVerification(_, _ string) error  { return nil }
func (noopMailer) SendPasswordReset(_, _ string) error { return nil }
func (noopMailer) SendNewMessage(_, _ string) error    { return nil }

// authedHandler mounts chat behind ResolveUser, plus listing (to create a
// listing via the API) and auth (to register and obtain real session cookies).
func authedHandler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)},
		listing.Routes(q), chat.Routes(pool, noopMailer{}, "http://test.local"),
		auth.Routes(pool, noopMailer{}, "http://test.local"))
}

// registerUser creates a throwaway verified account and returns its session
// token and id. Chat (like listing) requires a verified email.
func registerUser(t *testing.T, h http.Handler, pool *pgxpool.Pool, name string, verified bool) (token string, id uuid.UUID) {
	t.Helper()
	email := "u-" + ids.New().String() + "@test.local"
	body := `{"email":"` + email + `","password":"correct horse battery staple","display_name":"` + name + `"}`
	rec := doJSON(t, h, http.MethodPost, "/api/v1/auth/register", body, "")
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
			t.Fatalf("verify user: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, resp.User.ID)
	})
	return resp.Token, resp.User.ID
}

// seedRace inserts a published race with the evidence its policy requires, and
// cleans up every chat/listing row hanging off it (FK order matters).
func seedRace(t *testing.T, pool *pgxpool.Pool, policy string) sqlcgen.Race {
	t.Helper()
	id := ids.New()
	params := sqlcgen.CreateRaceParams{
		ID: id, Slug: "tr-" + id.String(), Name: "Test Race", Sport: "running",
		EventDate: future(), City: "Testville", Country: "ZX",
		TransferPolicy: policy, Status: "published",
	}
	src := "https://example.org/policy"
	switch policy {
	case "platform_sale":
		params.PolicySourceUrl = &src
	case "official_only":
		params.OfficialTransferUrl = &src
	}
	race, err := sqlcgen.New(pool).CreateRace(context.Background(), params)
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM messages WHERE thread_id IN
			(SELECT t.id FROM chat_threads t JOIN listings l ON l.id = t.listing_id WHERE l.race_id = $1)`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM chat_threads WHERE listing_id IN (SELECT id FROM listings WHERE race_id = $1)`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM policy_acks WHERE race_id = $1`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE race_id = $1`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
	})
	return race
}

// createListing posts a listing as the (verified) seller and returns its id.
func createListing(t *testing.T, h http.Handler, raceID uuid.UUID, token string) string {
	t.Helper()
	body := `{"race_id":"` + raceID.String() + `","price_cents":5000,"original_price_cents":6000}`
	rec := doJSON(t, h, http.MethodPost, "/api/v1/listings", body, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create listing: status = %d, body = %s", rec.Code, rec.Body)
	}
	var got listing.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("create listing: bad JSON: %v", err)
	}
	return got.ID.String()
}

func future() time.Time { return time.Now().UTC().AddDate(0, 1, 0).Truncate(24 * time.Hour) }

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

// startThread posts a first message and returns the created thread id.
func startThread(t *testing.T, h http.Handler, listingID, token, body string) string {
	t.Helper()
	rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads",
		`{"body":"`+body+`"}`, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("start thread: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		ThreadID string `json:"thread_id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("start thread: bad JSON: %v", err)
	}
	return resp.ThreadID
}

type inboxThread struct {
	ID          string `json:"id"`
	Role        string `json:"role"`
	OtherParty  string `json:"other_party"`
	UnreadCount int    `json:"unread_count"`
}

func inbox(t *testing.T, h http.Handler, token string) []inboxThread {
	t.Helper()
	rec := doJSON(t, h, http.MethodGet, "/api/v1/threads", "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("inbox: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Items []inboxThread `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("inbox: bad JSON: %v", err)
	}
	return resp.Items
}

type message struct {
	ID       string `json:"id"`
	SenderID string `json:"sender_id"`
	Body     string `json:"body"`
}

func messages(t *testing.T, h http.Handler, threadID, since, token string) []message {
	t.Helper()
	path := "/api/v1/threads/" + threadID + "/messages"
	if since != "" {
		path += "?since=" + since
	}
	rec := doJSON(t, h, http.MethodGet, path, "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("messages: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Items []message `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("messages: bad JSON: %v", err)
	}
	return resp.Items
}

// TestStartThreadAckGate proves the #8 acceptance: in a gated mode the first
// message is impossible until the ack is stored, then succeeds.
func TestStartThreadAckGate(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, buyerID := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "connect_only")
	listingID := createListing(t, h, race.ID, sellerTok)

	// No ack yet: blocked.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads",
		`{"body":"hi"}`, buyerTok); rec.Code != http.StatusForbidden {
		t.Fatalf("pre-ack start: status = %d, want 403", rec.Code)
	}

	// Acknowledge, then it works.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/races/"+race.Slug+"/ack", "", buyerTok); rec.Code != http.StatusNoContent {
		t.Fatalf("ack: status = %d, body = %s", rec.Code, rec.Body)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads",
		`{"body":"hi"}`, buyerTok); rec.Code != http.StatusCreated {
		t.Fatalf("post-ack start: status = %d, body = %s", rec.Code, rec.Body)
	}

	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM policy_acks WHERE user_id = $1 AND race_id = $2`,
		buyerID, race.ID).Scan(&n); err != nil || n != 1 {
		t.Fatalf("ack row: n = %d, err = %v", n, err)
	}
}

// TestConverseAndInbox proves two participants exchange messages and that unread
// counts track who has read what (#8 top-line acceptance), on a non-gated race.
func TestConverseAndInbox(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, sellerID := registerUser(t, h, pool, "Seller", true)
	buyerTok, buyerID := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "platform_sale") // non-gated: no ack needed
	listingID := createListing(t, h, race.ID, sellerTok)

	threadID := startThread(t, h, listingID, buyerTok, "hello")

	// Seller's inbox: one thread from the buyer, one unread.
	sInbox := inbox(t, h, sellerTok)
	if len(sInbox) != 1 || sInbox[0].UnreadCount != 1 || sInbox[0].Role != "seller" || sInbox[0].OtherParty != "Buyer" {
		t.Fatalf("seller inbox after first message: %+v", sInbox)
	}
	// Buyer's own message is not unread for the buyer.
	if bInbox := inbox(t, h, buyerTok); len(bInbox) != 1 || bInbox[0].UnreadCount != 0 || bInbox[0].Role != "buyer" {
		t.Fatalf("buyer inbox after own message: %+v", bInbox)
	}

	// Seller reads the thread, clearing their unread.
	msgs := messages(t, h, threadID, "", sellerTok)
	if len(msgs) != 1 || msgs[0].Body != "hello" || msgs[0].SenderID != buyerID.String() {
		t.Fatalf("seller sees: %+v", msgs)
	}
	if s := inbox(t, h, sellerTok); s[0].UnreadCount != 0 {
		t.Fatalf("seller unread after read: %d, want 0", s[0].UnreadCount)
	}

	// Seller replies; now the buyer has one unread.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/threads/"+threadID+"/messages",
		`{"body":"hi back"}`, sellerTok); rec.Code != http.StatusCreated {
		t.Fatalf("seller reply: status = %d, body = %s", rec.Code, rec.Body)
	}
	if b := inbox(t, h, buyerTok); b[0].UnreadCount != 1 {
		t.Fatalf("buyer unread after reply: %d, want 1", b[0].UnreadCount)
	}

	// Buyer polls with the first message as cursor: only the reply comes back.
	newer := messages(t, h, threadID, msgs[0].ID, buyerTok)
	if len(newer) != 1 || newer[0].Body != "hi back" || newer[0].SenderID != sellerID.String() {
		t.Fatalf("buyer poll since first: %+v", newer)
	}
	if b := inbox(t, h, buyerTok); b[0].UnreadCount != 0 {
		t.Fatalf("buyer unread after read: %d, want 0", b[0].UnreadCount)
	}
}

// TestThreadAuthz proves only participants can read/write a thread and a seller
// cannot open a thread on their own listing.
func TestThreadAuthz(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	strangerTok, _ := registerUser(t, h, pool, "Stranger", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)

	// Seller cannot start a thread on their own listing.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads",
		`{"body":"mine"}`, sellerTok); rec.Code != http.StatusForbidden {
		t.Errorf("seller self-thread: status = %d, want 403", rec.Code)
	}

	threadID := startThread(t, h, listingID, buyerTok, "hello")

	// A non-participant can neither read nor write the thread.
	if rec := doJSON(t, h, http.MethodGet, "/api/v1/threads/"+threadID+"/messages", "", strangerTok); rec.Code != http.StatusForbidden {
		t.Errorf("stranger read: status = %d, want 403", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/threads/"+threadID+"/messages",
		`{"body":"intrude"}`, strangerTok); rec.Code != http.StatusForbidden {
		t.Errorf("stranger write: status = %d, want 403", rec.Code)
	}
}

// TestStartThreadGates covers the auth, verification, and listing-state guards.
func TestStartThreadGates(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	unverifiedTok, _ := registerUser(t, h, pool, "Unverified", false)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads", `{"body":"x"}`, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no session: status = %d, want 401", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads", `{"body":"x"}`, unverifiedTok); rec.Code != http.StatusForbidden {
		t.Errorf("unverified: status = %d, want 403", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+ids.New().String()+"/threads", `{"body":"x"}`, unverifiedTok); rec.Code != http.StatusForbidden {
		// Unverified is rejected before the listing is even looked up.
		t.Errorf("unverified unknown listing: status = %d, want 403", rec.Code)
	}

	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+ids.New().String()+"/threads", `{"body":"x"}`, buyerTok); rec.Code != http.StatusNotFound {
		t.Errorf("unknown listing: status = %d, want 404", rec.Code)
	}
	// Empty body is rejected.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads", `{"body":"   "}`, buyerTok); rec.Code != http.StatusBadRequest {
		t.Errorf("blank body: status = %d, want 400", rec.Code)
	}

	// A cancelled listing cannot be contacted.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/cancel", "", sellerTok); rec.Code != http.StatusOK {
		t.Fatalf("cancel for setup: status = %d, body = %s", rec.Code, rec.Body)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listingID+"/threads", `{"body":"x"}`, buyerTok); rec.Code != http.StatusConflict {
		t.Errorf("cancelled listing: status = %d, want 409", rec.Code)
	}
}
