package listing

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/batchjob"
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
// done. See batchjob.Run for the loop/lock/shutdown semantics shared with the
// auth janitor and chat retention jobs.
func StartExpiryJob(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, every time.Duration) {
	batchjob.Run(ctx, logger, "listing expiry failed", "expired past-race listings", every,
		func(ctx context.Context) (int64, bool, error) {
			return expirePastRaceListings(ctx, pool, time.Now().UTC(), expiryBatchSize)
		})
}

// expirePastRaceListings flips active listings of finished races to 'expired',
// batchSize at a time via the shared drain (each batch its own transaction
// holding expiryLockKey). It reports the total expired and whether this
// instance held the lock.
func expirePastRaceListings(ctx context.Context, pool *pgxpool.Pool, now time.Time, batchSize int32) (int64, bool, error) {
	// Start of today UTC: a race whose day is before today is over.
	cutoff := now.Truncate(24 * time.Hour)
	return batchjob.Drain(ctx, batchSize, func(ctx context.Context) (int64, bool, error) {
		return expireListingBatch(ctx, pool, cutoff, batchSize)
	})
}

// expireListingBatch flips up to batchSize past-race listings in one
// transaction holding expiryLockKey, reporting rows affected and whether the
// lock was held (false = another instance has it).
func expireListingBatch(ctx context.Context, pool *pgxpool.Pool, cutoff time.Time, batchSize int32) (int64, bool, error) {
	return batchjob.Batch(ctx, pool, expiryLockKey, func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
		return q.ExpirePastRaceListings(ctx, sqlcgen.ExpirePastRaceListingsParams{Cutoff: cutoff, BatchSize: batchSize})
	})
}
