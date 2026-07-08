package main

import (
	"context"
	"strings"
	"testing"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
)

// The wipe guard (#159): seed must refuse any database that does not carry
// the dev_marker stamp (created by `make infra`'s stamp step; prod and
// staging never get one). Each case toggles the marker inside a transaction
// that is rolled back - DDL is transactional in Postgres - so the shared
// test database is never actually mutated.
func TestEnsureDevMarker(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	inTx := func(t *testing.T, setup string, fn func(err error)) {
		t.Helper()
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("begin: %v", err)
		}
		defer func() { _ = tx.Rollback(ctx) }()
		if _, err := tx.Exec(ctx, setup); err != nil {
			t.Fatalf("setup %q: %v", setup, err)
		}
		fn(ensureDevMarker(ctx, tx))
	}

	inTx(t, `DROP TABLE IF EXISTS dev_marker`, func(err error) {
		if err == nil {
			t.Fatal("unstamped database accepted: ensureDevMarker = nil, want an error")
		}
		// The refusal must carry the copy-pasteable stamp for local/non-compose
		// servers, byte-identical to the `make infra` stamp (#185). Assert the
		// exact printed line (two-space indent, no trailing punctuation) spelled
		// out literally, so a drifted constant or reformatted message is caught.
		const wantLine = "  CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())"
		if !strings.Contains(err.Error(), wantLine) {
			t.Errorf("refusal message is missing the exact stamp line:\n%v", err)
		}
	})

	inTx(t, `CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())`, func(err error) {
		if err != nil {
			t.Errorf("stamped database refused: ensureDevMarker = %v, want nil", err)
		}
	})
}

// The -guard-only path (make guard-dev-db) gates migrate/migrate-down (#184),
// so the check itself must never write: it runs before goose against a target
// that might be prod. Prove ensureDevMarker leaves the marker table's rows
// untouched (it only reads via to_regclass). Rolled back, so the shared test
// DB is never mutated.
func TestEnsureDevMarkerIsReadOnly(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		t.Fatalf("stamp: %v", err)
	}
	// Seed a row so the assertion is "row count unchanged", not "empty" - the
	// latter would pass even if the marker table legitimately held rows.
	if _, err := tx.Exec(ctx, `INSERT INTO dev_marker DEFAULT VALUES`); err != nil {
		t.Fatalf("seed row: %v", err)
	}
	countRows := func() int {
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM dev_marker`).Scan(&n); err != nil {
			t.Fatalf("count: %v", err)
		}
		return n
	}
	before := countRows()
	if err := ensureDevMarker(ctx, tx); err != nil {
		t.Fatalf("stamped database refused: %v", err)
	}
	if after := countRows(); after != before {
		t.Errorf("ensureDevMarker mutated the target: dev_marker rows %d -> %d, want unchanged", before, after)
	}
}

// The -guard-only path (make guard-dev-db) must fail closed: if it cannot
// reach and verify the target, migrate/migrate-down must abort rather than run
// blind. An unreachable DB is the cheap, deterministic proxy for "cannot
// verify" (#184).
func TestRunGuardOnlyFailsClosed(t *testing.T) {
	if err := runGuardOnly("postgres://nope:nope@127.0.0.1:1/nope?sslmode=disable"); err == nil {
		t.Error("runGuardOnly against an unreachable DB = nil, want an error so migrate aborts")
	}
}
