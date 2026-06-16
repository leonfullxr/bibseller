package auth

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

// Brute-force / spam budget per client IP on the argon2id-heavy endpoints.
// Generous enough for a shared NAT, tight enough to make online password
// guessing and bulk registration impractical.
const (
	rateLimitMax    = 10
	rateLimitWindow = time.Minute
)

// rateLimiter is a per-key fixed-window counter. In-process, no external store
// (docs/ARCHITECTURE.md → no new infra): v1 runs a single API instance, and a
// fixed window resets without per-request bookkeeping. The trigger to revisit
// (Redis-backed) is horizontal scaling of the API.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*window
	limit  int
	window time.Duration
}

type window struct {
	count int
	start time.Time
}

func newRateLimiter(limit int, w time.Duration) *rateLimiter {
	return &rateLimiter{hits: map[string]*window{}, limit: limit, window: w}
}

// allow reports whether key may proceed at time now, and (when it may not) the
// whole seconds until its window resets — the Retry-After hint.
func (rl *rateLimiter) allow(key string, now time.Time) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	w := rl.hits[key]
	if w == nil || now.Sub(w.start) >= rl.window {
		rl.hits[key] = &window{count: 1, start: now}
		return true, 0
	}
	if w.count >= rl.limit {
		// Whole seconds (rounded up, ≥1) until this window resets.
		return false, int(math.Ceil((rl.window - now.Sub(w.start)).Seconds()))
	}
	w.count++
	return true, 0
}

// sweep evicts expired windows so the map can't grow without bound. Runs for
// the life of the process.
func (rl *rateLimiter) sweep(every time.Duration) {
	for range time.Tick(every) {
		rl.mu.Lock()
		now := time.Now()
		for k, w := range rl.hits {
			if now.Sub(w.start) >= rl.window {
				delete(rl.hits, k)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware throttles the brute-forceable auth endpoints per client
// IP. Mount it on the v1 sub-mux, where paths are v1-relative ("/auth/login");
// every other route (the public catalog, /auth/me, /auth/logout) passes
// straight through.
func RateLimit() func(http.Handler) http.Handler {
	rl := newRateLimiter(rateLimitMax, rateLimitWindow)
	go rl.sweep(rateLimitWindow)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/auth/login", "/auth/register", "/auth/verify/resend":
				if ok, retry := rl.allow(clientIP(r), time.Now()); !ok {
					w.Header().Set("Retry-After", strconv.Itoa(retry))
					httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "too many requests, slow down")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
