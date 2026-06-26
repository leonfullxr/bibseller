// Package db owns the pgx connection pool.
package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// maxConns caps the pool per API instance. pgx otherwise silently defaults
	// MaxConns to max(4, NumCPU) - the real backend concurrency ceiling, set by
	// accident. v1 runs a single instance (docs/ARCHITECTURE.md), so this stays
	// well under postgresMaxConnections with headroom for migrations, backups and
	// psql. Raise this and Postgres max_connections together if we ever scale out
	// (#93, #100).
	maxConns = 20

	// postgresMaxConnections is stock postgres:16's default max_connections. Kept
	// here only so a test can assert maxConns x instances stays under it; it is
	// not applied to Postgres (resource limits live with #94/compose.prod.yml).
	postgresMaxConnections = 100
)

// NewPool builds a lazy connection pool: the URL is validated here but no
// connection happens until first use, so the API boots even while Postgres
// is still starting.
func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		// Never wrap err: pgx echoes the raw connection string (which can carry
		// the password) in its parse error, and cmd/api logs startup errors
		// verbatim - a secret must never reach the logs (#63).
		return nil, errors.New("invalid DATABASE_URL")
	}
	cfg.MaxConns = maxConns
	return pgxpool.NewWithConfig(ctx, cfg)
}
