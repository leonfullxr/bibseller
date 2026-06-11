# Roadmap

Live status: [tracking issue #13](https://github.com/leonfullxr/bibseller/issues/13). Open decisions (none block the start): [#3](https://github.com/leonfullxr/bibseller/issues/3).

## Milestones

| Milestone | Issue | Effort | Depends on |
|---|---|---|---|
| M0 — Scaffold: monorepo, dev env, CI | [#1](https://github.com/leonfullxr/bibseller/issues/1) | M | — |
| M1 — Schema v1: migrations, sqlc, seed | [#2](https://github.com/leonfullxr/bibseller/issues/2) | M | M0 |
| M2 — Public catalog (read-only, policy-aware) | [#4](https://github.com/leonfullxr/bibseller/issues/4) | M | M1 |
| M3 — Auth & accounts | [#5](https://github.com/leonfullxr/bibseller/issues/5) | M | M1 |
| M4 — Seller flows: listings + uploads | [#7](https://github.com/leonfullxr/bibseller/issues/7) | L | M2, M3 |
| M5 — Chat (polling v1) | [#8](https://github.com/leonfullxr/bibseller/issues/8) | L | M3 |
| M6 — Payments: Stripe Connect + state machine | [#10](https://github.com/leonfullxr/bibseller/issues/10) | XL | M4 (M5 first recommended) |
| M7 — Trust, safety & compliance | [#11](https://github.com/leonfullxr/bibseller/issues/11) | L | lite: M3 · full: M5/M6 |
| M8 — i18n (Paraglide, en + es) | [#9](https://github.com/leonfullxr/bibseller/issues/9) | M | M2 |
| M9 — Production: hosting, observability, backups | [#12](https://github.com/leonfullxr/bibseller/issues/12) | L | skeleton: M2 · full: pre-beta |
| Ops — Race catalog & policy verification | [#6](https://github.com/leonfullxr/bibseller/issues/6) | M, ongoing | M1; fully parallel |

## Phases & gates

```
Phase 1  M0 → M1                      foundations
Phase 2  M2 → M9-lite (deploy!)       first end-to-end slice, walking skeleton in prod
         Ops catalog research starts and never stops
Phase 3  M3 → M4 → M5 → M8 → M7-lite  accounts, listings, chat, i18n, legal pages
─────────🚀 GATE 1: PUBLIC BETA ───────────────────────────────────────────────
         Chat-only marketplace. No payments. Smallest legal surface.
Phase 4  M6 → M7-full → M9-full       payments, full compliance, hardened prod
─────────🚀 GATE 2: PAYMENTS GA ───────────────────────────────────────────────
         Secure-payment toggle live for verified platform_sale races.
```

## Why this order

- **Read-only first (M2):** teaches SvelteKit SSR + Go handlers + sqlc with zero risk and no irreversible decisions.
- **Deploy a walking skeleton right after M2:** ops surprises surface while the stack is tiny, not the week before launch.
- **Payments dead last (M6):** the hardest, most regulated part — and the beta genuinely works without it. Payments-first is the classic way to get demoralized in week two.
- **Catalog research (Ops) in parallel:** non-coding work; the catalog's quality is the product's quality.
- **Beta before payments:** validates the matching idea with minimal compliance surface; payment volume arrives only after there's evidence people want the matches.

## Post-v1 backlog (parked, by design)

SSE/WebSocket chat upgrade (trigger metrics in [ARCHITECTURE.md](ARCHITECTURE.md#decision-log)) · OAuth: Google, Strava · more locales (de, fr, it, nl) · multi-currency (SEK/PLN/DKK) · Meilisearch if Postgres FTS maxes out · "suggest a race" self-service intake · PWA polish · race-organizer portal (official_only races could manage their link-outs themselves).
