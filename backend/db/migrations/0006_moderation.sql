-- +goose Up
CREATE TABLE reports (
    id             uuid PRIMARY KEY,
    reporter_id    uuid REFERENCES users (id), -- nullable: DSA notices can be anonymous
    reporter_email text,
    subject_type   text NOT NULL
        CONSTRAINT reports_subject_type_check
        CHECK (subject_type IN ('listing', 'message', 'user')),
    subject_id     uuid NOT NULL,
    reason         text NOT NULL
        CONSTRAINT reports_reason_check
        CHECK (reason IN ('forbidden_transfer', 'scam', 'offensive', 'other')),
    details        text,
    status         text NOT NULL DEFAULT 'open'
        CONSTRAINT reports_status_check CHECK (status IN ('open', 'actioned', 'dismissed')),
    created_at     timestamptz NOT NULL DEFAULT now(),
    resolved_at    timestamptz,
    resolved_by    uuid REFERENCES users (id)
);

CREATE INDEX reports_open_idx ON reports (created_at) WHERE status = 'open';

CREATE TABLE audit_log (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    admin_id     uuid NOT NULL REFERENCES users (id),
    action       text NOT NULL,
    subject_type text NOT NULL,
    subject_id   uuid,
    details      jsonb,
    created_at   timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE audit_log;
DROP TABLE reports;
