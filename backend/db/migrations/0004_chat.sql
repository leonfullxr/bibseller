-- +goose Up
CREATE TABLE chat_threads (
    id                  uuid PRIMARY KEY,
    listing_id          uuid NOT NULL REFERENCES listings (id),
    buyer_id            uuid NOT NULL REFERENCES users (id),
    created_at          timestamptz NOT NULL DEFAULT now(),
    last_message_at     timestamptz,
    buyer_last_read_at  timestamptz,
    seller_last_read_at timestamptz,
    CONSTRAINT chat_threads_one_per_buyer UNIQUE (listing_id, buyer_id)
);

CREATE TABLE messages (
    id         uuid PRIMARY KEY, -- uuidv7: time-ordered, doubles as polling cursor
    thread_id  uuid NOT NULL REFERENCES chat_threads (id),
    sender_id  uuid NOT NULL REFERENCES users (id),
    body       text NOT NULL
        CONSTRAINT messages_body_length CHECK (char_length(body) BETWEEN 1 AND 4000),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX messages_thread_idx ON messages (thread_id, id);

-- Timestamped evidence that a buyer accepted the venue-only terms.
CREATE TABLE policy_acks (
    id       uuid PRIMARY KEY,
    user_id  uuid NOT NULL REFERENCES users (id),
    race_id  uuid NOT NULL REFERENCES races (id),
    policy   text NOT NULL,
    acked_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT policy_acks_once_per_race UNIQUE (user_id, race_id)
);

-- +goose Down
DROP TABLE policy_acks;
DROP TABLE messages;
DROP TABLE chat_threads;
