package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ratelimit"
)

// Brute-force / spam budget per client IP on the argon2id-heavy endpoints.
// Generous enough for a shared NAT, tight enough to make online password
// guessing and bulk registration impractical.
const (
	rateLimitMax    = 10
	rateLimitWindow = time.Minute
)

// RateLimit middleware throttles the brute-forceable auth endpoints per client
// IP. Mount it on the v1 sub-mux, where paths are v1-relative ("/auth/login");
// every other route (the public catalog, /auth/me, /auth/logout) passes
// straight through.
func RateLimit() func(http.Handler) http.Handler {
	rl := ratelimit.New(rateLimitMax, rateLimitWindow)
	go rl.Sweep(rateLimitWindow)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/auth/login", "/auth/register", "/auth/verify/resend",
				"/auth/password/reset/request", "/auth/password/reset":
				// /auth/password/reset is included because it runs an argon2 hash
				// before the token is validated - without a cap, spamming bogus
				// tokens with valid-looking passwords is an unauthenticated CPU DoS.
				if ok, retry := rl.Allow(httpx.ClientIPKey(r), time.Now()); !ok {
					w.Header().Set("Retry-After", strconv.Itoa(retry))
					httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "too many requests, slow down")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
