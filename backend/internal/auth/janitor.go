package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

// janitorLockKey serializes the auth janitor across API instances (distinct
// from the listing expiry and chat retention jobs' keys).
const janitorLockKey int64 = 4_920_003

// janitorRetention: sessions are kept until expiry + 30 days
// (docs/DATA_MODEL.md retention table); the token tables get the same buffer -
// they have no documented retention of their own, and a short post-expiry
// window costs nothing while keeping recent rows around for debugging.
const janitorRetention = 30 * 24 * time.Hour

// janitorBatchSize bounds each transaction/lock hold, so a large backlog of
// dead rows doesn't become one unbounded DELETE (#99).
const janitorBatchSize int32 = 500

// StartJanitorJob purges auth rows past the retention horizon - expired
// sessions plus stale email-verification and password-reset tokens - on a
// ticker until ctx is done (fires once immediately, then every `every`). Safe
// on every instance: a transaction-scoped advisory lock means only one does
// the work per tick (architecture invariant 5).
func StartJanitorJob(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := purgeExpiredAuthRows(ctx, pool, time.Now().UTC(), janitorBatchSize); err != nil {
			logger.Error("auth janitor failed", "err", err)
		} else if ran && n > 0 {
			logger.Info("purged expired auth rows", "count", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// purgeExpiredAuthRows deletes rows whose expires_at is older than the
// retention buffer, batchSize at a time per table (each batch its own
// transaction, releasing the advisory lock between batches). It reports the
// total deleted across all three tables and whether this instance held the
// lock for at least the first batch (false = another instance is running it
// this tick).
func purgeExpiredAuthRows(ctx context.Context, pool *pgxpool.Pool, now time.Time, batchSize int32) (count int64, ran bool, err error) {
	if batchSize <= 0 {
		// A batch that deletes 0 rows never satisfies the loop's "fewer than a
		// full batch" exit, so this would otherwise spin forever re-acquiring
		// the lock every iteration.
		return 0, false, fmt.Errorf("purgeExpiredAuthRows: batchSize must be positive, got %d", batchSize)
	}
	cutoff := now.Add(-janitorRetention)
	deletes := []func(context.Context, *sqlcgen.Queries) (int64, error){
		func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredSessionsBatch(ctx, sqlcgen.DeleteExpiredSessionsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		},
		func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredEmailVerificationsBatch(ctx, sqlcgen.DeleteExpiredEmailVerificationsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		},
		func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredPasswordResetsBatch(ctx, sqlcgen.DeleteExpiredPasswordResetsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		},
	}
	for _, del := range deletes {
		for {
			n, batchRan, err := janitorBatch(ctx, pool, del)
			if err != nil {
				return count, ran, err
			}
			if !batchRan {
				return count, ran, nil // another instance holds the lock this tick
			}
			ran = true
			count += n
			if ctx.Err() != nil {
				return count, ran, nil // shutting down; don't start another table's batch
			}
			if n < int64(batchSize) {
				break // this table is drained; on to the next
			}
		}
	}
	return count, ran, nil
}

func janitorBatch(ctx context.Context, pool *pgxpool.Pool, del func(context.Context, *sqlcgen.Queries) (int64, error)) (int64, bool, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed
	q := sqlcgen.New(tx)

	locked, err := q.TryAdvisoryXactLock(ctx, janitorLockKey)
	if err != nil {
		return 0, false, err
	}
	if !locked {
		return 0, false, nil
	}

	n, err := del(ctx, q)
	if err != nil {
		return 0, true, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, true, err
	}
	return n, true, nil
}
