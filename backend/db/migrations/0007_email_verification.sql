-- +goose Up
-- Email verification tokens. Like sessions, the table holds only the SHA-256 of
-- the token mailed to the user, so a DB leak yields no working links. The token
-- gates listing/chat creation from M4 (browsing stays open); a NULL
-- users.email_verified_at means "not yet verified".
CREATE TABLE email_verifications (
    token_hash bytea PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX email_verifications_user_idx ON email_verifications (user_id);

-- +goose Down
DROP TABLE email_verifications;
