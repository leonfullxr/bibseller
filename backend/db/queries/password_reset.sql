-- Password reset lifecycle. The 1-hour TTL lives here (both the INSERT and the
-- not-expired guard), tighter than the 24-hour verification token because a
-- reset grants full account access.

-- name: CreatePasswordReset :one
INSERT INTO password_resets (token_hash, user_id, expires_at)
VALUES ($1, $2, now() + interval '1 hour')
RETURNING expires_at;

-- name: ConsumePasswordReset :one
-- Atomically validates and consumes a token: only one of two concurrent
-- requests with the same token can delete the row and get the user_id back.
DELETE FROM password_resets
WHERE token_hash = $1
  AND expires_at > now()
RETURNING user_id;

-- name: DeletePasswordResetsForUser :exec
DELETE FROM password_resets WHERE user_id = $1;

-- name: DeleteExpiredPasswordResetsBatch :execrows
-- Janitor batch (#142): ConsumePasswordReset removes used tokens, but a
-- requested-and-never-clicked reset stays forever. Same retention buffer as
-- sessions (expiry + 30 days); batch by primary key (#99).
DELETE FROM password_resets
WHERE token_hash IN (
    SELECT p.token_hash FROM password_resets p
    WHERE p.expires_at < sqlc.arg('cutoff')
    LIMIT sqlc.arg('batch_size')
);
