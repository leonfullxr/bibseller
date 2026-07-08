package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

// trustProxyHeader enables CF-Connecting-IP trust for one test, restoring the
// default-off state afterwards (#182). Callers must not use t.Parallel: the
// switch is package-global.
func trustProxyHeader(t *testing.T) {
	t.Helper()
	httpx.SetTrustProxyHeader(true)
	t.Cleanup(func() { httpx.SetTrustProxyHeader(false) })
}

// newLimitedSender returns a send func hitting a RateLimit()-wrapped 204
// handler as the given client IP, from Caddy's fixed RemoteAddr (#189).
func newLimitedSender(t *testing.T) func(ip string) int {
	t.Helper()
	h := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	return func(ip string) int {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "172.18.0.9:52048" // Caddy's address, the same for every client
		req.Header.Set("CF-Connecting-IP", ip)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}
}

// Behind the prod proxy chain every request reaches the API from Caddy, so
// RemoteAddr is identical for all clients; with header trust enabled the
// limiter must key on the Cloudflare-provided client address for per-client
// budgets to exist (#133). Fails when the key reverts to RemoteAddr-only.
func TestRateLimitIndependentBudgetsPerClientIP(t *testing.T) {
	trustProxyHeader(t)
	send := newLimitedSender(t)

	for i := range rateLimitMax {
		if code := send("203.0.113.7"); code != http.StatusNoContent {
			t.Fatalf("request %d from first client: status = %d", i+1, code)
		}
	}
	if code := send("203.0.113.7"); code != http.StatusTooManyRequests {
		t.Fatalf("over-limit request from first client: status = %d, want 429", code)
	}
	if code := send("198.51.100.4"); code != http.StatusNoContent {
		t.Fatalf("second client shares the first client's budget: status = %d, want 204", code)
	}
}

// An IPv6 client owns its whole routed /64, so per-address budgets would be
// rotated around for free; addresses within one /64 must share a budget
// (#133 follow-up). Fails when the limiter keys on the full address.
func TestRateLimitBucketsIPv6ByPrefix(t *testing.T) {
	trustProxyHeader(t)
	send := newLimitedSender(t)

	// Exhaust the budget rotating the low 64 bits within one /64.
	for i := range rateLimitMax {
		if code := send(fmt.Sprintf("2001:db8:2:3::%x", i+1)); code != http.StatusNoContent {
			t.Fatalf("request %d from the /64: status = %d", i+1, code)
		}
	}
	if code := send("2001:db8:2:3::dead"); code != http.StatusTooManyRequests {
		t.Fatalf("rotated address inside the exhausted /64: status = %d, want 429", code)
	}
	// A different /64 is a different client.
	if code := send("2001:db8:2:4::1"); code != http.StatusNoContent {
		t.Fatalf("neighboring /64 shares the budget: status = %d, want 204", code)
	}
}

// Without the TRUST_PROXY_HEADER opt-in, CF-Connecting-IP is
// attacker-controlled, so the limiter must key on RemoteAddr: rotating the
// header must not mint fresh budgets (#182). Fails when the header is trusted
// unconditionally.
func TestRateLimitIgnoresHeaderWithoutTrust(t *testing.T) {
	send := newLimitedSender(t) // trust deliberately left off

	// Rotate the spoofed header every request; all share RemoteAddr's budget.
	for i := range rateLimitMax {
		if code := send(fmt.Sprintf("203.0.113.%d", i+1)); code != http.StatusNoContent {
			t.Fatalf("request %d: status = %d", i+1, code)
		}
	}
	if code := send("198.51.100.99"); code != http.StatusTooManyRequests {
		t.Fatalf("spoofed header minted a fresh budget: status = %d, want 429", code)
	}
}
