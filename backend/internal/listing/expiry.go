package listing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

// expiryLockKey is the advisory-lock key that serializes the expiry job across
// API instances. Each background job picks a distinct, stable arbitrary int64;
// this one is the listing expiry job.
const expiryLockKey int64 = 4_920_001

// expiryBatchSize bounds each transaction/lock hold, so a large backlog of
// past-race listings doesn't become one unbounded UPDATE holding every
// matching row lock at once (#99).
const expiryBatchSize int32 = 500

// StartExpiryJob runs the past-race listing expiry on a ticker until ctx is
// done (it fires once immediately, then every `every`). It is safe to start on
// every API instance: a transaction-scoped advisory lock ensures only one
// instance does the work each tick (architecture invariant 5).
func StartExpiryJob(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := expirePastRaceListings(ctx, pool, time.Now().UTC(), expiryBatchSize); err != nil {
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

// expirePastRaceListings flips active listings of finished races to 'expired',
// batchSize at a time (each batch its own transaction, releasing the advisory
// lock between batches so a large backlog never holds one huge lock set). It
// reports the total expired and whether this instance held the lock for at
// least the first batch (false = another instance is running it this tick).
func expirePastRaceListings(ctx context.Context, pool *pgxpool.Pool, now time.Time, batchSize int32) (count int64, ran bool, err error) {
	if batchSize <= 0 {
		// A batch that expires 0 rows never satisfies the loop's "fewer than a
		// full batch" exit, so this would otherwise spin forever re-acquiring
		// the lock every iteration.
		return 0, false, fmt.Errorf("expirePastRaceListings: batchSize must be positive, got %d", batchSize)
	}
	// Start of today UTC: a race whose day is before today is over.
	cutoff := now.Truncate(24 * time.Hour)
	for {
		n, batchRan, err := expireListingBatch(ctx, pool, cutoff, batchSize)
		if err != nil {
			return count, ran, err
		}
		if !batchRan {
			return count, ran, nil
		}
		ran = true
		count += n
		if n < int64(batchSize) || ctx.Err() != nil {
			return count, ran, nil
		}
	}
}

func expireListingBatch(ctx context.Context, pool *pgxpool.Pool, cutoff time.Time, batchSize int32) (int64, bool, error) {
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

	n, err := q.ExpirePastRaceListings(ctx, sqlcgen.ExpirePastRaceListingsParams{Cutoff: cutoff, BatchSize: batchSize})
	if err != nil {
		return 0, true, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, true, err
	}
	return n, true, nil
}
