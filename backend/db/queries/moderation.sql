-- name: CreateReport :one
INSERT INTO reports (id, reporter_id, subject_type, subject_id, reason, details)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateBlock :exec
-- Idempotent: blocking an already-blocked user is a no-op.
INSERT INTO blocks (blocker_id, blocked_id)
VALUES ($1, $2)
ON CONFLICT (blocker_id, blocked_id) DO NOTHING;

-- name: DeleteBlock :exec
DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2;

-- name: IsBlocked :one
-- True if either user has blocked the other - a block silences the conversation
-- both ways.
SELECT EXISTS (
    SELECT 1 FROM blocks
    WHERE (blocker_id = $1 AND blocked_id = $2)
       OR (blocker_id = $2 AND blocked_id = $1)
);
