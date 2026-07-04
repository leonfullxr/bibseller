package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// seedJanitorUser creates a throwaway user; sessions/tokens go with it via
// ON DELETE CASCADE.
func seedJanitorUser(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	id := ids.New()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, 'x', 'Janitor User')`,
		id, id.String()+"@test.local"); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	})
	return id
}

// seedRow inserts one row into table (sessions / email_verifications /
// password_resets - identical token_hash/user_id/expires_at shape) with the
// given expiry, returning the token hash for existence checks.
func seedRow(t *testing.T, pool *pgxpool.Pool, table string, userID uuid.UUID, expiresAt time.Time) []byte {
	t.Helper()
	hash := []byte(ids.New().String())
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO `+table+` (token_hash, user_id, expires_at) VALUES ($1, $2, $3)`,
		hash, userID, expiresAt); err != nil {
		t.Fatalf("seed %s: %v", table, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM `+table+` WHERE token_hash = $1`, hash)
	})
	return hash
}

func rowExists(t *testing.T, pool *pgxpool.Pool, table string, hash []byte) bool {
	t.Helper()
	var exists bool
	if err := pool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM `+table+` WHERE token_hash = $1)`, hash).Scan(&exists); err != nil {
		t.Fatalf("row exists %s: %v", table, err)
	}
	return exists
}

// The deletion semantics (D6: this test fails without the job): rows past
// expiry + 30 days go, everything fresher stays - including rows that are
// expired but still inside the retention buffer.
func TestPurgeExpiredAuthRows(t *testing.T) {
	pool := testdb.Pool(t)
	user := seedJanitorUser(t, pool)
	now := time.Now().UTC()

	old := now.Add(-janitorRetention - 24*time.Hour) // past retention: delete
	buffer := now.Add(-time.Hour)                    // expired, inside buffer: keep
	fresh := now.Add(24 * time.Hour)                 // live: keep
	tables := []string{"sessions", "email_verifications", "password_resets"}

	oldRows := map[string][]byte{}
	bufferRows := map[string][]byte{}
	freshRows := map[string][]byte{}
	for _, tbl := range tables {
		oldRows[tbl] = seedRow(t, pool, tbl, user, old)
		bufferRows[tbl] = seedRow(t, pool, tbl, user, buffer)
		freshRows[tbl] = seedRow(t, pool, tbl, user, fresh)
	}

	n, ran, err := purgeExpiredAuthRows(context.Background(), pool, now, janitorBatchSize)
	if err != nil {
		t.Fatalf("purge: %v", err)
	}
	if !ran {
		t.Fatal("janitor did not run (lock not held)")
	}
	if n < 3 {
		t.Errorf("deleted count = %d, want >= 3 (one old row per table)", n)
	}

	for _, tbl := range tables {
		if rowExists(t, pool, tbl, oldRows[tbl]) {
			t.Errorf("%s: row past retention survived", tbl)
		}
		if !rowExists(t, pool, tbl, bufferRows[tbl]) {
			t.Errorf("%s: expired row inside the 30-day buffer was deleted", tbl)
		}
		if !rowExists(t, pool, tbl, freshRows[tbl]) {
			t.Errorf("%s: live row was deleted", tbl)
		}
	}
}

// Batching: with batchSize=1 draining this test's three old sessions can only
// happen across multiple batches; drain across repeated calls like the other
// jobs' batch tests (shared test DB, #117).
func TestPurgeExpiredAuthRowsBatches(t *testing.T) {
	pool := testdb.Pool(t)
	user := seedJanitorUser(t, pool)
	now := time.Now().UTC()
	old := now.Add(-janitorRetention - 24*time.Hour)

	hashes := [][]byte{
		seedRow(t, pool, "sessions", user, old),
		seedRow(t, pool, "sessions", user, old),
		seedRow(t, pool, "sessions", user, old),
	}
	allGone := func() bool {
		for _, h := range hashes {
			if rowExists(t, pool, "sessions", h) {
				return false
			}
		}
		return true
	}

	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, ran, err := purgeExpiredAuthRows(context.Background(), pool, now, 1); err != nil {
			t.Fatalf("purge: %v", err)
		} else if !ran {
			t.Fatal("janitor did not run (lock not held)")
		}
		if allGone() {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("this test's sessions were not all purged after repeated batchSize=1 runs")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// The guard that stops batchSize <= 0 from spinning the batch loop forever.
func TestPurgeExpiredAuthRowsRejectsNonPositiveBatchSize(t *testing.T) {
	pool := testdb.Pool(t)
	for _, bad := range []int32{0, -1} {
		if _, ran, err := purgeExpiredAuthRows(context.Background(), pool, time.Now().UTC(), bad); err == nil || ran {
			t.Errorf("batchSize=%d: err = %v, ran = %v, want an error and ran=false", bad, err, ran)
		}
	}
}

// Advisory-lock exclusion: while another holder has the janitor lock, the job
// skips (ran=false), it neither blocks nor deletes.
func TestJanitorAdvisoryLockSerializes(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var held bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, janitorLockKey).Scan(&held); err != nil {
		t.Fatalf("acquire lock: %v", err)
	}
	if !held {
		t.Fatal("could not acquire the lock to set up the test")
	}

	n, ran, err := purgeExpiredAuthRows(ctx, pool, time.Now().UTC(), janitorBatchSize)
	if err != nil {
		t.Fatalf("purge while locked: %v", err)
	}
	if ran {
		t.Error("janitor ran while another holder had the lock")
	}
	if n != 0 {
		t.Errorf("count = %d, want 0 when lock not acquired", n)
	}
}
