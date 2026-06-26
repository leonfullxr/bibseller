# Deployment

Two supported ways to run Bibseller in production. The application, its
configuration, and the Cloudflare front are identical in both; only the machine
and the ingress differ. Keeping them the same is deliberate: moving from one to
the other is a data copy plus a DNS change, no app changes.

- Self-host (the current beta plan, zero hosting cost per D8): the founder's own
  machine reached through a Cloudflare Tunnel. The only recurring cost is the
  domain.
- Managed prod (the migration target - the "real", always-on site): a Hetzner
  VPS in the EU behind Cloudflare.

Status: design reference for M9 (#12). The runnable artifacts (a prod compose
file, a Caddyfile, the cloudflared config) are M9 deliverables; the snippets
here are ready to lift into them. Placeholders are written as `<domain>`.

## Environments and the release workflow

Three environments; two of them are the same compose stack. A release moves one
tested commit forward - you never ship something staging did not run.

| Environment | What runs it | Branch | Reached at |
|---|---|---|---|
| dev | `make dev` - Vite + hot reload + Mailpit, on localhost | feature branch | `localhost:5173` |
| staging | the prod compose stack, ephemeral, project `bibseller-staging` | `main` | `https://<staging-host>` |
| prod | the prod compose stack, always-on, project `bibseller-prod` | `production` | `https://<domain>` |

Staging is not a separate configuration: it is `deploy/compose.prod.yml` and the
same `Caddyfile`, isolated by a different compose project name (its own network
and `pgdata` / `miniodata` volumes) plus its own `deploy/.env.staging`. So
"green on staging" means the production build, behind Cloudflare, with the real
`__Host-session` cookie and `cf-ipcountry` header.

Why localhost is not enough for the final check: the session cookie is
`__Host-session` with `Secure`, so the browser rejects it without HTTPS; the
locale banner reads Cloudflare's `cf-ipcountry` (D18); and `PUBLIC_ORIGIN` is
baked into the web image (D24). Only a Cloudflare-fronted deploy exercises all
three, which is exactly what staging is for.

### One box, two stacks

Both stacks run on the self-host machine. Three things keep them apart:

- Project name: `make staging-*` passes `-p bibseller-staging`, so Docker
  namespaces its containers, network, and volumes. Prod's data lives on separate
  volumes and cannot be touched.
- Postgres host port: prod publishes `127.0.0.1:5432`, staging sets
  `DB_HOST_PORT=5433`. Nothing else publishes a host port (web, api, minio, and
  caddy are reached only through the tunnel), so that is the only possible clash.
- Resource limits: every service in `compose.prod.yml` sets
  `deploy.resources.limits` (memory + CPU), so a staging migration, seed, or
  runaway query physically cannot starve prod - without this a staging OOM could
  make the kernel OOM-killer take down prod's Postgres. Memory is the hard
  boundary. Budget for a ~16G box: per-stack steady state is ~4.9G (db 2G,
  minio 1G, api + web 768M each, caddy/mail/cloudflared 128M each), so two stacks
  are ~10G, leaving ~6G for the OS, page cache, and the transient migrate (1G) /
  minio-init one-shots. Postgres is tuned (`shared_buffers`, `work_mem`,
  `max_connections=50`) to stay inside its 2G limit. Adjust the numbers in
  `compose.prod.yml` to your box's RAM and core count.

Each stack gets its own Cloudflare Tunnel and token: create a second tunnel,
route `<staging-host>` to `http://caddy:80`, and put that token in
`deploy/.env.staging`. The two ingresses stay independent, and because `__Host-`
cookies are pinned to one host, a staging login never reaches prod.

Pick `<staging-host>` as a single label under your zone apex (for example
`bibsellertest.leonfuller.com`), not `test.<app-host>`. Cloudflare's free
Universal SSL issues only `apex` and `*.apex`, which covers a one-label host but
not a two-label one like `test.bibseller.leonfuller.com` - that host gets no edge
certificate and every request fails the TLS handshake. A deeper wildcard needs
the paid Advanced Certificate Manager, so the free path is a single-label host,
which reuses the certificate prod already has with no provisioning wait.

Staging is ephemeral on purpose: bring it up to test, tear it down after. It
keeps no data worth saving (its own volumes, the preview seed), so idle cost is
zero.

### One box: a git worktree per branch

The stack builds from the checked-out source, and one clone cannot have both
branches checked out at once. Use a git worktree so each environment builds its
own branch:

```
git checkout production                       # the existing clone -> prod
git worktree add ../bibseller-staging main    # a second dir -> staging
```

Run `make prod-*` from the prod clone and `make staging-*` from
`../bibseller-staging`.

### Test, then deploy

```
# 1. merge: feature branch -> PR -> CI green -> squash-merge to main

# 2. test the release on staging (in ../bibseller-staging, checked out on main):
git pull
make staging-up && make staging-migrate && make staging-seed
#    click through it, run a smoke check, rehearse the migration (below)
make staging-down

# 3. promote the SAME commit to prod:
git checkout production && git merge --ff-only main && git push
#    then on the box, in the prod clone (checked out on production):
git pull && make prod-migrate && make prod-up
```

`merge --ff-only` guarantees `production` is a commit `main` already carried,
never a fresh merge. If you want a deploy record, open a PR from `main` into
`production` instead of the fast-forward push: same commit, plus a CI run.

### Migration rehearsal (before a risky migration)

goose migrations are forward-only and run by hand against the live DB, so the
main deploy risk is a migration that fails on real-shaped data. Rehearse it on a
fresh staging DB loaded from a prod dump:

```
make prod-backup                                 # writes backups/db-<ts>.sql.gz
make staging-down && docker volume rm bibseller-staging_pgdata
make staging-up                                  # an empty staging DB
gunzip -c backups/db-<ts>.sql.gz | \
  docker exec -i bibseller-staging-db-1 sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
make staging-migrate                             # the new migration, on prod-shaped data
```

If it succeeds here, `make prod-migrate` on the live box is safe.

## The app, and the one rule that shapes deployment

Four processes plus a mail relay:

| Process | Port | Notes |
|---|---|---|
| web (SvelteKit, `adapter-node`) | 3000 | `npm run build` then `node build/index.js`; served at the domain root (`paths.relative: false`) |
| api (Go) | 8080 | `go build ./cmd/api` |
| Postgres 16 | 5432 | persisted volume; `goose` migrations |
| object storage (S3-compatible) | 9000 | MinIO self-hosted, or Cloudflare R2 / Scaleway managed (D16) |
| SMTP relay | 25 | Postfix (self-host) or a managed relay |

The rule: the browser calls the API with relative `/api/*` URLs (same-origin, no
CORS anywhere - see ARCHITECTURE.md). In dev, Vite proxies `/api` to the Go
server. Production has no Vite, so a reverse proxy must do that routing:
`/api/*` to the api (:8080), everything else to the web server (:3000).
Cloudflare sits in front of that proxy for three things:

- TLS termination, so the browser is on HTTPS. The session cookie is
  `__Host-session` with `Secure` always on, so production simply does not work
  over plain HTTP - the browser would reject the cookie. Cloudflare provides the
  HTTPS leg; the local hops can stay HTTP.
- CDN for the public, cacheable catalog pages (anonymous responses are
  `public, max-age=60`; signed-in ones are `private, no-store`).
- The `cf-ipcountry` header that the locale-suggestion banner reads (D18).

## Shared building blocks

### Reverse proxy (Caddy)

Required in both models. Cloudflare terminates TLS, so Caddy only speaks plain
HTTP locally and just routes by path. It forwards headers (including
`CF-IPCountry` and `X-Forwarded-*`) unchanged by default.

```
# Caddyfile - listens on :80, fronted by Cloudflare
:80 {
	handle /api/* {
		reverse_proxy api:8080
	}
	handle {
		reverse_proxy web:3000
	}
}
```

### Web server (adapter-node)

```
cd frontend && npm ci && npm run build   # -> build/
node build/index.js
```

Environment for the node server:

- `PORT=3000`
- `ORIGIN=https://<domain>` - so `url.origin` and SvelteKit's CSRF check are
  correct behind the proxy/tunnel (the server never sees the public host
  otherwise).
- `API_URL=http://api:8080` - where server-side `load` functions reach the Go
  API internally (not the public URL).

### Environment (the Go API)

Defaults live in `backend/internal/platform/config/config.go`; override per
deploy:

| Var | Production value |
|---|---|
| `ENV` | `production` |
| `PORT` | `8080` |
| `DATABASE_URL` | `postgres://bibseller:<pw>@db:5432/bibseller?sslmode=disable` (TLS terminates inside the host/network) |
| `SMTP_ADDR` | `localhost:25` (Postfix) or the relay's `host:port` |
| `EMAIL_FROM` | `Bibseller <noreply@<domain>>` |
| `APP_URL` | `https://<domain>` - used to build verification / reset links |
| `S3_ENDPOINT` | `http://minio:9000`, or the R2 / Scaleway endpoint |
| `S3_ACCESS_KEY` / `S3_SECRET_KEY` | object-storage credentials |
| `S3_BUCKET` | `bibseller` |

### Postgres

A `postgres:16` with a persisted volume. Create the role/db, then run the
migrations (`make migrate`, i.e. goose up). Back it up nightly with `pg_dump`
and copy the dump off the machine. On self-host this is not optional: the
machine is a single point of failure.

### Object storage

S3-compatible, same `minio-go` client either way (D16). Self-host: the MinIO
image from the dev compose with a persisted volume and a created bucket.
Managed: point `S3_ENDPOINT` / keys / bucket at Cloudflare R2 or Scaleway
Object Storage - no code change.

### Email and deliverability (read before self-hosting mail)

Transactional mail must actually arrive: a user cannot verify their account if
the email lands in spam. A laptop or any residential IP sending mail
direct-to-recipient-MX is a losing battle - residential ranges are on
blocklists and most ISPs block outbound port 25. So do not deliver mail
directly from the self-hosted box.

Instead run Postfix as a send-only null client that relays through an
authenticated smarthost (Brevo or SES). Postfix gives local queueing and retry
(useful when a laptop or home connection blips) and a single `SMTP_ADDR` for
the app; the smarthost does the actual delivery and signs DKIM. Set SPF, DKIM,
and DMARC records on `<domain>` (the smarthost documents the exact values).

```
# /etc/postfix/main.cf (relay-only sketch)
relayhost = [smtp-relay.brevo.com]:587
smtp_sasl_auth_enable = yes
smtp_sasl_password_maps = static:<relay-user>:<relay-key>
smtp_sasl_security_options = noanonymous
smtp_tls_security_level = encrypt
inet_interfaces = loopback-only   # accept only local mail from the app
```

(The app could also point `SMTP_ADDR` straight at the smarthost and skip
Postfix; Postfix earns its place mainly through queueing on a flaky/offline
box.)

## Model A - Self-host (laptop + Cloudflare Tunnel)  [current beta plan]

No public IP and no port forwarding. `cloudflared` makes an outbound connection
to Cloudflare's edge; Cloudflare routes the public hostname down the tunnel to
local Caddy.

```
browser
  -> Cloudflare edge        (TLS, CDN, sets CF-IPCountry)
  -> cloudflared tunnel     (outbound from the laptop; no inbound ports)
  -> Caddy :80              (/api/* -> api:8080, else -> web:3000)
  -> api -> Postgres, MinIO
     api -> Postfix(:25) -> smarthost -> recipient
```

Setup:

1. Put `<domain>` on Cloudflare (Cloudflare must manage its DNS - a Tunnel
   requirement).
2. Install `cloudflared`, authenticate, create a named tunnel, and route the
   hostname to local Caddy:
   ```
   cloudflared tunnel login
   cloudflared tunnel create bibseller
   # ~/.cloudflared/config.yml
   tunnel: bibseller
   credentials-file: /home/<user>/.cloudflared/<tunnel-id>.json
   ingress:
     - hostname: <domain>
       service: http://localhost:80
     - service: http_status:404
   cloudflared tunnel route dns bibseller <domain>
   ```
3. Bring up the app (Postgres, MinIO, api, web, Caddy) via a prod compose file,
   and run `cloudflared` as a service (or a container in the same compose).
4. Set `ENV=production`, `ORIGIN` / `APP_URL=https://<domain>`, and the secrets
   above.
5. Configure the Postfix relay; add SPF/DKIM/DMARC.
6. Backups: a cron `pg_dump` and an `mc mirror` of the MinIO bucket, both copied
   off the laptop (external drive or a cloud bucket).

Caveats, stated honestly:

- Availability: a laptop sleeps, reboots, and rides a home connection, so expect
  downtime. Disable sleep/suspend; `cloudflared` reconnects automatically when
  the box is back; the Postfix queue survives restarts. Fine for a small beta,
  not for a launch - graduate to an always-on mini-PC, or to Model B, when
  uptime starts to matter.
- Single point of failure: the laptop. Offsite backups are mandatory, not nice
  to have.
- `cf-ipcountry` still works, because Cloudflare is in front via the tunnel.

## Model B - Managed prod (Hetzner + Cloudflare)  [migration target]

```
browser
  -> Cloudflare (proxied DNS / orange cloud: TLS, CDN, CF-IPCountry)
  -> Hetzner VPS public IP
  -> Caddy (TLS via a Cloudflare Origin certificate) -> /api, else
  -> api -> Postgres; api -> R2/Scaleway; api -> Brevo/SES
```

Setup differs from self-host only at the edges:

1. A Hetzner CPX VPS in an EU region (Falkenstein or Nuremberg), Debian, with
   the firewall open only on 80/443 (and SSH).
2. Cloudflare: an `A` record to the VPS IP, proxied; TLS mode Full (strict)
   with a Cloudflare Origin certificate installed on Caddy. (You can instead run
   `cloudflared` here too, exactly as in Model A, and skip the public IP.)
3. The same app services via the prod compose file; managed object storage
   (R2 or Scaleway) and a managed email relay (Brevo or SES) for reliable
   deliverability.
4. `ENV=production`, `ORIGIN` / `APP_URL`, secrets - identical to Model A.
5. Observability (centralized logs, basic metrics) and nightly backups with a
   rehearsed restore drill. These are the M9-full deliverables.

What you gain over self-host: always-on, a public IP, managed-relay email
deliverability, and managed/automated backups, for roughly EUR 10-20/month.

## Migrating self-host -> managed

Because both run the same app on the same compose behind the same Cloudflare
account, migration is: stand up Model B, `pg_dump | psql` the database, `mc
mirror` the object-storage bucket, then repoint the Cloudflare hostname (DNS
record, or move the tunnel) at the VPS. No application changes.

## Backups and restore

The self-host box is a single point of failure (D20), so backups must leave the
box and must be proven to restore. Two scripts cover this (D26); both read
`deploy/.env.prod` and need only docker + curl on the host.

### What runs

`scripts/offsite-backup.sh` (the nightly job, `make prod-backup-offsite`):

1. `pg_dump --no-owner | gzip` to `backups/db-<timestamp>.sql.gz` (the password
   stays in the container env, never on the host argv - D23; `--no-owner` lets
   the dump restore into any role).
2. Copies that dump to Cloudflare R2 under `db/`, via the pinned `minio/mc`
   image (nothing to install on the host).
3. Mirrors the MinIO bucket to R2 under `media/`. The mirror is append-only (no
   `--remove`): a deletion or wipe in MinIO can never propagate and erase the
   offsite copy. Enable R2 object versioning on the bucket for a further layer.
4. Prunes local dumps older than 14 days. Offsite retention is an R2 lifecycle
   rule (e.g. expire `db/` after 30 days) - set it in the R2 dashboard.

On any failure it emails `BACKUP_ALERT_EMAIL` straight through the Brevo
smarthost (not via the Postfix container, so the alert still sends when the
stack itself is down) and exits non-zero.

### Schedule it (cron)

The box does not sleep (Model A caveats), so plain cron suffices. As the deploy
user (the one in the `docker` group), `crontab -e`:

```
0 3 * * * /home/<user>/bibseller/scripts/offsite-backup.sh >> /home/<user>/bibseller-backup.log 2>&1
```

The script lives on the `production` branch, so it is present in the prod clone
after the first promotion. If the box ever does sleep, use a systemd user timer
with `Persistent=true` instead, so a missed run fires on the next wake.

### The restore drill (monthly)

A backup you have never restored is not a backup. `scripts/restore-drill.sh`
(`make prod-restore-drill`) proves it:

1. Pulls the latest dump from R2 (proving the offsite artifact, not just a local
   copy).
2. Restores it into a throwaway `postgres:16` container - prod and staging are
   untouched.
3. Times the restore and asserts the schema restored (goose version present) and
   the key tables are queryable (`races`, `users`).
4. Tears the container down; exits non-zero if any assertion fails.

```
make prod-restore-drill                            # latest from R2
./scripts/restore-drill.sh backups/db-<ts>.sql.gz  # or a specific local dump
```

Note the date and restore time each month; a sudden jump in restore time is an
early warning. PITR with wal-g is the eventual upgrade (deferred) - this
dump-and-drill loop is the floor that must exist first.

### Restoring prod for real

If you ever need to bring the live database back from a dump:

```
./scripts/restore-drill.sh backups/db-<ts>.sql.gz   # confirm the dump is good first
make prod-down
docker volume rm bibseller-prod_pgdata              # discard the damaged data
make prod-up                                        # fresh empty DB
gunzip -c backups/db-<ts>.sql.gz \
  | docker exec -i bibseller-prod-db-1 sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
```

Media is restored by mirroring R2 back the other way (`mc mirror r2/<bucket>/media
src/<bucket>`).
