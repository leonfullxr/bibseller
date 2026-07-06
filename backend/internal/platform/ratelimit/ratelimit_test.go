package ratelimit

import (
	"testing"
	"time"
)

// The limiter is pure given an injected clock, so this needs no DB or HTTP.
func TestLimiterFixedWindow(t *testing.T) {
	rl := New(2, time.Minute)
	t0 := time.Now()

	if ok, _ := rl.Allow("1.2.3.4", t0); !ok {
		t.Fatal("1st request denied")
	}
	if ok, _ := rl.Allow("1.2.3.4", t0); !ok {
		t.Fatal("2nd request denied")
	}
	ok, retry := rl.Allow("1.2.3.4", t0)
	if ok {
		t.Fatal("3rd request allowed past the limit")
	}
	if retry <= 0 || retry > 60 {
		t.Errorf("Retry-After = %d, want 1..60", retry)
	}

	// A different key has its own budget.
	if ok, _ := rl.Allow("5.6.7.8", t0); !ok {
		t.Error("unrelated key throttled")
	}

	// The window resets once it elapses.
	if ok, _ := rl.Allow("1.2.3.4", t0.Add(time.Minute)); !ok {
		t.Error("still throttled after the window elapsed")
	}
}
