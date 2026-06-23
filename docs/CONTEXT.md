# Project context & decision log

The single source of truth for product decisions and working agreements.
Read together with [PRODUCT.md](PRODUCT.md), [ARCHITECTURE.md](ARCHITECTURE.md), [DATA_MODEL.md](DATA_MODEL.md), [ROADMAP.md](ROADMAP.md), [DEPLOYMENT.md](DEPLOYMENT.md). When code or docs disagree with this file, this file wins until amended. Every founder decision lands here with its date and rationale - append, don't rewrite history.

## Mission

Race bibs go to waste while runners refresh sold-out registration pages. Bibseller connects sellers and buyers of race bibs EU-wide, peer-to-peer, at zero commission - always within each race's own transfer rules. The goal is the match, not the margin.

## The one concept everything hangs on

Every race carries a `transfer_policy` (`platform_sale | official_only | connect_only | unknown`). Every behavior - which buttons exist, whether money can flow, which disclaimers render, what liability posture applies - derives from that one field, enforced server-side and in the database, never just hidden in UI. If a change can't answer "what does this do in each of the four modes?", it isn't done.

## Language

Policy view - the frontend's single derivation of a race's `transfer_policy` into presentation facts: the CTA affordance, the disclaimer block and its tone, and the badge label. Lives in `$lib/policy.ts`; components read it instead of branching on the policy string, and the words it labels sit in a separate table (the M8 i18n seed). Avoid: policy descriptor, policy config, policy helper.

## Founder decisions log

| # | Date | Decision | Choice | Rationale / consequences |
|---|------|----------|--------|--------------------------|
| D1 | 2026-06-12 | Stripe processing fee | Buyer pays it, at cost - itemized "payment processing" line at checkout | Zero commission stays true and sustainable; platform nets ≈ €0 per transaction. Gross-up formula in [PAYMENTS_AND_COMPLIANCE.md](PAYMENTS_AND_COMPLIANCE.md#fee-reality--pass-through). Shapes M6 checkout math. |
| D2 | 2026-06-12 | Anti-scalping price cap | Hard cap at face value when the seller provides the original price; warning banner when they don't | Matches the non-profit mission and most races' own rules. Enforced in the M4 service layer + tests (cross-row context makes a DB CHECK awkward). |
| D3 | 2026-06-12 | Launch strategy | Chat-only public beta first (Gate 1), payments enabled later (Gate 2) | Validates the matching idea with the smallest legal/compliance surface. Beta = M0-M5 + M7-lite + M8 + M9. |
| D4 | 2026-06-12 | First market | Spain first, then expand | Ops catalog research (#6) prioritizes ES races; Spanish is the first translation (M8). Seed data already reflects this. |
| D5 | 2026-06-12 | PR workflow | Foundation merged as PR #14; one PR per milestone afterwards | Smaller reviews, milestone issues auto-close per PR, cleaner revert points. |
| D6 | 2026-06-12 | Verification bar | `make verify` + `make smoke` harness + per-issue acceptance checklists (Playwright deferred) | Every PR proves its issue's acceptance criteria with commands/tests. Revisit Playwright when M4 introduces real form flows. |
| D7 | 2026-06-12 | Division of labor | Claude implements, founder reviews to learn | Founder learns Go/Svelte by reviewing milestone PRs; can grab any issue at will. |
| D8 | 2026-06-12 | Deployment timing | Deferred - no walking-skeleton deploy now; M9 starts before the beta gate | Saves cost/ops now; accepts that ops surprises arrive later. Hosting choice still open (see Open questions). |
| D9 | 2026-06-11 | Commission | Zero. Non-negotiable. | Founding stance; the platform may only pass through processor fees (D1). |
| D10 | 2026-06-11 | Races that forbid transfers | Stay listed, chat-only, strong disclaimer + recorded acknowledgment (`policy_acks`) | Founder intent: provide the channel, take no part in money or transfer. Revisit with counsel during M7. |
| D11 | 2026-06-11 | Currency | EUR-only v1; `currency` column exists everywhere | Add SEK/PLN/DKK post-launch by demand. |
| D12 | 2026-06-11 | Auth | Own email+password + Postgres sessions, no auth SaaS; OAuth (Google/Strava) post-v1 | Cost zero, EU data residency, learning value. Spec in ARCHITECTURE.md. |
| D13 | 2026-06-11 | Chat transport | HTTP polling first; SSE/WebSocket upgrade is transport-only | Upgrade trigger: sustained poll QPS > ~2k or p95 > 100ms. |
| D14 | 2026-06-12 | Stdlib-first stack | Drop chi and Tailwind; no superforms/zod/i18n libraries. Routing on Go 1.22 `http.ServeMux` (method + `{wildcard}` patterns), styling in scoped Svelte `<style>` blocks with tokens in `layout.css`, forms via SvelteKit form actions + native HTML5 validation | Founder learning goal (D7): core Go/Svelte primitives over third-party abstractions. pgx/sqlc/goose stay - typed SQL without ORM bloat. |
| D15 | 2026-06-17 | Listing photos and handover confirmation | Photos are private chat artifacts (built in M5, server-proxied through the API), not public listing images. Buyer and seller confirm the handover with an in-chat acknowledgment button; a photo is optional. | Proof-of-registration photos can contain personal data, so they stay private to the buyer-seller conversation, never the public catalog. Image upload therefore moves from M4 (#7) to M5 (#8); the held-funds confirmation/escrow is the M6 order flow. M4 ships listings only. |
| D16 | 2026-06-18 | Object-storage client (first deliberate dependency under D14) | `minio-go` | In-chat images are private, server-proxied artifacts (D15): the API authorizes each fetch by thread participation, so it must hold an S3 client - presigned-PUT-only (the old M4 security-checklist line) does not fit. `minio-go` is the lightweight S3-compatible client, identical across MinIO (dev) and R2/Scaleway (prod); first justified exception to D14's no-new-dependencies stance. Lands in M5.3 (#38). |
| D17 | 2026-06-18 | i18n approach (M8) | Hand-rolled message dictionaries + English-root / `/es`-prefix routing | Reaffirms D14 (no i18n library) over issue #9's "Paraglide" title: en+es is small enough to hand-roll a typed dictionary + `t()` accessor, keeping the stdlib-first learning goal. Locale URL strategy: English at the root, Spanish under `/es` with `hreflang` alternates (best for search traffic); locale resolved by signed-in `users.locale` > `locale` cookie > `Accept-Language` > `en`. Lands in M8.1 (#45) + M8.2 (#46). |
| D18 | 2026-06-19 | Locale detection: suggest, don't force (refines D17) | Land in English; a *settled* choice (signed-in `users.locale`, else the `locale` cookie) redirects the URL to match it; soft signals never auto-redirect. A no-cookie visitor who looks Spanish by location (`cf-ipcountry`=ES, Cloudflare's prod geo header; `Accept-Language` as fallback when geo is absent) gets a dismissible "switch to Spanish" banner: accept -> cookie=es + `/es`, dismiss -> cookie=en. | Auto-redirecting by detected locale is discouraged by Google (can block crawling of the other language) and fights the English landing page the founder wants; a banner is the recommended pattern and keeps both URLs crawlable. Server-side error-string i18n shipped in #49 (API returns codes; frontend maps `t(apiError.${code})`). Lands in M8.1 (#45). |
| D19 | 2026-06-22 | Cookie posture (M7-lite) | Essential-only: the single `__Host-session` cookie; no analytics or tracking cookies, so no consent banner | Follows D12 (own auth + Postgres sessions). A strictly-necessary session cookie needs no consent under the ePrivacy Directive; stated on the Privacy page (#54). Revisit if analytics or third-party embeds are added. |
| D20 | 2026-06-22 | Beta hosting | Self-host on the founder's machine via a Cloudflare Tunnel for the public beta; Hetzner + Cloudflare is the documented always-on migration target | Extends D8 (zero hosting cost) and keeps `cf-ipcountry` (D18) working since Cloudflare fronts the tunnel. Accept the laptop as a single point of failure with downtime at beta scale; offsite backups mandatory; email relays through a smarthost (residential IPs cannot deliver direct). Both models + the migration path live in DEPLOYMENT.md; the self-host build is M9.1 (#60). |
| D21 | 2026-06-22 | Object storage image (M9.1) | Self-hosted MinIO pinned to the last published official OSS release (`RELEASE.2025-09-07`); app + dev both pinned | MinIO's OSS edition went maintenance-only (Dec 2025) then was archived; no new images, so `:latest` is stale/unpatched. Pinning gives a trusted, reproducible image for the beta - we drive MinIO via `mc`, so the stripped-down console is irrelevant. If a CVE forces a move: the `pgsty/minio` community fork (patched) or Cloudflare R2 (free tier); the minio-go client (D16) works unchanged against any S3-compatible target. Beta keeps the app on MinIO's root credentials (single-box); a bucket-scoped access key is a noted future hardening. |
| D22 | 2026-06-23 | Prod email provider (M9.1) | Brevo as the SMTP smarthost; send from a subdomain (`admin@bibseller.leonfuller.com`) authenticated with SPF + DKIM on Cloudflare DNS, DMARC inherited from the apex | Concretizes D20's "relay through a smarthost". A residential IP cannot deliver mail directly, and the apex `leonfuller.com` (Proton MX, `p=quarantine`) does not authorize Brevo, so unauthenticated sends were silently quarantined. Sending from an authenticated subdomain isolates the app's sending reputation from the founder's other mail and makes DKIM align. Postfix (`boky/postfix`) relays app -> Brevo; the app -> Postfix hop is plaintext on the internal network (no STARTTLS, #65). |
| D23 | 2026-06-23 | Prod DB connection string (M9.1) | libpq keyword/value DSN (not a `postgres://` URL); migrations run via a one-shot `migrate` compose container, not host shell-sourcing | A random `POSTGRES_PASSWORD` with `@ & %` corrupts a URL DSN and breaks `set -a; . .env.prod` shell-sourcing (#64). The keyword/value form needs no URL-escaping (pgx parses it) and routing the password through compose interpolation keeps it out of any host shell. The api also stopped logging the raw DSN on a parse failure, so a malformed value cannot leak the password (#63). |
| D24 | 2026-06-23 | Prod web asset URLs (M9.1) | Bake `PUBLIC_ORIGIN` into the web image at build time via `kit.paths.assets`; the image is therefore origin-specific | With `paths.relative:false` (kept so the i18n `link()` helper's `/es` prefix works, D17), Vite's `__vitePreload` emitted root-relative dep strings that resolved against the entry-module URL and doubled the path (`/_app/immutable/entry/_app/...`) -> 404 -> no CSS/JS (#66). Setting `paths.assets` to the origin emits absolute asset URLs (and importer-relative dynamic imports), fixing it without touching the ~57 `link()` call sites. Trade-off: rebuild the web image if the domain changes. Replaces the Caddy rewrite workaround. |
| D25 | 2026-06-23 | Test/release workflow (post-M9.1) | Staging is the prod compose stack run as a second, ephemeral project (`bibseller-staging`) on the same box behind a single-label staging host (e.g. `bibsellertest.leonfuller.com`); `main` runs staging, a long-lived `production` branch runs prod, and a release is a `--ff-only` of the tested commit. Promote by fast-forwarding `production` to `main` (or a PR `main` -> `production`). | Localhost cannot exercise the `__Host-session` cookie (needs HTTPS), `cf-ipcountry` (D18), or the baked `PUBLIC_ORIGIN` (D24), so the final pre-deploy check must be a Cloudflare-fronted deploy. Reusing one compose file + Caddyfile (isolated by compose project name + `DB_HOST_PORT`) avoids a second config to maintain; ephemeral staging keeps idle cost zero on the single self-host box (D20). Box uses a git worktree per branch. The staging host must be a single label under the zone apex: Cloudflare free Universal SSL covers apex + `*.apex` only, so a two-label host like `test.bibseller.leonfuller.com` has no edge certificate and fails the TLS handshake (deeper wildcards need the paid Advanced Certificate Manager). No CD pipeline yet - deploys stay `make` over SSH. Full flow in DEPLOYMENT.md. |

## Scope boundaries

In for beta (Gate 1): race catalog with verified policies, accounts, listings with photos, buyer-seller chat with policy acknowledgments, report/block, legal pages, GDPR export/delete, en+es, EU deployment.

In for Gate 2: Stripe Connect payments for `platform_sale` races only, order state machine with held funds, full DSA moderation queue, DAC7 memo signed off.

Out (parked): native apps, ads/premium tiers, race-organizer SaaS, multi-currency, OAuth, WebSocket chat, search engine beyond Postgres FTS, anonymous contact.

## Architecture invariants (always true, test-enforced where possible)

1. No payment code path exists outside `platform_sale` - checkout rejects other policies at the API; there is no PaymentIntent construction reachable for them (M6 tests make this structural).
2. The database enforces the load-bearing rules - policy enum, evidence URLs required per policy, order-state CHECK, `total = item + fee`, one live order per listing. `backend/internal/platform/db/schema_test.go` proves each.
3. Money is integer cents + currency code. Never floats.
4. Order state transitions happen only in `internal/order` via guarded `UPDATE … WHERE state = $from`, appending `order_events` in the same transaction. Terminal states are frozen.
5. Stateless API - any instance count works; jobs coordinate via Postgres advisory locks.
6. IDs are app-generated UUIDv7 - time-ordered; the id is the pagination cursor.
7. Policy gating is server-side. UI reflects policy; it never is the enforcement.

## Milestones & status

| Milestone | Issue | Status |
|---|---|---|
| M0 scaffold - M1 schema - M2 public catalog | #1 #2 #4 | done (PR #14) |
| M3 auth & accounts | #5 | auth + accounts complete (register/login/logout, sessions, email verify, CSRF, password reset, change password, log-out-all, per-IP + per-account throttle, settings display name/locale/country, delete stub); remaining acceptance "unverified cannot list/chat" is enforced when M4/M5 add those endpoints |
| M4 seller flows (listings) | #7 | done - listings CRUD, past-race expiry job, /sell + /account/listings (sub-issues #29-#31, PRs #33-#35); image upload moved to M5 per D15 |
| M5 chat (the beta's core) | #8 | done - threads/polling/inbox + ack gate, private images (minio-go), report/block + retention job (sub-issues #36-#39, PRs #41-#44) |
| M8 i18n (en + es per D4) | #9 | done - hand-rolled per D14/D17; sub-issues #45 (foundation + en, PR #48) #46 (Spanish + localized emails, PR #50); server/API error-string i18n via codes shipped in #49 (PR #51) |
| M7 trust/safety (lite gates beta; full gates payments) | #11 | M7-lite done - legal pages: ToS (#53, PR #55) + Privacy/contact/cookies (#54, PR #58); report button shipped in M5. M7-full (GDPR/DSA, DAC7 memo, moderation queue) gates payments |
| M9 production (starts pre-beta per D8) | #12 | M9.1 self-host beta LIVE - Cloudflare Tunnel + Caddy + MinIO + Brevo (D20-D24). M9-full remaining: automated backups + rehearsed restore drill, observability (Sentry/uptime), CI/CD, staging, CSP |
| Gate 1: chat-only public beta (D3) | #13 | - |
| M6 payments (Stripe Connect) | #10 | - |
| Gate 2: payments GA | #13 | - |
| Ops: race catalog & policy verification (ES-first per D4) | #6 | ongoing |

## Verification protocol (D6)

The pyramid, bottom-up:
1. Unit tests - pure logic (state machines, formatters, cursors). No I/O.
2. Schema constraint tests - real Postgres, each in a rolled-back transaction; prove the DB rejects illegal data.
3. Handler integration tests - `httptest` against real Postgres with self-cleaning fixtures; prove endpoints, authz, pagination, error envelopes.
4. Smoke harness (`make smoke`) - boots the actual stack against seed data and asserts end-to-end behavior, including the policy-gating matrix (Buy exists only on `platform_sale`, disclaimers render per mode, drafts 404).
5. CI - all of the above plus lint, type-check, sqlc drift, migrations round-trip (up -> 0 -> up), on every PR.

Working rules:
- `make verify` must pass before any commit; `make smoke` before any PR that touches behavior.
- Every issue's acceptance criteria get a line in the PR description naming the command or test that proves each one.
- New behavior ⇒ a test that fails without the change. Bug fix ⇒ a test that reproduces it first.
- A milestone is done when: acceptance criteria proven, `verify` + `smoke` green, CI green, docs/CONTEXT.md updated if a decision was touched.

## Known notes & accepted trade-offs

- `ListRaces` correlated subquery for `active_listings`: bounded by page size (<=100 index-only lookups via the partial index), not unbounded N+1. Revisit if p95 on `/api/v1/races` degrades (then: `LEFT JOIN … COUNT(*) FILTER` or a counter column).
- `Cache-Control` on catalog pages is gated on auth (resolved, was the M3 checklist item): signed-in responses are `private, no-store` because the page embeds the layout nav, anonymous stay `public, max-age=60`. Applies to `/races`, `/races/[slug]`, `/listings/[id]`; a smoke assertion guards it.
- Exact-page-size pagination edge: when results == page size, a `next_cursor` is emitted whose page is empty. Harmless; the empty state renders. Fix only if users notice.
- `text-scale` meta tag in `app.html` comes from the official SvelteKit scaffold; inert in browsers, deliberately left untouched.
- sqlc/goose run as pinned `go run pkg@version` (not go.mod tools): sqlc's tree forces a `go 1.26` directive that golangci-lint can't lint yet. Revisit when golangci targets current Go.

## Open questions (decide by the milestone that needs them)

| Question | Decide by | Current lean |
|---|---|---|
| Refund fee absorption (Stripe keeps the original fee on refunds) | M6 | Platform eats it at beta volumes |
| Hosting provider | decided (D20) | Self-host for beta; Hetzner the managed migration target |
| Prod email provider | decided (D22) | Brevo smarthost, subdomain-authenticated (SPF+DKIM) |
| Prod object storage | Pre-beta (MinIO covers dev) | Cloudflare R2 |
| Branding / domain ("Bibseller" is a working name) | Pre-beta | - |
| Playwright browser E2E | Revisit at M4 | Deferred (D6) |
