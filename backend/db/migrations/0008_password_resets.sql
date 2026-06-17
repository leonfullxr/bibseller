-- +goose Up
-- Password reset tokens. Same shape and threat model as email_verifications:
-- the table holds only the SHA-256 of the token mailed to the user, so a DB
-- leak yields no working links. Shorter-lived than verification (1 hour, set
-- in the query) because a reset hands over full control of the account.
CREATE TABLE password_resets (
    token_hash bytea PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX password_resets_user_idx ON password_resets (user_id);

-- +goose Down
DROP TABLE password_resets;
