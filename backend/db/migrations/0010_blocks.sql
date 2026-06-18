-- +goose Up
-- One user blocking another. A block silences the conversation both ways
-- (enforced in the chat service), so the row direction records who initiated it.
CREATE TABLE blocks (
    blocker_id uuid NOT NULL REFERENCES users (id),
    blocked_id uuid NOT NULL REFERENCES users (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (blocker_id, blocked_id),
    CONSTRAINT blocks_not_self CHECK (blocker_id <> blocked_id)
);

CREATE INDEX blocks_blocked_idx ON blocks (blocked_id);

-- +goose Down
DROP TABLE blocks;
