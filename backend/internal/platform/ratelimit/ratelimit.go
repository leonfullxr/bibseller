// Package ratelimit is the per-key fixed-window counter shared by the auth,
// chat, and moderation endpoints (#122; formerly three verbatim copies).
// In-process, no external store (docs/ARCHITECTURE.md -> no new infra): v1
// runs a single API instance, and a fixed window resets without per-request
// bookkeeping. The trigger to revisit (a central store) is horizontal scaling
// of the API (#100).
package ratelimit

import (
	"math"
	"sync"
	"time"
)

// Limiter counts hits per key inside a fixed window.
type Limiter struct {
	mu     sync.Mutex
	hits   map[string]*window
	limit  int
	window time.Duration
}

type window struct {
	count int
	start time.Time
}

func New(limit int, w time.Duration) *Limiter {
	return &Limiter{hits: map[string]*window{}, limit: limit, window: w}
}

// Allow reports whether key may proceed at time now, and (when it may not) the
// whole seconds until its window resets - the Retry-After hint.
func (l *Limiter) Allow(key string, now time.Time) (bool, int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	w := l.hits[key]
	if w == nil || now.Sub(w.start) >= l.window {
		l.hits[key] = &window{count: 1, start: now}
		return true, 0
	}
	if w.count >= l.limit {
		// Whole seconds (rounded up, >=1) until this window resets.
		return false, int(math.Ceil((l.window - now.Sub(w.start)).Seconds()))
	}
	w.count++
	return true, 0
}

// Sweep evicts expired windows so the map cannot grow without bound. Runs for
// the life of the process - call it in a goroutine.
func (l *Limiter) Sweep(every time.Duration) {
	for range time.Tick(every) {
		l.mu.Lock()
		now := time.Now()
		for k, w := range l.hits {
			if now.Sub(w.start) >= l.window {
				delete(l.hits, k)
			}
		}
		l.mu.Unlock()
	}
}
