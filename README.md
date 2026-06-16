# Bibseller

> Working name - final branding is tracked in [Decisions #3](https://github.com/leonfullxr/bibseller/issues/3).

A non-profit, peer-to-peer marketplace for race bibs, EU-wide: runners who can't make it to the start line meet runners who missed registration. Zero commission - the goal is connecting people, not profit.

## The core idea

Everything in this product hinges on one fact: every race has its own transfer rules. Some allow resale, some run their own official name-change process, some forbid transfers outright. Each race in the catalog therefore carries a `transfer_policy` that drives every downstream behavior - which buttons exist, whether money can flow, which disclaimers render:

| Policy | Meaning | What the platform provides |
|---|---|---|
| `platform_sale` | Race allows bib resale | Listing + chat + optional secure payment (funds held until the transfer is confirmed) |
| `official_only` | Race runs its own name-change process | Listing + chat + prominent link to the official procedure - we never touch money |
| `connect_only` | Resale restricted or gray-area | Listing + chat only, strong disclaimer, recorded buyer acknowledgment |
| `unknown` | Policy not yet verified | Treated exactly like `connect_only` |

Payments, where available, run through Stripe Connect at zero commission. The platform structurally prevents payment flows for races that don't allow transfers - it's not a UI toggle, it's enforced at the API and database layer.

## Stack

SvelteKit (Svelte 5, TypeScript) - Go - PostgreSQL 16 - sqlc + goose - Stripe Connect - S3-compatible object storage - Paraglide i18n. One monorepo, no microservices, Postgres for everything until scale demands otherwise. Full rationale: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Documentation

| Doc | Contents |
|---|---|
| [docs/CONTEXT.md](docs/CONTEXT.md) | Decision log, working agreements, verification protocol - read first |
| [docs/PRODUCT.md](docs/PRODUCT.md) | Mission, personas, user journeys, the policy matrix as UX |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Stack & rationale, repo layout, dev environment, API conventions, auth design, scaling path |
| [docs/DATA_MODEL.md](docs/DATA_MODEL.md) | Schema, constraints, state machines, indexes, retention |
| [docs/PAYMENTS_AND_COMPLIANCE.md](docs/PAYMENTS_AND_COMPLIANCE.md) | Money flows, fee policy, GDPR / DSA / DAC7 |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Milestones mapped to issues, launch gates |

## Status

Done so far: M0 scaffold ([#1](https://github.com/leonfullxr/bibseller/issues/1)), M1 schema ([#2](https://github.com/leonfullxr/bibseller/issues/2)), M2 public catalog ([#4](https://github.com/leonfullxr/bibseller/issues/4)). Next: M3 auth ([#5](https://github.com/leonfullxr/bibseller/issues/5)). Live status: [#13](https://github.com/leonfullxr/bibseller/issues/13).

## Development

Prereqs: Go 1.25+, Node 22+, Docker. Optional: [`air`](https://github.com/air-verse/air) for Go hot-reload (`make dev` falls back to `go run`), `golangci-lint`, Stripe CLI (M6+). No install needed for `goose`/`sqlc` - the Makefile runs pinned versions via `go run`.

```sh
cp .env.example .env    # optional - defaults work without it
make dev                # Postgres/MinIO/Mailpit + Go API :8080 + SvelteKit :5173
```

Verify: <http://localhost:5173> shows “API connected”, and `curl localhost:5173/api/healthz` returns `{"status":"ok"}` through the Vite proxy.

| Target | Does |
|---|---|
| `make dev` | infra (compose) + both apps with hot reload |
| `make migrate` / `make migrate-down` | goose migrations |
| `make sqlc` | regenerate type-safe query code |
| `make seed` | wipe + load dev data (20 races, all policy modes) |
| `make test` / `make lint` | both halves |
| `make verify` | pre-commit gate: lint + typecheck + tests + sqlc drift |
| `make smoke` | end-to-end assertions against the seeded stack (wipes dev data) |

Mailpit UI: `localhost:8025` - MinIO console: `localhost:9001` (minioadmin / minioadmin). More detail: [docs/ARCHITECTURE.md -> Dev environment](docs/ARCHITECTURE.md#dev-environment).

## License

See [LICENSE](LICENSE).
