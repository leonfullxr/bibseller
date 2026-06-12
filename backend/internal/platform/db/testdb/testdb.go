// Package testdb hands integration tests a pool against the migrated dev/CI
// database, skipping the test cleanly when none is reachable.
package testdb

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/config"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db"
)

func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := db.NewPool(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatalf("pool config: %v", err)
	}
	t.Cleanup(pool.Close)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		t.Skipf("postgres not reachable, skipping integration test: %v", err)
	}

	var migrated bool
	if err := pool.QueryRow(ctx, `SELECT to_regclass('races') IS NOT NULL`).Scan(&migrated); err != nil || !migrated {
		t.Skip("schema not migrated — run `make migrate` first")
	}
	return pool
}
