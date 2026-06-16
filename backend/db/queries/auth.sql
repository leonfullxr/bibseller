-- Session lifecycle. The session TTL (30-day idle expiry, docs/ARCHITECTURE.md
-- → Auth & sessions) lives in this file: both INSERT and touch use the same
-- interval literal.

-- name: CreateSession :one
INSERT INTO sessions (token_hash, user_id, expires_at, ip, user_agent)
VALUES ($1, $2, now() + interval '30 days', $3, $4)
RETURNING expires_at;

-- name: GetSessionWithUser :one
SELECT s.last_seen_at, u.*
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1
  AND s.expires_at > now();

-- name: TouchSession :exec
UPDATE sessions
SET last_seen_at = now(), expires_at = now() + interval '30 days'
WHERE token_hash = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token_hash = $1;
