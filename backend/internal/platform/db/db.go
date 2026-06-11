// Package db owns the pgx connection pool.
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool builds a lazy connection pool: the URL is validated here but no
// connection happens until first use, so the API boots even while Postgres
// is still starting.
func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, cfg)
}
