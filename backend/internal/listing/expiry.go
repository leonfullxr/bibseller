package listing

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

// expiryLockKey is the advisory-lock key that serializes the expiry job across
// API instances. Each background job picks a distinct, stable arbitrary int64;
// this one is the listing expiry job.
const expiryLockKey int64 = 4_920_001

// StartExpiryJob runs the past-race listing expiry on a ticker until ctx is
// done (it fires once immediately, then every `every`). It is safe to start on
// every API instance: a transaction-scoped advisory lock ensures only one
// instance does the work each tick (architecture invariant 5).
func StartExpiryJob(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := expirePastRaceListings(ctx, pool, time.Now().UTC()); err != nil {
			logger.Error("listing expiry failed", "err", err)
		} else if ran && n > 0 {
			logger.Info("expired past-race listings", "count", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// expirePastRaceListings flips active listings of finished races to 'expired'
// under a transaction-scoped advisory lock, so concurrent instances don't
// double-run. It reports how many it expired and whether this instance held the
// lock (false = another instance is running it this tick).
func expirePastRaceListings(ctx context.Context, pool *pgxpool.Pool, now time.Time) (count int64, ran bool, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed
	q := sqlcgen.New(tx)

	locked, err := q.TryAdvisoryXactLock(ctx, expiryLockKey)
	if err != nil {
		return 0, false, err
	}
	if !locked {
		return 0, false, nil // another instance holds the lock this tick
	}

	// Start of today UTC: a race whose day is before today is over.
	cutoff := now.Truncate(24 * time.Hour)
	count, err = q.ExpirePastRaceListings(ctx, cutoff)
	if err != nil {
		return 0, true, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, true, err
	}
	return count, true, nil
}
