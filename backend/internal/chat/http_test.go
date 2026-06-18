package chat_test

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"log/slog"
	"mime/multipart"
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
	"github.com/leonfullxr/bibseller/backend/internal/platform/config"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/platform/storage"
)

type noopMailer struct{}

func (noopMailer) SendVerification(_, _ string) error  { return nil }
func (noopMailer) SendPasswordReset(_, _ string) error { return nil }
func (noopMailer) SendNewMessage(_, _ string) error    { return nil }

// authedHandler mounts chat behind ResolveUser, plus listing (to create a
// listing via the API) and auth (to register and obtain real session cookies).
func authedHandler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	store := newStorage()
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)},
		listing.Routes(q), chat.Routes(pool, noopMailer{}, store, "http://test.local"),
		auth.Routes(pool, noopMailer{}, "http://test.local"))
}

// newStorage builds the object-storage client from env (MinIO defaults). New
// does not connect, so this never fails on a healthy config; image tests gate
// on reachability via requireStorage.
func newStorage() *storage.Client {
	cfg := config.Load()
	store, err := storage.New(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	if err != nil {
		panic(err)
	}
	return store
}

// requireStorage skips the test unless the object store is reachable and the
// bucket exists (mirrors testdb.Pool for Postgres).
func requireStorage(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := newStorage().EnsureBucket(ctx); err != nil {
		t.Skipf("object storage not reachable, skipping image test: %v", err)
	}
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
	HasImage bool   `json:"has_image"`
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

// TestStartThreadReuse proves find-or-create is idempotent: a buyer contacting
// the same listing twice reuses the one thread (the ON CONFLICT DO NOTHING path)
// and both messages land in it.
func TestStartThreadReuse(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)

	first := startThread(t, h, listingID, buyerTok, "hello")
	second := startThread(t, h, listingID, buyerTok, "still here")
	if first != second {
		t.Fatalf("re-contact created a new thread: %s != %s", first, second)
	}
	if msgs := messages(t, h, first, "", buyerTok); len(msgs) != 2 {
		t.Fatalf("thread should hold both messages, got %d", len(msgs))
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
		t.Errorf("unverified write: status = %d, want 403", rec.Code)
	}
	// Reads are gated on verification too, for consistency with writes.
	if rec := doJSON(t, h, http.MethodGet, "/api/v1/threads", "", unverifiedTok); rec.Code != http.StatusForbidden {
		t.Errorf("unverified inbox read: status = %d, want 403", rec.Code)
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

func tinyJPEG(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4)), nil); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}
	return buf.Bytes()
}

// jpegWithEXIF returns a JPEG carrying an APP1 "Exif" segment, so a test can
// prove the server strips it when it re-encodes the upload.
func jpegWithEXIF(t *testing.T) []byte {
	t.Helper()
	base := tinyJPEG(t) // begins with the SOI marker FF D8
	payload := append([]byte("Exif\x00\x00"), make([]byte, 16)...)
	seg := []byte{0xFF, 0xE1, byte((len(payload) + 2) >> 8), byte(len(payload) + 2)}
	seg = append(seg, payload...)
	out := append([]byte{}, base[:2]...) // SOI
	out = append(out, seg...)            // injected APP1/EXIF
	return append(out, base[2:]...)
}

func multipartImage(t *testing.T, filename string, data []byte, caption string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if _, err := fw.Write(data); err != nil {
		t.Fatalf("multipart write: %v", err)
	}
	if caption != "" {
		if err := mw.WriteField("body", caption); err != nil {
			t.Fatalf("multipart field: %v", err)
		}
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("multipart close: %v", err)
	}
	return &buf, mw.FormDataContentType()
}

func postMultipart(t *testing.T, h http.Handler, path, token string, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", contentType)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

func uploadImage(t *testing.T, h http.Handler, threadID, token, caption string, data []byte) string {
	t.Helper()
	body, ct := multipartImage(t, "bib.jpg", data, caption)
	rec := postMultipart(t, h, "/api/v1/threads/"+threadID+"/messages", token, body, ct)
	if rec.Code != http.StatusCreated {
		t.Fatalf("upload image: status = %d, body = %s", rec.Code, rec.Body)
	}
	var m struct {
		ID       string `json:"id"`
		HasImage bool   `json:"has_image"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("upload image: bad JSON: %v", err)
	}
	if !m.HasImage {
		t.Fatal("uploaded message is missing has_image")
	}
	return m.ID
}

func imageResponse(t *testing.T, h http.Handler, threadID, msgID, token string) *httptest.ResponseRecorder {
	t.Helper()
	return doJSON(t, h, http.MethodGet, "/api/v1/threads/"+threadID+"/messages/"+msgID+"/image", "", token)
}

// TestImageMessage proves a participant can attach an image, both participants
// retrieve it (a non-participant cannot), the bytes are re-encoded with the EXIF
// stripped, and the message advertises the image plus its caption (#8 / #38).
func TestImageMessage(t *testing.T) {
	pool := testdb.Pool(t)
	requireStorage(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	strangerTok, _ := registerUser(t, h, pool, "Stranger", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hello")

	msgID := uploadImage(t, h, threadID, buyerTok, "my bib", jpegWithEXIF(t))

	rec := imageResponse(t, h, threadID, msgID, sellerTok)
	if rec.Code != http.StatusOK {
		t.Fatalf("seller fetch image: status = %d", rec.Code)
	}
	got := rec.Body.Bytes()
	if bytes.Contains(got, []byte("Exif")) {
		t.Error("retrieved image still carries an EXIF segment")
	}
	if _, _, err := image.Decode(bytes.NewReader(got)); err != nil {
		t.Errorf("retrieved image does not decode: %v", err)
	}
	if rec := imageResponse(t, h, threadID, msgID, buyerTok); rec.Code != http.StatusOK {
		t.Errorf("buyer fetch image: status = %d, want 200", rec.Code)
	}
	if rec := imageResponse(t, h, threadID, msgID, strangerTok); rec.Code != http.StatusForbidden {
		t.Errorf("stranger fetch image: status = %d, want 403", rec.Code)
	}

	for _, m := range messages(t, h, threadID, "", buyerTok) {
		if m.ID == msgID && (!m.HasImage || m.Body != "my bib") {
			t.Errorf("image message in list: %+v", m)
		}
	}

	// If the object vanishes from storage, the download is a 404, not a 500.
	var key string
	if err := pool.QueryRow(context.Background(),
		`SELECT image_key FROM messages WHERE id = $1`, uuid.MustParse(msgID)).Scan(&key); err != nil {
		t.Fatalf("read image_key: %v", err)
	}
	if err := newStorage().Delete(context.Background(), key); err != nil {
		t.Fatalf("delete object: %v", err)
	}
	if rec := imageResponse(t, h, threadID, msgID, buyerTok); rec.Code != http.StatusNotFound {
		t.Errorf("missing object fetch: status = %d, want 404", rec.Code)
	}
}

func TestImageMessageRejectsBadUploads(t *testing.T) {
	pool := testdb.Pool(t)
	requireStorage(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hello")

	// Non-image bytes are rejected once decoding fails.
	body, ct := multipartImage(t, "notes.jpg", []byte("definitely not an image"), "")
	if rec := postMultipart(t, h, "/api/v1/threads/"+threadID+"/messages", buyerTok, body, ct); rec.Code != http.StatusBadRequest {
		t.Errorf("non-image upload: status = %d, want 400", rec.Code)
	}

	// An oversized upload is rejected before decoding.
	big, ctBig := multipartImage(t, "big.jpg", make([]byte, 7<<20), "")
	if rec := postMultipart(t, h, "/api/v1/threads/"+threadID+"/messages", buyerTok, big, ctBig); rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized upload: status = %d, want 413", rec.Code)
	}
}

// TestImageMessageWithoutCaption covers an image-only message (no caption): the
// API returns body:null with has_image, and it round-trips through the list.
func TestImageMessageWithoutCaption(t *testing.T) {
	pool := testdb.Pool(t)
	requireStorage(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hello")

	body, ct := multipartImage(t, "bib.jpg", tinyJPEG(t), "")
	rec := postMultipart(t, h, "/api/v1/threads/"+threadID+"/messages", buyerTok, body, ct)
	if rec.Code != http.StatusCreated {
		t.Fatalf("image-only upload: status = %d, body = %s", rec.Code, rec.Body)
	}
	var created struct {
		ID       string  `json:"id"`
		Body     *string `json:"body"`
		HasImage bool    `json:"has_image"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if created.Body != nil || !created.HasImage {
		t.Errorf("image-only message: body = %v, has_image = %v; want nil, true", created.Body, created.HasImage)
	}

	// It also round-trips through the list (a null body unmarshals to "" here).
	found := false
	for _, m := range messages(t, h, threadID, "", buyerTok) {
		if m.ID == created.ID {
			found = m.HasImage && m.Body == ""
		}
	}
	if !found {
		t.Error("image-only message not reflected with has_image and empty body in the list")
	}
}

// TestBlockedCannotChat proves a block stops new threads and messages both ways.
func TestBlockedCannotChat(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, sellerID := registerUser(t, h, pool, "Seller", true)
	buyerTok, buyerID := registerUser(t, h, pool, "Buyer", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hello")

	if _, err := pool.Exec(context.Background(),
		`INSERT INTO blocks (blocker_id, blocked_id) VALUES ($1, $2)`, buyerID, sellerID); err != nil {
		t.Fatalf("insert block: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM blocks WHERE blocker_id = $1 OR blocked_id = $1`, buyerID)
	})

	// Neither side can send in the existing thread (the block is symmetric).
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/threads/"+threadID+"/messages", `{"body":"hi"}`, buyerTok); rec.Code != http.StatusForbidden {
		t.Errorf("blocker send: status = %d, want 403", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/threads/"+threadID+"/messages", `{"body":"hi"}`, sellerTok); rec.Code != http.StatusForbidden {
		t.Errorf("blocked send: status = %d, want 403", rec.Code)
	}

	// And no new thread can be started with the same seller.
	listing2 := createListing(t, h, race.ID, sellerTok)
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/listings/"+listing2+"/threads", `{"body":"hi"}`, buyerTok); rec.Code != http.StatusForbidden {
		t.Errorf("blocked start thread: status = %d, want 403", rec.Code)
	}
}
