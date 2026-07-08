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
		// servers, byte-identical to the `make infra` stamp (#185). Spelled out
		// literally so a drifted constant cannot hide the regression.
		const stamp = "CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())"
		if !strings.Contains(err.Error(), stamp) {
			t.Errorf("refusal message is missing the stamp SQL:\n%v", err)
		}
	})

	inTx(t, `CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())`, func(err error) {
		if err != nil {
			t.Errorf("stamped database refused: ensureDevMarker = %v, want nil", err)
		}
	})
}
