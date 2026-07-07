-- name: CreateRace :one
INSERT INTO races (
    id, slug, name, series, sport, distance, event_date, city, country,
    website_url, transfer_policy, official_transfer_url, policy_source_url,
    policy_notes, policy_verified_at, policy_verified_by, status, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: GetRaceByID :one
SELECT * FROM races WHERE id = $1;

-- name: GetRaceBySlug :one
SELECT r.*,
    (SELECT count(*) FROM listings l
      WHERE l.race_id = r.id AND l.status = 'active') AS active_listings
FROM races r
WHERE r.slug = $1;

-- name: ListRaces :many
SELECT r.*,
    (SELECT count(*) FROM listings l
      WHERE l.race_id = r.id AND l.status = 'active') AS active_listings
FROM races r
WHERE r.status = 'published'
  AND (sqlc.narg('country')::text IS NULL OR r.country = sqlc.narg('country'))
  AND (sqlc.narg('sport')::text IS NULL OR r.sport = sqlc.narg('sport'))
  AND (sqlc.narg('distance')::text IS NULL OR r.distance = sqlc.narg('distance'))
  AND (sqlc.narg('transfer_policy')::text IS NULL OR r.transfer_policy = sqlc.narg('transfer_policy'))
  AND (sqlc.narg('date_from')::date IS NULL OR r.event_date >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::date IS NULL OR r.event_date <= sqlc.narg('date_to'))
  AND (sqlc.narg('search')::text IS NULL
       OR to_tsvector('simple', r.name || ' ' || r.city)
          @@ plainto_tsquery('simple', sqlc.narg('search')))
  AND (sqlc.narg('cursor_date')::date IS NULL
       OR (r.event_date, r.id) > (sqlc.narg('cursor_date'), sqlc.narg('cursor_id')::uuid))
ORDER BY r.event_date, r.id
LIMIT sqlc.arg('page_size');

-- name: MapUpcomingRaces :many
-- Per-city upcoming published races for the /races map. One row per race,
-- capped to the first @per_city_limit per city (enough for the hover popover),
-- each carrying that city's full upcoming-race count (city_total). The handler
-- sums the distinct city totals into per-country totals. This replaces fetching
-- a page of full race rows (and their discarded per-row listing counts) just to
-- colour the map, and is not bounded by page size.
SELECT country, city, name, slug, city_total
FROM (
    SELECT r.country, r.city, r.name, r.slug,
        count(*) OVER (PARTITION BY r.country, r.city) AS city_total,
        row_number() OVER (PARTITION BY r.country, r.city ORDER BY r.event_date, r.id) AS rn
    FROM races r
    WHERE r.status = 'published'
      AND r.event_date >= sqlc.arg('date_from')::date
) ranked
WHERE rn <= sqlc.arg('per_city_limit')::int
ORDER BY country, city, rn;
