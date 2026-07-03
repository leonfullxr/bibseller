-- +goose Up
-- Query-level observability (#136): the SCALING.md trigger table attaches
-- actions to measurable signals; pg_stat_statements is the zero-new-
-- infrastructure way to measure them. Collection needs the library preloaded
-- (deploy/compose.prod.yml db command); in dev/CI it is not preloaded - the
-- extension still creates (and drops) fine, but SELECTing from the view
-- errors with "pg_stat_statements must be loaded via shared_preload_libraries"
-- until the server is restarted with the preload. Nothing in the app or its
-- tests queries the view; it is read manually via psql (docs/SCALING.md).
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- +goose Down
DROP EXTENSION IF EXISTS pg_stat_statements;
