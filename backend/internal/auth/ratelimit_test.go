package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// The limiter is pure given an injected clock, so this needs no DB or HTTP.
func TestRateLimiterFixedWindow(t *testing.T) {
	rl := newRateLimiter(2, time.Minute)
	t0 := time.Now()

	if ok, _ := rl.allow("1.2.3.4", t0); !ok {
		t.Fatal("1st request denied")
	}
	if ok, _ := rl.allow("1.2.3.4", t0); !ok {
		t.Fatal("2nd request denied")
	}
	ok, retry := rl.allow("1.2.3.4", t0)
	if ok {
		t.Fatal("3rd request allowed past the limit")
	}
	if retry <= 0 || retry > 60 {
		t.Errorf("Retry-After = %d, want 1..60", retry)
	}

	// A different IP has its own budget.
	if ok, _ := rl.allow("5.6.7.8", t0); !ok {
		t.Error("unrelated IP throttled")
	}

	// The window resets once it elapses.
	if ok, _ := rl.allow("1.2.3.4", t0.Add(time.Minute)); !ok {
		t.Error("still throttled after the window elapsed")
	}
}

// Behind the prod proxy chain every request reaches the API from Caddy, so
// RemoteAddr is identical for all clients; the limiter must key on the
// Cloudflare-provided client address for per-client budgets to exist (#133).
// Fails when the key reverts to RemoteAddr-only.
func TestRateLimitIndependentBudgetsPerClientIP(t *testing.T) {
	h := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	send := func(ip string) int {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "172.18.0.9:52048" // Caddy's address, the same for every client
		req.Header.Set("CF-Connecting-IP", ip)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

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
	h := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	send := func(ip string) int {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "172.18.0.9:52048"
		req.Header.Set("CF-Connecting-IP", ip)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

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
