package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/batchjob"
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
// ticker until ctx is done. See batchjob.Run for the loop/lock/shutdown
// semantics shared with the listing expiry and chat retention jobs.
func StartJanitorJob(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, every time.Duration) {
	batchjob.Run(ctx, logger, "auth janitor failed", "purged expired auth rows", every,
		func(ctx context.Context) (int64, bool, error) {
			return purgeExpiredAuthRows(ctx, pool, time.Now().UTC(), janitorBatchSize)
		})
}

// purgeExpiredAuthRows deletes rows whose expires_at is older than the
// retention buffer, batchSize at a time per table via the shared drain (each
// batch its own transaction holding janitorLockKey). It reports the total
// deleted across all three tables and whether this instance held the lock.
func purgeExpiredAuthRows(ctx context.Context, pool *pgxpool.Pool, now time.Time, batchSize int32) (int64, bool, error) {
	cutoff := now.Add(-janitorRetention)
	batchOf := func(del batchjob.Del) batchjob.Step {
		return func(ctx context.Context) (int64, bool, error) {
			return batchjob.Batch(ctx, pool, janitorLockKey, del)
		}
	}
	return batchjob.Drain(ctx, batchSize,
		batchOf(func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredSessionsBatch(ctx, sqlcgen.DeleteExpiredSessionsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		}),
		batchOf(func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredEmailVerificationsBatch(ctx, sqlcgen.DeleteExpiredEmailVerificationsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		}),
		batchOf(func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
			return q.DeleteExpiredPasswordResetsBatch(ctx, sqlcgen.DeleteExpiredPasswordResetsBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		}),
	)
}
