# CLAUDE.md

Project instructions for AI-assisted work on Bibseller. Behavioral guidelines below are the contract; the project section makes them concrete.

Tradeoff: these guidelines bias toward caution over speed. For trivial tasks, use judgment.

## Project: Bibseller

Zero-commission, EU-wide P2P marketplace for race bibs. Everything derives from `races.transfer_policy` (`platform_sale | official_only | connect_only | unknown`) - which buttons exist, whether money can flow, which disclaimers render. Policy gating lives server-side and in DB constraints, never only in UI.

Read before significant work: [docs/CONTEXT.md](docs/CONTEXT.md) (decisions log - binding), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), [docs/DATA_MODEL.md](docs/DATA_MODEL.md). Status lives in issue #13.

### Stack & layout

Go 1.25 (stdlib `net/http` ServeMux, pgx, sqlc - no ORM) - SvelteKit (Svelte 5 runes, TS strict, scoped component CSS, adapter-node) - Postgres 16 - goose migrations. Monorepo: `backend/` (package-by-domain under `internal/`, shared infra in `internal/platform/`), `frontend/`. Heavy tools run as pinned `go run pkg@version` from the Makefile, deliberately not in `go.mod`.

### Commands

```sh
make dev            # infra (compose) + API :8080 + web :5173, hot reload
make migrate        # goose up        (migrate-down: one step back)
make seed           # dev-only wipe + load (20 races, all 4 policy modes)
make verify         # the pre-commit gate: lint + typecheck + tests + sqlc drift
make smoke          # boots the stack, asserts end-to-end policy behavior
make sqlc           # regenerate after editing db/queries/*.sql - commit the output
```

### Branching & deploy

`main` is the trunk and runs staging; a long-lived `production` branch runs the live site. Work on a feature branch, PR into `main`, verify on staging (`make staging-up`), then `make promote` (fast-forwards `production` to the tested `main` and pushes) and deploy with `make prod-migrate && make prod-up`. Staging and prod are the same compose stack on one self-host box, isolated by compose project name; each runs from its own git worktree. Backups: `make prod-backup-offsite` (nightly cron) and `make prod-restore-drill`. Full detail in [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) and CONTEXT.md (D20, D25, D26).

### Conventions

- SQL lives in `backend/db/queries/*.sql`; `make sqlc`; commit generated code (CI checks drift).
- IDs: app-generated UUIDv7 via `internal/platform/ids` (id doubles as pagination cursor). Money: integer cents + currency. Time: `timestamptz`, UTC.
- API: `/api/v1`, snake_case JSON, error envelope via `httpx.Error` (stable `code` strings).
- Tests that need Postgres use `internal/platform/db/testdb` (skips cleanly when absent) and self-cleaning fixtures - never depend on seed rows.
- Frontend internal links use `resolve()` (`svelte/no-navigation-without-resolve` is on). Server `load` functions call the Go API via `$lib/api/server.ts`.
- Lint/format are CI-enforced: `golangci-lint`, Prettier, ESLint, `svelte-check`.

## Agent skills

### Issue tracker

Issues and PRDs live in this repo's GitHub Issues, driven via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage roles map to their default label names (reusing the existing `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `CONTEXT.md` lives at `docs/CONTEXT.md`, ADRs at `docs/adr/`. See `docs/agents/domain.md`.

## 1. Think Before Coding

Don't assume. Don't hide confusion. Surface tradeoffs.

- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.
- Bibseller-specific: any feature touching listings, chat, or money must answer "what happens in each of the four policy modes?" before implementation starts.

## 2. Simplicity First

Minimum code that solves the problem. Nothing speculative.

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.
- Bibseller-specific: no new infrastructure (Redis, queues, services) - the decision log in ARCHITECTURE.md names the trigger that would change each "no".

Ask: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

Touch only what you must. Clean up only your own mess.

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.
- Remove imports/variables/functions that YOUR changes made unused; leave pre-existing dead code alone unless asked.

The test: every changed line traces directly to the request.

## 4. Goal-Driven Execution

Define success criteria. Loop until verified.

- Transform tasks into verifiable goals: "add validation" -> "write tests for invalid inputs, make them pass"; "fix the bug" -> "write the reproducing test first".
- For multi-step tasks, state a brief plan: `[step] -> verify: [check]` per step.
- Bibseller protocol (decision D6): `make verify` before every commit; `make smoke` before any behavior-touching PR; every issue's acceptance criteria get a line in the PR description naming the command or test that proves it; new behavior ⇒ a test that fails without the change.

---

These guidelines are working if: diffs contain only necessary changes, clarifying questions come before implementation rather than after mistakes, and every PR demonstrates its acceptance criteria rather than asserting them.
