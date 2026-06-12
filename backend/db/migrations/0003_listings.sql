-- +goose Up
CREATE TABLE listings (
    id                   uuid PRIMARY KEY,
    race_id              uuid NOT NULL REFERENCES races (id),
    seller_id            uuid NOT NULL REFERENCES users (id),
    status               text NOT NULL DEFAULT 'active'
        CONSTRAINT listings_status_check
        CHECK (status IN ('active', 'reserved', 'sold', 'cancelled', 'expired', 'removed')),
    price_cents          integer
        CONSTRAINT listings_price_nonnegative CHECK (price_cents >= 0),
    currency             char(3) NOT NULL DEFAULT 'EUR',
    original_price_cents integer
        CONSTRAINT listings_original_price_nonnegative CHECK (original_price_cents >= 0),
    description          text,
    image_key            text,
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now(),
    expires_at           timestamptz NOT NULL
);

CREATE INDEX listings_race_active_idx ON listings (race_id) WHERE status = 'active';
CREATE INDEX listings_seller_idx ON listings (seller_id);

-- +goose Down
DROP TABLE listings;
