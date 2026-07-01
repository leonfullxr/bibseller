package listing

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// seedRaceListing creates a race on eventDate with one listing, and (optionally)
// transitions the listing out of 'active'. Returns the listing id.
func seedRaceListing(t *testing.T, pool *pgxpool.Pool, sellerID uuid.UUID, eventDate time.Time, cancel bool) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	q := sqlcgen.New(pool)
	src := "https://example.org/source"
	raceID := ids.New()
	race, err := q.CreateRace(ctx, sqlcgen.CreateRaceParams{
		ID: raceID, Slug: "ex-" + raceID.String(), Name: "Expiry Race", Sport: "running",
		EventDate: eventDate, City: "Testville", Country: "ZX",
		TransferPolicy: "platform_sale", PolicySourceUrl: &src, Status: "published",
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	price, orig := int32(5000), int32(6000)
	l, err := q.CreateListing(ctx, sqlcgen.CreateListingParams{
		ID: ids.New(), RaceID: race.ID, SellerID: sellerID,
		PriceCents: &price, Currency: "EUR", OriginalPriceCents: &orig,
		ExpiresAt: race.EventDate,
	})
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
	if cancel {
		if _, err := q.UpdateListingStatus(ctx, sqlcgen.UpdateListingStatusParams{
			ID: l.ID, FromStatus: "active", ToStatus: "cancelled",
		}); err != nil {
			t.Fatalf("cancel seed listing: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id = $1`, l.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
	})
	return l.ID
}

func statusOf(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM listings WHERE id = $1`, id).Scan(&s); err != nil {
		t.Fatalf("status of %s: %v", id, err)
	}
	return s
}

func seedSeller(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	id := ids.New()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, 'x', 'Expiry Seller')`,
		id, id.String()+"@test.local"); err != nil {
		t.Fatalf("seed seller: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	})
	return id
}

func TestExpirePastRaceListings(t *testing.T) {
	pool := testdb.Pool(t)
	seller := seedSeller(t, pool)

	pastActive := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 0, -3), false)
	futureActive := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 1, 0), false)
	pastCancelled := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 0, -3), true)

	n, ran, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC(), expiryBatchSize)
	if err != nil {
		t.Fatalf("expire: %v", err)
	}
	if !ran {
		t.Fatal("expiry did not run (lock not held)")
	}
	if n < 1 {
		t.Errorf("expired count = %d, want >= 1", n)
	}

	if s := statusOf(t, pool, pastActive); s != "expired" {
		t.Errorf("past active listing = %q, want expired", s)
	}
	if s := statusOf(t, pool, futureActive); s != "active" {
		t.Errorf("future active listing = %q, want active", s)
	}
	if s := statusOf(t, pool, pastCancelled); s != "cancelled" {
		t.Errorf("past cancelled listing = %q, want cancelled (untouched)", s)
	}

	// A race happening today is not past: re-running expires nothing new from
	// the today boundary. (Sanity: the function is idempotent for this set.)
	if _, _, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC(), expiryBatchSize); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if s := statusOf(t, pool, futureActive); s != "active" {
		t.Errorf("future listing changed on second run: %q", s)
	}
}

// TestExpirePastRaceListingsBatches proves #99: a backlog bigger than one
// batch is fully drained within a single call, one batch at a time - a
// batchSize of 1 against 3 due listings can only reach n == 3 by running
// three internal batches.
func TestExpirePastRaceListingsBatches(t *testing.T) {
	pool := testdb.Pool(t)
	seller := seedSeller(t, pool)

	past := time.Now().UTC().AddDate(0, 0, -3)
	ids := []uuid.UUID{
		seedRaceListing(t, pool, seller, past, false),
		seedRaceListing(t, pool, seller, past, false),
		seedRaceListing(t, pool, seller, past, false),
	}

	n, ran, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC(), 1)
	if err != nil {
		t.Fatalf("expire: %v", err)
	}
	if !ran {
		t.Fatal("expiry did not run (lock not held)")
	}
	// >= not == : testdb.Pool runs against a database shared with other
	// packages' tests (CI runs them concurrently), so an unrelated past-dated
	// active listing can legitimately also get swept up here. What proves
	// multi-batch draining is that batchSize=1 still reaches every one of
	// *this test's* ids (checked below) in a single call.
	if n < int64(len(ids)) {
		t.Fatalf("expired count = %d, want >= %d (batchSize=1 forces multiple batches)", n, len(ids))
	}
	for _, id := range ids {
		if s := statusOf(t, pool, id); s != "expired" {
			t.Errorf("listing %s = %q, want expired", id, s)
		}
	}
}

// TestExpirePastRaceListingsRejectsNonPositiveBatchSize proves the guard that
// stops a batchSize <= 0 from spinning the loop forever (0 rows expired never
// satisfies the "fewer than a full batch" exit).
func TestExpirePastRaceListingsRejectsNonPositiveBatchSize(t *testing.T) {
	pool := testdb.Pool(t)
	for _, bad := range []int32{0, -1} {
		if _, ran, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC(), bad); err == nil || ran {
			t.Errorf("batchSize=%d: err = %v, ran = %v, want an error and ran=false", bad, err, ran)
		}
	}
}

// waitForRowLockContention polls pg_stat_activity until another backend is
// actually blocked waiting on a row lock, so a test that depends on lock
// contention isn't just hoping a goroutine got scheduled in time. A session
// waiting on a row another transaction holds shows up as wait_event_type
// 'Lock' / wait_event 'transactionid' (Postgres implements row-lock waits as
// waiting on the blocking transaction's xid, not a tuple-level pg_locks row).
func waitForRowLockContention(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var waiting bool
		if err := pool.QueryRow(context.Background(), `
			SELECT EXISTS (
				SELECT 1 FROM pg_stat_activity
				WHERE wait_event_type = 'Lock' AND wait_event = 'transactionid'
			)`).Scan(&waiting); err == nil && waiting {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("timed out waiting for a blocked lock on listings")
}

// TestExpireListingRechecksStatusOnConcurrentChange proves the outer WHERE's
// l.status = 'active' recheck (#99 round 2, Copilot review): batching by id
// alone isn't safe under concurrent writes. If a listing changes away from
// 'active' after the batch's candidate-id subquery already ran but before its
// row lock is free, Postgres only re-evaluates conditions that are literally
// in the outer UPDATE's WHERE clause (EvalPlanQual) - id membership alone
// would still match and force the row to 'expired', silently overwriting
// whatever the concurrent change was.
func TestExpireListingRechecksStatusOnConcurrentChange(t *testing.T) {
	pool := testdb.Pool(t)
	seller := seedSeller(t, pool)
	listingID := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 0, -3), false)

	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer conn.Release()
	lockTx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatalf("begin lock tx: %v", err)
	}
	defer func() { _ = lockTx.Rollback(ctx) }()
	if _, err := lockTx.Exec(ctx, `SELECT * FROM listings WHERE id = $1 FOR UPDATE`, listingID); err != nil {
		t.Fatalf("lock row: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _, _ = expireListingBatch(context.Background(), pool, time.Now().UTC().Truncate(24*time.Hour), 500)
	}()

	// The batch's UPDATE is now genuinely blocked on this row's lock - only
	// once we're sure of that do we change status and release it, so this
	// test can't pass for the wrong reason (the subquery alone excluding an
	// already-cancelled row, without ever exercising the outer recheck).
	waitForRowLockContention(t, pool)
	if _, err := lockTx.Exec(ctx, `UPDATE listings SET status = 'cancelled' WHERE id = $1`, listingID); err != nil {
		t.Fatalf("concurrent status change: %v", err)
	}
	if err := lockTx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("expiry batch never returned after the lock was released")
	}

	if s := statusOf(t, pool, listingID); s != "cancelled" {
		t.Errorf("listing = %q, want cancelled (the concurrent change must win, not get overwritten to expired)", s)
	}
}

func TestExpiryAdvisoryLockSerializes(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	// Hold the expiry lock on a separate connection/transaction.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var held bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, expiryLockKey).Scan(&held); err != nil {
		t.Fatalf("acquire lock: %v", err)
	}
	if !held {
		t.Fatal("could not acquire the lock to set up the test")
	}

	// The job must see the lock as taken and skip, not block or run.
	n, ran, err := expirePastRaceListings(ctx, pool, time.Now().UTC(), expiryBatchSize)
	if err != nil {
		t.Fatalf("expire while locked: %v", err)
	}
	if ran {
		t.Error("expiry ran while another holder had the lock")
	}
	if n != 0 {
		t.Errorf("count = %d, want 0 when lock not acquired", n)
	}
}
