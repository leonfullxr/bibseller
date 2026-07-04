-- Email verification lifecycle. The 24-hour TTL lives here (both the INSERT and
-- the not-expired guard), mirroring the session TTL convention in auth.sql.

-- name: CreateEmailVerification :one
INSERT INTO email_verifications (token_hash, user_id, expires_at)
VALUES ($1, $2, now() + interval '24 hours')
RETURNING expires_at;

-- name: GetEmailVerificationUser :one
SELECT user_id FROM email_verifications
WHERE token_hash = $1
  AND expires_at > now();

-- name: DeleteEmailVerificationsForUser :exec
DELETE FROM email_verifications WHERE user_id = $1;

-- name: MarkEmailVerified :exec
UPDATE users
SET email_verified_at = now(), updated_at = now()
WHERE id = $1 AND email_verified_at IS NULL;

-- name: DeleteExpiredEmailVerificationsBatch :execrows
-- Janitor batch (#142): tokens for users who never verify are otherwise only
-- removed by account deletion (CASCADE), so expired rows accumulate. Same
-- retention buffer as sessions (expiry + 30 days); batch by primary key (#99).
DELETE FROM email_verifications
WHERE token_hash IN (
    SELECT v.token_hash FROM email_verifications v
    WHERE v.expires_at < sqlc.arg('cutoff')
    LIMIT sqlc.arg('batch_size')
);
