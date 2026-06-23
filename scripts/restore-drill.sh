#!/usr/bin/env bash
# Restore drill: prove a backup actually restores (docs/DEPLOYMENT.md).
# Pulls the latest Postgres dump from R2 (or takes a local .sql.gz as $1),
# restores it into a throwaway postgres:16 container, times the restore, and
# asserts the schema + key tables. Prod and staging are untouched.
#
#   ./scripts/restore-drill.sh                          # latest dump from R2
#   ./scripts/restore-drill.sh backups/db-....sql.gz    # a specific local dump
set -euo pipefail

cd "$(dirname "$0")/.."
ENV_FILE=deploy/.env.prod
MC_IMAGE=minio/mc:RELEASE.2025-08-13T08-35-41Z
PG_IMAGE=postgres:16
CONTAINER=bibseller-restore-drill
WORK=$(mktemp -d)
getenv() { grep -E "^$1=" "$ENV_FILE" 2>/dev/null | head -n1 | cut -d= -f2- || true; }
cleanup() { docker rm -f "$CONTAINER" >/dev/null 2>&1 || true; rm -rf "$WORK"; }
trap cleanup EXIT

# 1. obtain a dump: explicit local file, else the newest db/*.sql.gz from R2.
if [ "${1:-}" ]; then
  DUMP="$1"; [ -s "$DUMP" ] || { echo "no such dump: $DUMP" >&2; exit 1; }
  echo "Using local dump: $DUMP"
else
  echo "Fetching the latest dump from R2..."
  docker run --rm -v "$WORK:/out" \
    -e E="$(getenv R2_ENDPOINT)" -e K="$(getenv R2_ACCESS_KEY)" \
    -e S="$(getenv R2_SECRET_KEY)" -e B="$(getenv R2_BUCKET)" \
    --entrypoint /bin/sh "$MC_IMAGE" -c '
      set -e
      mc alias set r2 "$E" "$K" "$S" >/dev/null
      latest=$(mc ls "r2/$B/db/" | awk "{print \$NF}" | sort | tail -n1)
      [ -n "$latest" ] || { echo "no dumps under r2/$B/db/" >&2; exit 1; }
      echo "latest on R2: $latest"
      mc cp "r2/$B/db/$latest" /out/dump.sql.gz
    '
  DUMP="$WORK/dump.sql.gz"
fi

# 2. throwaway postgres; wait until ready.
echo "Starting throwaway $PG_IMAGE..."
docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" -e POSTGRES_PASSWORD=drill -e POSTGRES_DB=drill "$PG_IMAGE" >/dev/null
# Wait for the REAL server: the official image first runs a socket-only bootstrap
# server during init, so check TCP (127.0.0.1) - only the final server listens there.
ready=
for i in $(seq 1 60); do
  if docker exec "$CONTAINER" pg_isready -h 127.0.0.1 -U postgres >/dev/null 2>&1; then ready=1; break; fi
  sleep 1
done
[ -n "$ready" ] || { echo "postgres did not become ready" >&2; exit 1; }
# A plain (owner-ful) dump references role "bibseller"; create it so the restore
# does not trip on ownership. --no-owner dumps simply ignore this.
docker exec "$CONTAINER" psql -U postgres -d drill -c "CREATE ROLE bibseller LOGIN" >/dev/null 2>&1 || true

# 3. restore, timed.
echo "Restoring..."
start=$(date +%s)
gunzip -c "$DUMP" | docker exec -i "$CONTAINER" psql -v ON_ERROR_STOP=1 -U postgres -d drill >/dev/null
elapsed=$(( $(date +%s) - start ))

# 4. assert: schema restored (goose) + key tables queryable.
q() { docker exec "$CONTAINER" psql -tA -U postgres -d drill -c "$1" 2>/dev/null || true; }
ver=$(q "select max(version_id) from goose_db_version")
races=$(q "select count(*) from races")
users=$(q "select count(*) from users")

echo "--------------------------------------------"
echo "dump:          $(basename "$DUMP")"
echo "restore time:  ${elapsed}s"
echo "goose version: ${ver:-<none>}"
echo "races rows:    ${races:-<error>}"
echo "users rows:    ${users:-<error>}"
echo "--------------------------------------------"
[ -n "$ver" ]              || { echo "FAIL: goose_db_version empty - schema did not restore" >&2; exit 1; }
[[ "$races" =~ ^[0-9]+$ ]] || { echo "FAIL: races table not queryable" >&2; exit 1; }
[[ "$users" =~ ^[0-9]+$ ]] || { echo "FAIL: users table not queryable" >&2; exit 1; }
echo "PASS - backup restores cleanly."
