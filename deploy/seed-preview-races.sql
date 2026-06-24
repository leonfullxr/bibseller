-- Preview/template races for the live site (not real events).
-- All transfer_policy = 'unknown': chat-only, "policy unverified" badge, no money flow.
-- Idempotent (ON CONFLICT (slug) DO NOTHING) and removable in one query:
--   DELETE FROM races WHERE slug LIKE 'preview-%';
-- Run with:
--   docker exec bibseller-prod-db-1 sh -c \
--     'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < deploy/seed-preview-races.sql
INSERT INTO races (id, slug, name, sport, distance, event_date, city, country, transfer_policy, status, policy_notes)
VALUES
  (gen_random_uuid(), 'preview-granada-half-2027',   'Granada Media Maraton 2027', 'running',   'half',     DATE '2027-03-28', 'Granada',  'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-valencia-marathon-2026','Valencia Marathon 2026',    'running',   'marathon', DATE '2026-12-06', 'Valencia', 'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-paris-marathon-2027',  'Paris Marathon 2027',       'running',   'marathon', DATE '2027-04-11', 'Paris',    'FR', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-berlin-marathon-2026', 'Berlin Marathon 2026',      'running',   'marathon', DATE '2026-09-27', 'Berlin',   'DE', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-garda-trail-42k-2026', 'Garda Trail 42K 2026',      'trail',     '42k',      DATE '2026-09-19', 'Riva del Garda', 'IT', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-mallorca-70-3-2027',   'Mallorca 70.3 Triathlon 2027','triathlon','70.3',    DATE '2027-05-08', 'Alcudia',  'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-valencia-half-2026',   'Valencia Half Marathon 2026','running',   'half',     DATE '2026-10-25', 'Valencia', 'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.')
ON CONFLICT (slug) DO NOTHING;
