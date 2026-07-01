package chat_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
)

// imageRequestWithHeaders is imageResponse but lets the caller set conditional
// request headers (If-None-Match, Range).
func imageRequestWithHeaders(t *testing.T, h http.Handler, threadID, msgID, token string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/threads/"+threadID+"/messages/"+msgID+"/image", nil)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	h.ServeHTTP(rec, req)
	return rec
}

// TestImageMessageETagAndRange proves #98: a first fetch carries a stable
// ETag and a bounded private max-age with Content-Length, a repeat fetch with
// If-None-Match gets a bodyless 304, a Range request gets partial content,
// and authz is unchanged (a non-participant still gets 403).
func TestImageMessageETagAndRange(t *testing.T) {
	pool := testdb.Pool(t)
	requireStorage(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	strangerTok, _ := registerUser(t, h, pool, "Stranger", true)
	race := seedRace(t, pool, "platform_sale")
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hello")
	msgID := uploadImage(t, h, threadID, buyerTok, "", tinyJPEG(t))

	first := imageRequestWithHeaders(t, h, threadID, msgID, buyerTok, nil)
	if first.Code != http.StatusOK {
		t.Fatalf("first fetch: status = %d", first.Code)
	}
	etag := first.Header().Get("ETag")
	if etag == "" {
		t.Fatal("first fetch: missing ETag")
	}
	if cc := first.Header().Get("Cache-Control"); !strings.Contains(cc, "private") || !strings.Contains(cc, "max-age") {
		t.Errorf("first fetch: Cache-Control = %q, want a bounded private directive", cc)
	}
	if first.Header().Get("Content-Length") == "" {
		t.Error("first fetch: missing Content-Length")
	}

	repeat := imageRequestWithHeaders(t, h, threadID, msgID, buyerTok, map[string]string{"If-None-Match": etag})
	if repeat.Code != http.StatusNotModified {
		t.Fatalf("repeat fetch with If-None-Match: status = %d, want 304", repeat.Code)
	}
	if len(repeat.Body.Bytes()) != 0 {
		t.Errorf("304 response carries a body: %d bytes", len(repeat.Body.Bytes()))
	}

	ranged := imageRequestWithHeaders(t, h, threadID, msgID, buyerTok, map[string]string{"Range": "bytes=0-3"})
	if ranged.Code != http.StatusPartialContent {
		t.Fatalf("range fetch: status = %d, want 206", ranged.Code)
	}
	if len(ranged.Body.Bytes()) != 4 {
		t.Errorf("range fetch: got %d bytes, want 4", len(ranged.Body.Bytes()))
	}

	if rec := imageRequestWithHeaders(t, h, threadID, msgID, strangerTok, map[string]string{"If-None-Match": etag}); rec.Code != http.StatusForbidden {
		t.Errorf("stranger with If-None-Match: status = %d, want 403", rec.Code)
	}
}
