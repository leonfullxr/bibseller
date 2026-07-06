-- Drop the five schema-v1 columns no code path reads or meaningfully writes
-- (#124). Each returns with the feature that actually uses it: role -> first
-- admin-gated endpoint, stripe_* -> M6, anonymized_at -> M7 delete flow,
-- listings.image_key -> never (D15 moved photos to messages.image_key).

-- +goose Up
ALTER TABLE users
    DROP COLUMN role,
    DROP COLUMN stripe_account_id,
    DROP COLUMN stripe_customer_id,
    DROP COLUMN anonymized_at;

ALTER TABLE listings
    DROP COLUMN image_key;

-- +goose Down
ALTER TABLE users
    ADD COLUMN role text NOT NULL DEFAULT 'user'
        CONSTRAINT users_role_check CHECK (role IN ('user', 'admin')),
    ADD COLUMN stripe_account_id text,
    ADD COLUMN stripe_customer_id text,
    ADD COLUMN anonymized_at timestamptz;

ALTER TABLE listings
    ADD COLUMN image_key text;
