-- name: CreateThread :one
-- Find-or-create the (listing, buyer) thread, creating only while the listing is
-- still active at write time. This closes the TOCTOU between the handler's
-- active-check and this insert; ON CONFLICT DO NOTHING avoids writing a dead row
-- when the thread already exists (no MVCC churn on re-contact). The row comes
-- back with its real last_message_at - NULL only for a brand-new thread, which
-- is how the first message is detected. Zero rows means "no thread yet and the
-- listing is not active", which the handler maps to 409.
WITH created AS (
    INSERT INTO chat_threads (id, listing_id, buyer_id)
    SELECT $1, $2, $3
    WHERE EXISTS (SELECT 1 FROM listings WHERE id = $2 AND status = 'active')
    ON CONFLICT (listing_id, buyer_id) DO NOTHING
    RETURNING id, listing_id, buyer_id, created_at, last_message_at, buyer_last_read_at, seller_last_read_at
)
SELECT id, listing_id, buyer_id, created_at, last_message_at, buyer_last_read_at, seller_last_read_at
FROM created
UNION ALL
SELECT id, listing_id, buyer_id, created_at, last_message_at, buyer_last_read_at, seller_last_read_at
FROM chat_threads
WHERE listing_id = $2 AND buyer_id = $3
LIMIT 1;

-- name: GetThreadParticipants :one
-- The two participants (buyer, and the listing's seller) plus last_message_at,
-- for authz, first-message detection, and read marking on the message routes.
SELECT t.id, t.buyer_id, t.last_message_at, l.seller_id
FROM chat_threads t
JOIN listings l ON l.id = t.listing_id
WHERE t.id = $1;

-- name: InsertMessage :one
INSERT INTO messages (id, thread_id, sender_id, body, image_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMessageImage :one
-- A message's image plus the thread participants, for the authorized download.
-- Scoped by thread_id (from the URL) so a message id from another thread cannot
-- be fetched through this thread's route.
SELECT m.image_key, t.buyer_id, l.seller_id
FROM messages m
JOIN chat_threads t ON t.id = m.thread_id
JOIN listings l ON l.id = t.listing_id
WHERE m.id = $1 AND m.thread_id = $2;

-- name: TouchThreadOnMessage :exec
-- Bumps last_message_at and marks the sender's own side read (you have read what
-- you just sent). The handler guarantees the sender is a participant.
UPDATE chat_threads
SET last_message_at = now(),
    buyer_last_read_at  = CASE WHEN buyer_id  = $2 THEN now() ELSE buyer_last_read_at  END,
    seller_last_read_at = CASE WHEN buyer_id <> $2 THEN now() ELSE seller_last_read_at END
WHERE id = $1;

-- name: MarkThreadRead :exec
-- Advances the reader's last_read to the newest message they fetched. The
-- handler guarantees the reader is a participant.
UPDATE chat_threads
SET buyer_last_read_at  = CASE WHEN buyer_id  = sqlc.arg('reader') THEN sqlc.arg('read_at') ELSE buyer_last_read_at  END,
    seller_last_read_at = CASE WHEN buyer_id <> sqlc.arg('reader') THEN sqlc.arg('read_at') ELSE seller_last_read_at END
WHERE id = sqlc.arg('id');

-- name: ListMessages :many
-- Cursor-poll: ascending by id (UUIDv7, time-ordered), only newer than the
-- caller's cursor. Covered by messages_thread_idx (thread_id, id).
SELECT * FROM messages
WHERE thread_id = $1
  AND (sqlc.narg('cursor')::uuid IS NULL OR id > sqlc.narg('cursor'))
ORDER BY id
LIMIT sqlc.arg('page_size');

-- name: ListExpiredMessageImageKeys :many
-- Image keys of messages whose race finished before the cutoff, so the objects
-- can be removed from storage before DeleteExpiredMessages drops the rows.
SELECT m.image_key
FROM messages m
JOIN chat_threads t ON t.id = m.thread_id
JOIN listings l ON l.id = t.listing_id
JOIN races r ON r.id = l.race_id
WHERE r.event_date < sqlc.arg('cutoff') AND m.image_key IS NOT NULL;

-- name: DeleteExpiredMessages :execrows
-- Retention: delete messages whose race finished before the cutoff (12 months
-- after race event_date). Returns the number deleted.
DELETE FROM messages m
USING chat_threads t, listings l, races r
WHERE m.thread_id = t.id
  AND t.listing_id = l.id
  AND l.race_id = r.id
  AND r.event_date < sqlc.arg('cutoff');

-- name: GetPolicyAck :one
SELECT * FROM policy_acks WHERE user_id = $1 AND race_id = $2;

-- name: CreatePolicyAck :exec
-- Idempotent: one ack per (user, race). Re-acking is a no-op, not an error.
INSERT INTO policy_acks (id, user_id, race_id, policy)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, race_id) DO NOTHING;

-- name: ListThreadsForUser :many
-- Inbox: every thread where the caller is the buyer or the listing's seller,
-- with race context, both party names (the handler picks the "other" one), and
-- the caller's unread count. The unread subquery is bounded by a user's thread
-- count; revisit with a denormalized counter if an inbox grows large.
SELECT
    t.id, t.listing_id, t.buyer_id, t.last_message_at,
    l.seller_id,
    r.name AS race_name, r.slug AS race_slug,
    bu.display_name AS buyer_name,
    su.display_name AS seller_name,
    (SELECT count(*) FROM messages m
       WHERE m.thread_id = t.id
         AND m.sender_id <> $1
         AND (
           CASE WHEN t.buyer_id = $1 THEN t.buyer_last_read_at ELSE t.seller_last_read_at END IS NULL
           OR m.created_at > CASE WHEN t.buyer_id = $1 THEN t.buyer_last_read_at ELSE t.seller_last_read_at END
         )) AS unread_count
FROM chat_threads t
JOIN listings l ON l.id = t.listing_id
JOIN races r ON r.id = l.race_id
JOIN users bu ON bu.id = t.buyer_id
JOIN users su ON su.id = l.seller_id
WHERE t.buyer_id = $1 OR l.seller_id = $1
ORDER BY t.last_message_at DESC NULLS LAST;
