package batchjob

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// countHandler records how many log records were emitted at each level.
type countHandler struct {
	mu     sync.Mutex
	levels map[slog.Level]int
}

func newCountHandler() *countHandler { return &countHandler{levels: map[slog.Level]int{}} }

func (h *countHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h *countHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	h.levels[r.Level]++
	h.mu.Unlock()
	return nil
}
func (h *countHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *countHandler) WithGroup(string) slog.Handler      { return h }

func (h *countHandler) count(l slog.Level) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.levels[l]
}

// A tick error during shutdown (ctx already canceled) must NOT log at ERROR:
// cancellation mid-batch is a clean-shutdown artifact, not an operator alarm
// (#186). Fails on the pre-#186 code, which logged every tick error.
func TestRunSkipsErrorLogOnShutdown(t *testing.T) {
	h := newCountHandler()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already shutting down before the first tick

	// A 1h ticker never fires; the pre-canceled ctx makes Run return right
	// after the single immediate tick.
	Run(ctx, slog.New(h), "job failed", "job ran", time.Hour,
		func(context.Context) (int64, bool, error) {
			return 0, false, context.Canceled
		})

	if got := h.count(slog.LevelError); got != 0 {
		t.Errorf("shutdown tick error: ERROR records = %d, want 0", got)
	}
}

// A tick error while the context is still live IS a real failure and must log
// at ERROR. Guards against the #186 fix silencing genuine errors.
func TestRunLogsRealFailure(t *testing.T) {
	h := newCountHandler()
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	// Fast ticker so the second tick arrives promptly; the first tick errors
	// with the ctx still live (logs ERROR), the second cancels so Run exits.
	Run(ctx, slog.New(h), "job failed", "job ran", time.Millisecond,
		func(context.Context) (int64, bool, error) {
			calls++
			if calls >= 2 {
				cancel()
			}
			return 0, false, errors.New("boom")
		})

	if got := h.count(slog.LevelError); got != 1 {
		t.Errorf("live-context tick error: ERROR records = %d, want 1", got)
	}
}

// A non-positive tick interval is a caller misconfiguration: Run must decline
// to start (logging at ERROR) rather than let time.NewTicker panic.
func TestRunRejectsNonPositiveInterval(t *testing.T) {
	h := newCountHandler()
	called := false
	Run(context.Background(), slog.New(h), "job failed", "job ran", 0,
		func(context.Context) (int64, bool, error) {
			called = true
			return 0, false, nil
		})
	if called {
		t.Error("Run ran the tick with every=0, want it declined to start")
	}
	if got := h.count(slog.LevelError); got != 1 {
		t.Errorf("every=0: ERROR records = %d, want 1", got)
	}
}

// Drain refuses a non-positive batch size rather than spinning forever
// re-acquiring the lock on a step that can never satisfy the "< a full batch"
// exit.
func TestDrainRejectsNonPositiveBatchSize(t *testing.T) {
	for _, bad := range []int32{0, -1} {
		n, ran, err := Drain(context.Background(), bad,
			func(context.Context) (int64, bool, error) {
				t.Fatal("step must not run for a non-positive batch size")
				return 0, false, nil
			})
		if err == nil || ran || n != 0 {
			t.Errorf("batchSize=%d: (n, ran, err) = (%d, %v, %v), want (0, false, non-nil)", bad, n, ran, err)
		}
	}
}
