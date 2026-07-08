package main

import (
	"context"
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
			t.Error("unstamped database accepted: ensureDevMarker = nil, want an error")
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
	if err := ensureDevMarker(ctx, tx); err != nil {
		t.Fatalf("stamped database refused: %v", err)
	}
	var rows int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM dev_marker`).Scan(&rows); err != nil {
		t.Fatalf("count: %v", err)
	}
	if rows != 0 {
		t.Errorf("ensureDevMarker wrote to the target: dev_marker rows = %d, want 0", rows)
	}
}
