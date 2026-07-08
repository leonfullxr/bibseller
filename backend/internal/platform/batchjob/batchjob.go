// Package batchjob is the ticker/advisory-lock/batch scaffold shared by the
// listing expiry, chat retention, and auth janitor jobs (#187; formerly three
// verbatim copies). Domain packages keep their lock keys, retention constants,
// and per-batch closures; this package owns the loop shapes: the ticker (Run),
// the batch-by-batch drain with the spin guard (Drain), and the
// one-transaction-holding-the-advisory-lock batch (Batch).
package batchjob

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

// Step performs one unit of a job's work - one committed batch inside Drain,
// or a whole drain when used as Run's tick - reporting rows affected and
// whether this instance actually did the work (ran=false: another instance
// holds the advisory lock this tick).
type Step func(ctx context.Context) (count int64, ran bool, err error)

// Del deletes (or flips) up to one batch of rows using q, reporting how many
// rows it affected.
type Del func(ctx context.Context, q *sqlcgen.Queries) (int64, error)

// Run runs tick until ctx is done (it fires once immediately, then every
// `every`). Safe to start on every API instance: a transaction-scoped
// advisory lock inside each batch ensures only one instance does the work per
// tick (architecture invariant 5). Tick errors log at ERROR under failMsg -
// except while shutting down: cancellation mid-batch surfaces as a batch
// error on an otherwise clean shutdown, not a failure worth an operator's
// attention (#186). Ticks that did work log at INFO under doneMsg.
func Run(ctx context.Context, logger *slog.Logger, failMsg, doneMsg string, every time.Duration, tick Step) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		if n, ran, err := tick(ctx); err != nil {
			if ctx.Err() == nil { // a real failure, not shutdown (#186)
				logger.Error(failMsg, "err", err)
			}
		} else if ran && n > 0 {
			logger.Info(doneMsg, "count", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// Drain runs each step batch-by-batch until it affects fewer than a full
// batch, then moves to the next (each batch its own transaction, releasing
// the advisory lock between batches so a large backlog never holds one huge
// lock set, #99). It reports the total affected and whether this instance
// held the lock for at least the first batch (false = another instance is
// running it this tick), stopping early on shutdown rather than starting
// another batch.
func Drain(ctx context.Context, batchSize int32, steps ...Step) (count int64, ran bool, err error) {
	if batchSize <= 0 {
		// A batch that affects 0 rows never satisfies the loop's "fewer than a
		// full batch" exit, so this would otherwise spin forever re-acquiring
		// the lock every iteration.
		return 0, false, fmt.Errorf("batchjob: batchSize must be positive, got %d", batchSize)
	}
	for _, step := range steps {
		for {
			n, stepRan, err := step(ctx)
			if err != nil {
				return count, ran, err
			}
			if !stepRan {
				return count, ran, nil // another instance holds the lock this tick
			}
			ran = true
			count += n
			if ctx.Err() != nil {
				return count, ran, nil // shutting down; don't start another batch
			}
			if n < int64(batchSize) {
				break // this step is drained; on to the next
			}
		}
	}
	return count, ran, nil
}

// Batch runs del in its own transaction holding the job's advisory lock,
// reporting rows affected and whether the lock was held (false = another
// instance has it; nothing ran).
func Batch(ctx context.Context, pool *pgxpool.Pool, lockKey int64, del Del) (int64, bool, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed
	q := sqlcgen.New(tx)

	locked, err := q.TryAdvisoryXactLock(ctx, lockKey)
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
