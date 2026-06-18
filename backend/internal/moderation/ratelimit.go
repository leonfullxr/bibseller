package moderation

import (
	"math"
	"sync"
	"time"
)

// Per-account budget on the safety write endpoints (report, block).
const (
	reportRateMax    = 20
	reportRateWindow = time.Minute
)

// rateLimiter is a per-key fixed-window counter, in-process (no external store,
// per docs/ARCHITECTURE.md).
//
// ponytail: third copy of this pattern (auth, chat, here) - fold into a shared
// internal/platform/ratelimit in a dedicated refactor rather than a fourth.
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
