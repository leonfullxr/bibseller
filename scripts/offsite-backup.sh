#!/usr/bin/env bash
# Nightly offsite backup for the self-host prod stack (docs/DEPLOYMENT.md).
# pg_dump + an append-only MinIO mirror, both to Cloudflare R2; prunes old local
# dumps; emails via the Brevo relay on any failure. Run by cron, safe by hand.
# Needs only docker + curl on the host.
#
# Reads deploy/.env.prod with grep/cut, never `source`: a password with @ & %
# would break shell-sourcing (CONTEXT D23).
set -euo pipefail

cd "$(dirname "$0")/.."                          # repo root (the prod clone)
ENV_FILE=deploy/.env.prod
COMPOSE="docker compose --env-file $ENV_FILE -f deploy/compose.prod.yml"
MC_IMAGE=minio/mc:RELEASE.2025-08-13T08-35-41Z   # same pin as compose.prod.yml
NETWORK=bibseller-prod_default                   # prod compose default network
KEEP_DAYS=14                                     # local dump retention
STAMP=$(date +%F-%H%M%S)
DUMP="backups/db-$STAMP.sql.gz"

getenv() { grep -E "^$1=" "$ENV_FILE" 2>/dev/null | head -n1 | cut -d= -f2- || true; }

# Fail fast (and quietly) on an unconfigured offsite target, before the trap is
# set - a setup mistake should print, not try to email through unset relay creds.
for k in R2_ENDPOINT R2_ACCESS_KEY R2_SECRET_KEY R2_BUCKET; do
  case "$(getenv "$k")" in ""|*replace-with*|*"<accountid>"*)
    echo "offsite-backup: $k is not set in $ENV_FILE" >&2; exit 1 ;; esac
done

alert_email=$(getenv BACKUP_ALERT_EMAIL)
notify_fail() {
  trap - ERR
  echo "BACKUP FAILED: $1" >&2
  if [ -n "$alert_email" ]; then
    local relay from
    relay=$(getenv MAIL_RELAYHOST); relay=${relay//[\[\]]/}   # [host]:port -> host:port
    from=$(getenv EMAIL_FROM)
    case "$from" in *"<"*">"*) from=${from##*<}; from=${from%%>*} ;; esac
    printf 'From: %s\r\nTo: %s\r\nSubject: [bibseller] offsite backup FAILED\r\n\r\n%s (%s)\r\n' \
      "$(getenv EMAIL_FROM)" "$alert_email" "$1" "$STAMP" \
      | curl -sS --ssl-reqd --url "smtp://$relay" \
          --user "$(getenv MAIL_RELAYHOST_USERNAME):$(getenv MAIL_RELAYHOST_PASSWORD)" \
          --mail-from "$from" --mail-rcpt "$alert_email" --upload-file - \
      || echo "alert email ALSO failed" >&2
  fi
  exit 1
}
trap 'notify_fail "command failed near line $LINENO"' ERR

mkdir -p backups

# 1. Postgres dump. Password stays in the container env (compose --env-file),
#    never on the host argv (D23). --no-owner so it restores into any role.
$COMPOSE exec -T db sh -c 'pg_dump --no-owner -U "$POSTGRES_USER" "$POSTGRES_DB"' | gzip > "$DUMP"
[ -s "$DUMP" ] || notify_fail "pg_dump produced an empty file"

# 2+3. dump -> R2 (db/) and the MinIO bucket -> R2 (media/), via the pinned mc
#      image on the prod network. The mirror is append-only (no --remove): a
#      deletion or wipe in MinIO must never propagate and erase the offsite copy.
docker run --rm --network "$NETWORK" -v "$PWD/backups:/backups:ro" \
  -e SRC_KEY="$(getenv S3_ACCESS_KEY)" -e SRC_SECRET="$(getenv S3_SECRET_KEY)" \
  -e SRC_BUCKET="$(getenv S3_BUCKET)"  -e DUMP="$(basename "$DUMP")" \
  -e R2_ENDPOINT="$(getenv R2_ENDPOINT)" -e R2_KEY="$(getenv R2_ACCESS_KEY)" \
  -e R2_SECRET="$(getenv R2_SECRET_KEY)" -e R2_BUCKET="$(getenv R2_BUCKET)" \
  --entrypoint /bin/sh "$MC_IMAGE" -c '
    set -e
    mc alias set src "http://minio:9000" "$SRC_KEY" "$SRC_SECRET" >/dev/null
    mc alias set r2  "$R2_ENDPOINT"      "$R2_KEY"  "$R2_SECRET"  >/dev/null
    mc cp "/backups/$DUMP" "r2/$R2_BUCKET/db/$DUMP"
    mc mirror "src/$SRC_BUCKET" "r2/$R2_BUCKET/media"
  ' || notify_fail "mc upload/mirror to R2 failed"

# 4. prune local dumps (offsite retention is an R2 lifecycle rule - see docs).
find backups -name 'db-*.sql.gz' -mtime +"$KEEP_DAYS" -delete

echo "OK $STAMP: $DUMP -> r2/$(getenv R2_BUCKET)/db, MinIO mirrored to media/."
