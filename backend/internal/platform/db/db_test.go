package db

import (
	"context"
	"strings"
	"testing"
)

// A malformed DATABASE_URL must never surface its contents in the returned
// error (#63). pgx redacts a *recognized* password field, but echoes the raw
// connection string verbatim when it cannot structure it - and cmd/api logs
// startup errors verbatim. No DB connection is needed: parsing fails first.
func TestNewPoolDoesNotLeakConnStringOnBadURL(t *testing.T) {
	const secret = "le4kedSecret"
	_, err := NewPool(context.Background(), secret)
	if err == nil {
		t.Fatal("expected an error for the unparseable connection string")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error leaks the connection string: %v", err)
	}
}

// The pool must carry an explicit MaxConns, not pgx's silent max(4, NumCPU)
// default, and it must stay under Postgres' max_connections across all
// instances or the DB refuses connections. Asserting both ties the two numbers
// so they can't drift apart unnoticed (#93). Lazy pool: no DB needed - we only
// read the parsed config.
func TestNewPoolSetsExplicitMaxConns(t *testing.T) {
	pool, err := NewPool(context.Background(), "postgres://u:p@127.0.0.1:9/db")
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer pool.Close()

	if got := pool.Config().MaxConns; got != maxConns {
		t.Fatalf("MaxConns = %d, want the explicit cap %d (not pgx's silent default)", got, maxConns)
	}
	const instances = 1 // v1 single instance; bump with #100 when scaling out
	if maxConns*instances >= postgresMaxConnections {
		t.Fatalf("maxConns(%d) x instances(%d) must stay under Postgres max_connections(%d)",
			maxConns, instances, postgresMaxConnections)
	}
}
