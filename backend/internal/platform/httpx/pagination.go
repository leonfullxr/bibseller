package httpx

import (
	"fmt"
	"net/url"
	"strconv"
)

// Shared list defaults for the public catalog read endpoints (races, listings).
const (
	DefaultPageSize int32 = 24
	MaxPageSize     int32 = 100

	// CatalogCacheControl is the cache policy for anonymous catalog reads
	// (CDN/browser cacheable). Authed endpoints set no-store instead.
	CatalogCacheControl = "public, max-age=60, stale-while-revalidate=300"
)

// ParseLimit reads the "limit" query param: DefaultPageSize when absent, an
// error (whose message is caller-safe) when present but outside 1..MaxPageSize.
func ParseLimit(q url.Values) (int32, error) {
	v := q.Get("limit")
	if v == "" {
		return DefaultPageSize, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > int(MaxPageSize) {
		return 0, fmt.Errorf("limit must be 1..%d", MaxPageSize)
	}
	return int32(n), nil
}
