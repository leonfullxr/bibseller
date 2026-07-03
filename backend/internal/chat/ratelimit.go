package chat

import (
	"math"
	"sync"
	"time"
)

// Message-send budgets. Polling (GET) is unmetered; only writes are capped.
// Per-account is the primary anti-spam control; the per-IP cap is deliberately
// generous so it stops a single source flooding writes across many accounts
// without tripping a shared NAT of legitimate users.
const (
	msgRateMax    = 30
	ipRateMax     = 120
	msgRateWindow = time.Minute
)

// rateLimiter is a per-key fixed-window counter, in-process (no external store,
// per docs/ARCHITECTURE.md). Keyed by sender account, so it caps a user across
// every IP they post from.
//
// ponytail: duplicated from internal/auth's limiter; promote to a shared
// internal/platform/ratelimit on the third consumer (M5.4 reports) - two copies
// is below the extract-it threshold.
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
// whole seconds until its window resets - the Retry-After hint.
func (rl *rateLimiter) allow(key string, now time.Time) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	w := rl.hits[key]
	if w == nil || now.Sub(w.start) >= rl.window {
		rl.hits[key] = &window{count: 1, start: now}
		return true, 0
	}
	if w.count >= rl.limit {
		return false, int(math.Ceil((rl.window - now.Sub(w.start)).Seconds()))
	}
	w.count++
	return true, 0
}

// sweep evicts expired windows so the map cannot grow without bound. Runs for
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
