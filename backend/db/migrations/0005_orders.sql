-- +goose Up
-- Orders exist only for platform_sale races (enforced in the service layer
-- and by checkout tests; see docs/DATA_MODEL.md → order state machine).
CREATE TABLE orders (
    id                       uuid PRIMARY KEY,
    listing_id               uuid NOT NULL REFERENCES listings (id),
    buyer_id                 uuid NOT NULL REFERENCES users (id),
    seller_id                uuid NOT NULL REFERENCES users (id),
    state                    text NOT NULL DEFAULT 'pending_payment'
        CONSTRAINT orders_state_check
        CHECK (state IN ('pending_payment', 'paid_held', 'seller_marked_transferred',
                         'completed', 'disputed', 'cancelled', 'refunded')),
    item_amount_cents        integer NOT NULL
        CONSTRAINT orders_item_amount_positive CHECK (item_amount_cents > 0),
    processing_fee_cents     integer NOT NULL DEFAULT 0
        CONSTRAINT orders_fee_nonnegative CHECK (processing_fee_cents >= 0),
    total_amount_cents       integer NOT NULL,
    currency                 char(3) NOT NULL DEFAULT 'EUR',
    stripe_payment_intent_id text UNIQUE,
    stripe_transfer_id       text,
    reserved_until           timestamptz,
    seller_marked_at         timestamptz,
    completed_at             timestamptz,
    created_at               timestamptz NOT NULL DEFAULT now(),
    updated_at               timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT orders_total_consistent
        CHECK (total_amount_cents = item_amount_cents + processing_fee_cents)
);

-- At most one live order per listing.
CREATE UNIQUE INDEX orders_one_live_per_listing ON orders (listing_id)
    WHERE state IN ('pending_payment', 'paid_held', 'seller_marked_transferred', 'disputed');
CREATE INDEX orders_state_idx ON orders (state);
CREATE INDEX orders_buyer_idx ON orders (buyer_id);
CREATE INDEX orders_seller_idx ON orders (seller_id);

-- Append-only audit trail: every transition writes here in the same tx.
CREATE TABLE order_events (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id        uuid NOT NULL REFERENCES orders (id),
    from_state      text,
    to_state        text NOT NULL,
    actor           text NOT NULL,
    reason          text,
    stripe_event_id text,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX order_events_order_idx ON order_events (order_id);

-- Stripe webhook idempotency.
CREATE TABLE stripe_events (
    id           text PRIMARY KEY,
    type         text NOT NULL,
    payload      jsonb NOT NULL,
    processed_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE stripe_events;
DROP TABLE order_events;
DROP TABLE orders;
