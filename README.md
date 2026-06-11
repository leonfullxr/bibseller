# Bibseller

> Working name — final branding is tracked in [Decisions #3](https://github.com/leonfullxr/bibseller/issues/3).

A **non-profit, peer-to-peer marketplace for race bibs, EU-wide**: runners who can't make it to the start line meet runners who missed registration. Zero commission — the goal is connecting people, not profit.

## The core idea

Everything in this product hinges on one fact: **every race has its own transfer rules.** Some allow resale, some run their own official name-change process, some forbid transfers outright. Each race in the catalog therefore carries a `transfer_policy` that drives every downstream behavior — which buttons exist, whether money can flow, which disclaimers render:

| Policy | Meaning | What the platform provides |
|---|---|---|
| `platform_sale` | Race allows bib resale | Listing + chat + **optional secure payment** (funds held until the transfer is confirmed) |
| `official_only` | Race runs its own name-change process | Listing + chat + prominent link to the official procedure — we never touch money |
| `connect_only` | Resale restricted or gray-area | Listing + chat only, strong disclaimer, recorded buyer acknowledgment |
| `unknown` | Policy not yet verified | Treated exactly like `connect_only` |

Payments, where available, run through Stripe Connect at **zero commission**. The platform structurally prevents payment flows for races that don't allow transfers — it's not a UI toggle, it's enforced at the API and database layer.

## Stack

SvelteKit (Svelte 5, TypeScript) · Go · PostgreSQL 16 · sqlc + goose · Stripe Connect · S3-compatible object storage · Paraglide i18n. One monorepo, no microservices, Postgres for everything until scale demands otherwise. Full rationale: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Documentation

| Doc | Contents |
|---|---|
| [docs/PRODUCT.md](docs/PRODUCT.md) | Mission, personas, user journeys, the policy matrix as UX |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Stack & rationale, repo layout, dev environment, API conventions, auth design, scaling path |
| [docs/DATA_MODEL.md](docs/DATA_MODEL.md) | Schema, constraints, state machines, indexes, retention |
| [docs/PAYMENTS_AND_COMPLIANCE.md](docs/PAYMENTS_AND_COMPLIANCE.md) | Money flows, fee policy, GDPR / DSA / DAC7 |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Milestones mapped to issues, launch gates |

## Status

**Planning.** Implementation starts with the [M0 scaffold (#1)](https://github.com/leonfullxr/bibseller/issues/1); overall progress is tracked in [#13](https://github.com/leonfullxr/bibseller/issues/13).

## Development (lands with M0)

```sh
docker compose up -d   # Postgres + MinIO + Mailpit
make dev               # Go API (air hot-reload) + SvelteKit (vite)
```

Toolchain: Go 1.24+, Node 20+, `sqlc`, `goose`, `air`, Stripe CLI. See [docs/ARCHITECTURE.md → Dev environment](docs/ARCHITECTURE.md#dev-environment).

## License

See [LICENSE](LICENSE).
