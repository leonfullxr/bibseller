package chat

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
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
// ctx is done (fires once immediately, then every `every`). Safe on every
// instance: a transaction-scoped advisory lock means only one does the work per
// tick (architecture invariant 5).
func StartRetentionJob(ctx context.Context, pool *pgxpool.Pool, store Storage, logger *slog.Logger, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := purgeExpiredMessages(ctx, pool, store, logger, time.Now().UTC(), retentionBatchSize); err != nil {
			logger.Error("message retention failed", "err", err)
		} else if ran && n > 0 {
			logger.Info("deleted expired chat messages", "count", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// purgeExpiredMessages deletes messages whose race finished over retentionMonths
// ago, batchSize at a time (each batch its own transaction, releasing the
// advisory lock between batches), removing each batch's private image objects
// from storage before moving to the next. Reports the total deleted and
// whether this instance held the lock for at least the first batch.
func purgeExpiredMessages(ctx context.Context, pool *pgxpool.Pool, store Storage, logger *slog.Logger, now time.Time, batchSize int32) (count int64, ran bool, err error) {
	cutoff := now.AddDate(0, -retentionMonths, 0).Truncate(24 * time.Hour)
	for {
		n, keys, batchRan, err := purgeMessageBatch(ctx, pool, cutoff, batchSize)
		if err != nil {
			return count, ran, err
		}
		if !batchRan {
			return count, ran, nil // another instance holds the lock this tick
		}
		ran = true
		count += n

		if len(keys) > 0 {
			dctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			if derr := store.DeleteMany(dctx, keys); derr != nil {
				logger.Error("retention: batched image object delete failed", "err", derr, "count", len(keys))
			}
			cancel()
		}
		if n < int64(batchSize) || ctx.Err() != nil {
			return count, ran, nil // fewer than a full batch: nothing left; or shutting down
		}
	}
}

// purgeMessageBatch deletes up to batchSize expired messages in one
// transaction and returns their (now-orphaned) image keys for the caller to
// remove from storage after commit - ListExpiredMessageBatch and
// DeleteExpiredMessagesByID operate on the exact same id set, so a harvested
// key's message is always the one actually deleted.
func purgeMessageBatch(ctx context.Context, pool *pgxpool.Pool, cutoff time.Time, batchSize int32) (count int64, keys []string, ran bool, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed
	q := sqlcgen.New(tx)

	locked, err := q.TryAdvisoryXactLock(ctx, retentionLockKey)
	if err != nil {
		return 0, nil, false, err
	}
	if !locked {
		return 0, nil, false, nil
	}

	batch, err := q.ListExpiredMessageBatch(ctx, sqlcgen.ListExpiredMessageBatchParams{Cutoff: cutoff, BatchSize: batchSize})
	if err != nil {
		return 0, nil, true, err
	}
	if len(batch) == 0 {
		return 0, nil, true, nil
	}
	ids := make([]uuid.UUID, len(batch))
	for i, m := range batch {
		ids[i] = m.ID
	}
	count, err = q.DeleteExpiredMessagesByID(ctx, ids)
	if err != nil {
		return 0, nil, true, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, nil, true, err
	}

	for _, m := range batch {
		if m.ImageKey != nil {
			keys = append(keys, *m.ImageKey)
		}
	}
	return count, keys, true, nil
}
