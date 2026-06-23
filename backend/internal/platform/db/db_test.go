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
