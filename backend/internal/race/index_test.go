package race_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// bulkRaceCount is large enough that Postgres' cost model actually prefers an
// index over a full scan + sort - a handful of seed rows never would, so #95's
// "no full-table sort" claim needs a catalog this size to mean anything.
const bulkRaceCount = 5000

// seedBulkRaces inserts bulkRaceCount published races spread across countries,
// sports and two years of event dates, batched in one round trip. Self-cleaning
// via its distinct slug prefix, safe against any other DB state.
func seedBulkRaces(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	countries := []string{"ES", "FR", "DE", "IT", "PT"}
	sports := []string{"running", "trail", "triathlon", "cycling"}
	src := "https://example.org/source"

	batch := &pgx.Batch{}
	for i := range bulkRaceCount {
		id := ids.New()
		batch.Queue(
			`INSERT INTO races (id, slug, name, sport, event_date, city, country, transfer_policy, policy_source_url, official_transfer_url, status)
			 VALUES ($1, $2, $3, $4, '2027-01-01'::date + ($5 || ' days')::interval, $6, $7, 'platform_sale', $8, $8, 'published')`,
			id, "bulk-"+id.String(), "Bulk Race "+strconv.Itoa(i),
			sports[i%len(sports)], strconv.Itoa(i%730), "City"+strconv.Itoa(i%50), countries[i%len(countries)], src,
		)
	}
	if err := pool.SendBatch(ctx, batch).Close(); err != nil {
		t.Fatalf("bulk seed races: %v", err)
	}
	// The planner needs fresh stats to know the table just got 5000 rows bigger;
	// without this it costs plans off the pre-insert (near-empty) statistics.
	if _, err := pool.Exec(ctx, `ANALYZE races`); err != nil {
		t.Fatalf("analyze races: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM races WHERE slug LIKE 'bulk-%'`)
	})
}

// explainPlan runs the query under EXPLAIN ANALYZE and returns the plan text.
func explainPlan(t *testing.T, pool *pgxpool.Pool, query string, args ...any) string {
	t.Helper()
	rows, err := pool.Query(context.Background(), "EXPLAIN (ANALYZE, FORMAT TEXT) "+query, args...)
	if err != nil {
		t.Fatalf("explain: %v", err)
	}
	defer rows.Close()
	var b strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatalf("explain scan: %v", err)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

// listRacesSQL mirrors ListRaces (backend/db/queries/race.sql) - EXPLAIN needs
// the literal query text, which sqlc keeps unexported.
const listRacesSQL = `
SELECT r.id, r.slug, r.name, r.series, r.sport, r.distance, r.event_date, r.city, r.country, r.website_url, r.transfer_policy, r.official_transfer_url, r.policy_source_url, r.policy_notes, r.policy_verified_at, r.policy_verified_by, r.status, r.created_by, r.created_at, r.updated_at,
    (SELECT count(*) FROM listings l
      WHERE l.race_id = r.id AND l.status = 'active') AS active_listings
FROM races r
WHERE r.status = 'published'
  AND ($1::text IS NULL OR r.country = $1)
  AND ($2::text IS NULL OR r.sport = $2)
  AND ($3::text IS NULL OR r.transfer_policy = $3)
  AND ($4::date IS NULL OR r.event_date >= $4)
  AND ($5::date IS NULL OR r.event_date <= $5)
  AND ($6::text IS NULL
       OR to_tsvector('simple', r.name || ' ' || r.city)
          @@ plainto_tsquery('simple', $6))
  AND ($7::date IS NULL
       OR (r.event_date, r.id) > ($7, $8::uuid))
ORDER BY r.event_date, r.id
LIMIT $9
`

func TestListRacesUsesIndexNotFullSort(t *testing.T) {
	pool := testdb.Pool(t)
	seedBulkRaces(t, pool)

	// Default, unfiltered browse (#95 acceptance criterion 1).
	plan := explainPlan(t, pool, listRacesSQL,
		nil, nil, nil, nil, nil, nil, nil, uuid.Nil, 24)
	if !strings.Contains(plan, "Index") {
		t.Errorf("default browse: expected an index scan, got:\n%s", plan)
	}
	if strings.Contains(plan, "Sort") {
		t.Errorf("default browse: expected no full-table sort, got:\n%s", plan)
	}

	// Sport-filtered browse (#95 acceptance criterion 2).
	sport := "trail"
	plan = explainPlan(t, pool, listRacesSQL,
		nil, &sport, nil, nil, nil, nil, nil, uuid.Nil, 24)
	if !strings.Contains(plan, "Index") {
		t.Errorf("sport-filtered browse: expected an index scan, got:\n%s", plan)
	}
	if strings.Contains(plan, "Sort") {
		t.Errorf("sport-filtered browse: expected no full-table sort, got:\n%s", plan)
	}
}
