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
  (gen_random_uuid(), 'preview-valencia-half-2026',   'Valencia Half Marathon 2026','running',   'half',     DATE '2026-10-25', 'Valencia', 'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-vienna-marathon-2027',     'Vienna City Marathon 2027',  'running',   'marathon',   DATE '2027-04-18', 'Vienna',    'AT', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-brussels-20k-2026',        'Brussels 20K 2026',          'running',   '20k',        DATE '2026-09-20', 'Brussels',  'BE', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-amsterdam-marathon-2026',  'Amsterdam Marathon 2026',    'running',   'marathon',   DATE '2026-10-18', 'Amsterdam', 'NL', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-rotterdam-marathon-2027',  'Rotterdam Marathon 2027',    'running',   'marathon',   DATE '2027-04-11', 'Rotterdam', 'NL', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-cracovia-marathon-2027',   'Cracovia Marathon 2027',     'running',   'marathon',   DATE '2027-04-25', 'Kraków',    'PL', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-porto-marathon-2026',      'Porto Marathon 2026',        'running',   'marathon',   DATE '2026-11-08', 'Porto',     'PT', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-lisbon-half-2027',         'Lisbon Half Marathon 2027',  'running',   'half',       DATE '2027-03-07', 'Lisbon',    'PT', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-sevilla-marathon-2027',    'Sevilla Marathon 2027',      'running',   'marathon',   DATE '2027-02-21', 'Sevilla',   'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-bilbao-night-marathon-2026','Bilbao Night Marathon 2026','running',   'marathon',   DATE '2026-10-17', 'Bilbao',    'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-munich-marathon-2026',     'Munich Marathon 2026',       'running',   'marathon',   DATE '2026-10-11', 'Munich',    'DE', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-milano-granfondo-2027',    'Milano Gran Fondo 2027',     'cycling',   'gran fondo', DATE '2027-03-21', 'Milan',     'IT', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.'),
  (gen_random_uuid(), 'preview-madrid-triathlon-2027',    'Madrid Triathlon 2027',      'triathlon', 'olympic',    DATE '2027-06-05', 'Madrid',    'ES', 'unknown', 'published', 'Preview entry - transfer policy not yet verified.')
ON CONFLICT (slug) DO NOTHING;
