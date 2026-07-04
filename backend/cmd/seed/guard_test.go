package main

import (
	"context"
	"testing"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
)

// The wipe guard (#159): seed must refuse any database that does not carry
// the dev_marker stamp (created by `make infra`'s stamp step; prod and
// staging never get one). The test toggles the marker on the shared test DB
// and restores whatever state it found.
func TestEnsureDevMarker(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	var had bool
	if err := pool.QueryRow(ctx,
		`SELECT to_regclass('public.dev_marker') IS NOT NULL`).Scan(&had); err != nil {
		t.Fatalf("probe marker: %v", err)
	}
	t.Cleanup(func() {
		// A silent restore failure would strip a dev machine's marker and make
		// the next `make seed` refuse mysteriously - surface it.
		var err error
		if had {
			_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())`)
		} else {
			_, err = pool.Exec(ctx, `DROP TABLE IF EXISTS dev_marker`)
		}
		if err != nil {
			t.Errorf("restore dev_marker state (had=%v): %v - re-run `make infra` if seed now refuses", had, err)
		}
	})

	if _, err := pool.Exec(ctx, `DROP TABLE IF EXISTS dev_marker`); err != nil {
		t.Fatalf("drop marker: %v", err)
	}
	if err := ensureDevMarker(ctx, pool); err == nil {
		t.Error("unstamped database accepted: ensureDevMarker = nil, want an error")
	}

	if _, err := pool.Exec(ctx, `CREATE TABLE dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		t.Fatalf("create marker: %v", err)
	}
	if err := ensureDevMarker(ctx, pool); err != nil {
		t.Errorf("stamped database refused: ensureDevMarker = %v, want nil", err)
	}
}
