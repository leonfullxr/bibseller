-- +goose Up
-- ListRaces (backend/db/queries/race.sql) sorts published races by
-- (event_date, id) and keyset-paginates on the same tuple; races_browse_idx
-- only covers (country, event_date), so that sort/cursor pays for a step of
-- its own. This backs it directly (#95).
CREATE INDEX races_event_date_id_idx ON races (event_date, id) WHERE status = 'published';

-- +goose Down
DROP INDEX races_event_date_id_idx;
