-- name: CreateListing :one
INSERT INTO listings (
    id, race_id, seller_id, price_cents, currency, original_price_cents,
    description, image_key, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateListing :one
UPDATE listings
SET price_cents = $2, original_price_cents = $3, description = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetListingByID :one
SELECT sqlc.embed(listings), sqlc.embed(races), u.display_name AS seller_name
FROM listings
JOIN races ON races.id = listings.race_id
JOIN users u ON u.id = listings.seller_id
WHERE listings.id = $1;

-- name: ListActiveListingsByRace :many
SELECT l.*, u.display_name AS seller_name
FROM listings l
JOIN users u ON u.id = l.seller_id
WHERE l.race_id = $1
  AND l.status = 'active'
  AND (sqlc.narg('cursor_id')::uuid IS NULL OR l.id < sqlc.narg('cursor_id'))
ORDER BY l.id DESC
LIMIT sqlc.arg('page_size');

-- name: ListListingsBySeller :many
SELECT l.*, r.name AS race_name, r.slug AS race_slug, r.event_date
FROM listings l
JOIN races r ON r.id = l.race_id
WHERE l.seller_id = $1
ORDER BY l.created_at DESC;

-- name: UpdateListingStatus :one
UPDATE listings
SET status = sqlc.arg('to_status'), updated_at = now()
WHERE id = sqlc.arg('id') AND status = sqlc.arg('from_status')
RETURNING *;
