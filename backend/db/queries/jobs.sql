-- Generic background-job coordination. A transaction-scoped advisory lock lets
-- every API instance run the same ticker while only one does the work per tick
-- (architecture invariant 5: stateless API, jobs coordinate via Postgres
-- advisory locks). The lock auto-releases when the transaction ends.

-- name: TryAdvisoryXactLock :one
SELECT pg_try_advisory_xact_lock(sqlc.arg('key')::bigint);
