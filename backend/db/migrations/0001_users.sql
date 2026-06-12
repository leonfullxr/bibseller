-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id                 uuid PRIMARY KEY,
    email              citext UNIQUE NOT NULL,
    email_verified_at  timestamptz,
    password_hash      text NOT NULL,
    display_name       text NOT NULL,
    locale             text NOT NULL DEFAULT 'en',
    country            char(2),
    role               text NOT NULL DEFAULT 'user'
        CONSTRAINT users_role_check CHECK (role IN ('user', 'admin')),
    stripe_account_id  text,
    stripe_customer_id text,
    anonymized_at      timestamptz,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    token_hash   bytea PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL,
    ip           inet,
    user_agent   text
);

CREATE INDEX sessions_user_idx ON sessions (user_id);

-- +goose Down
DROP TABLE sessions;
DROP TABLE users;
DROP EXTENSION IF EXISTS citext;
