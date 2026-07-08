package chat

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/batchjob"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/storage"
)

// retentionLockKey serializes the message-retention job across API instances
// (distinct from the listing expiry job's key).
const retentionLockKey int64 = 4_920_002

// retentionMonths: chat is kept until 12 months after the race's event_date
// (docs/DATA_MODEL.md retention table), then deleted.
const retentionMonths = 12

// retentionBatchSize bounds each transaction/lock hold and each batched
// object-delete call, so a large matured cohort doesn't become one huge
// DELETE plus a long sequential object-delete loop (#99).
const retentionBatchSize = 500

// StartRetentionJob deletes chat past the retention horizon on a ticker until
// ctx is done. See batchjob.Run for the loop/lock/shutdown semantics shared
// with the auth janitor and listing expiry jobs.
func StartRetentionJob(ctx context.Context, pool *pgxpool.Pool, store *storage.Client, logger *slog.Logger, every time.Duration) {
	batchjob.Run(ctx, logger, "message retention failed", "deleted expired chat messages", every,
		func(ctx context.Context) (int64, bool, error) {
			return purgeExpiredMessages(ctx, pool, store, logger, time.Now().UTC(), retentionBatchSize)
		})
}

// purgeExpiredMessages deletes messages whose race finished over retentionMonths
// ago, batchSize at a time via the shared drain, removing each batch's private
// image objects from storage before the next. Unlike the other jobs the drain
// step does post-commit work (the object delete), which is why it does not use
// batchjob.Batch directly - the object delete belongs to retention, not the
// shared skeleton (#187). Reports the total deleted and whether this instance
// held the lock.
func purgeExpiredMessages(ctx context.Context, pool *pgxpool.Pool, store *storage.Client, logger *slog.Logger, now time.Time, batchSize int32) (int64, bool, error) {
	cutoff := now.AddDate(0, -retentionMonths, 0).Truncate(24 * time.Hour)
	return batchjob.Drain(ctx, batchSize, func(ctx context.Context) (int64, bool, error) {
		n, keys, ran, err := purgeMessageBatch(ctx, pool, cutoff, batchSize)
		if err != nil || !ran {
			return n, ran, err
		}
		if len(keys) > 0 {
			// Best-effort, post-commit: the rows are already gone, so an object
			// store hiccup must not fail (and re-run) the batch. A sample key
			// anchors the failing batch to a bucket/object when debugging (the
			// batched call collapses to one error, no per-key detail).
			dctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			if derr := store.DeleteMany(dctx, keys); derr != nil {
				logger.Error("retention: batched image object delete failed", "err", derr, "count", len(keys), "sample_key", keys[0])
			}
			cancel()
		}
		return n, ran, nil
	})
}

// purgeMessageBatch deletes up to batchSize expired messages in one
// transaction holding retentionLockKey and returns their (now-orphaned) image
// keys for the caller to remove from storage after commit -
// ListExpiredMessageBatch and DeleteExpiredMessagesByID operate on the exact
// same id set, so a harvested key's message is always the one deleted. Keys
// are returned only on a committed batch (batchjob.Batch surfaces a
// commit/lock failure as err/ran), so a rolled-back batch never triggers an
// object delete for rows that still exist.
func purgeMessageBatch(ctx context.Context, pool *pgxpool.Pool, cutoff time.Time, batchSize int32) (int64, []string, bool, error) {
	var keys []string
	n, ran, err := batchjob.Batch(ctx, pool, retentionLockKey, func(ctx context.Context, q *sqlcgen.Queries) (int64, error) {
		batch, err := q.ListExpiredMessageBatch(ctx, sqlcgen.ListExpiredMessageBatchParams{Cutoff: cutoff, BatchSize: batchSize})
		if err != nil {
			return 0, err
		}
		if len(batch) == 0 {
			return 0, nil
		}
		ids := make([]uuid.UUID, len(batch))
		for i, m := range batch {
			ids[i] = m.ID
		}
		count, err := q.DeleteExpiredMessagesByID(ctx, ids)
		if err != nil {
			return 0, err
		}
		for _, m := range batch {
			if m.ImageKey != nil {
				keys = append(keys, *m.ImageKey)
			}
		}
		return count, nil
	})
	if err != nil {
		return 0, nil, ran, err // batch did not commit; harvested keys are void
	}
	return n, keys, ran, nil
}
