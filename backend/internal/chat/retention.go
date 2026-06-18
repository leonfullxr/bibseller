package chat

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

// retentionLockKey serializes the message-retention job across API instances
// (distinct from the listing expiry job's key).
const retentionLockKey int64 = 4_920_002

// retentionMonths: chat is kept until 12 months after the race's event_date
// (docs/DATA_MODEL.md retention table), then deleted.
const retentionMonths = 12

// StartRetentionJob deletes chat past the retention horizon on a ticker until
// ctx is done (fires once immediately, then every `every`). Safe on every
// instance: a transaction-scoped advisory lock means only one does the work per
// tick (architecture invariant 5).
func StartRetentionJob(ctx context.Context, pool *pgxpool.Pool, store Storage, logger *slog.Logger, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := purgeExpiredMessages(ctx, pool, store, logger, time.Now().UTC()); err != nil {
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
// ago under an advisory lock, then removes their private image objects from
// storage. Reports how many it deleted and whether this instance held the lock.
func purgeExpiredMessages(ctx context.Context, pool *pgxpool.Pool, store Storage, logger *slog.Logger, now time.Time) (count int64, ran bool, err error) {
	cutoff := now.AddDate(0, -retentionMonths, 0).Truncate(24 * time.Hour)

	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed
	q := sqlcgen.New(tx)

	locked, err := q.TryAdvisoryXactLock(ctx, retentionLockKey)
	if err != nil {
		return 0, false, err
	}
	if !locked {
		return 0, false, nil // another instance holds the lock this tick
	}

	keys, err := q.ListExpiredMessageImageKeys(ctx, cutoff)
	if err != nil {
		return 0, true, err
	}
	count, err = q.DeleteExpiredMessages(ctx, cutoff)
	if err != nil {
		return 0, true, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, true, err
	}

	// Rows are gone; remove their objects (best-effort - an orphan is logged,
	// not fatal). Done after commit so a storage hiccup never blocks the purge.
	for _, key := range keys {
		if key == nil {
			continue
		}
		if ctx.Err() != nil {
			break // shutting down; the rows are already gone, leave the rest
		}
		dctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		if derr := store.Delete(dctx, *key); derr != nil {
			logger.Error("retention: image object delete failed", "err", derr, "key", *key)
		}
		cancel()
	}
	return count, true, nil
}
