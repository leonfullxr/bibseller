-- +goose Up
-- One row per sellable event (edition × distance); see docs/DATA_MODEL.md.
CREATE TABLE races (
    id                    uuid PRIMARY KEY,
    slug                  text UNIQUE NOT NULL,
    name                  text NOT NULL,
    series                text,
    sport                 text NOT NULL DEFAULT 'running'
        CONSTRAINT races_sport_check
        CHECK (sport IN ('running', 'trail', 'triathlon', 'cycling', 'obstacle', 'other')),
    distance              text,
    event_date            date NOT NULL,
    city                  text NOT NULL,
    country               char(2) NOT NULL,
    website_url           text,

    -- The load-bearing field: drives every downstream behavior.
    transfer_policy       text NOT NULL DEFAULT 'unknown'
        CONSTRAINT races_policy_check
        CHECK (transfer_policy IN ('platform_sale', 'official_only', 'connect_only', 'unknown')),
    official_transfer_url text,
    policy_source_url     text,
    policy_notes          text,
    policy_verified_at    timestamptz,
    policy_verified_by    uuid REFERENCES users (id),

    status                text NOT NULL DEFAULT 'draft'
        CONSTRAINT races_status_check CHECK (status IN ('draft', 'published', 'archived')),
    created_by            uuid REFERENCES users (id),
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    -- Policies with legal weight require their evidence.
    CONSTRAINT races_official_url_required
        CHECK (transfer_policy <> 'official_only' OR official_transfer_url IS NOT NULL),
    CONSTRAINT races_platform_source_required
        CHECK (transfer_policy <> 'platform_sale' OR policy_source_url IS NOT NULL)
);

CREATE INDEX races_browse_idx ON races (country, event_date) WHERE status = 'published';
CREATE INDEX races_fts_idx ON races
    USING gin (to_tsvector('simple', name || ' ' || city));

-- +goose Down
DROP TABLE races;
