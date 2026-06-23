// Package db owns the pgx connection pool.
package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
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
	return pgxpool.NewWithConfig(ctx, cfg)
}
