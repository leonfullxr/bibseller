# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

This is a single-context repo. The context doc lives under `docs/`, not at the repo root.

## Before exploring, read these

- **`docs/CONTEXT.md`** - the binding decisions log and domain glossary for Bibseller.
- **`docs/adr/`** - read ADRs that touch the area you're about to work in. This directory is created lazily when the first ADR is written; if it's absent, proceed silently.

Supporting docs worth reading for deeper work: `docs/ARCHITECTURE.md`, `docs/DATA_MODEL.md`, `docs/PAYMENTS_AND_COMPLIANCE.md`.

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The producer skill (`/grill-with-docs`) creates them lazily when terms or decisions actually get resolved.

## File structure

```
/
├── CLAUDE.md
├── docs/
│   ├── CONTEXT.md            ← decisions log + domain glossary
│   ├── ARCHITECTURE.md
│   ├── DATA_MODEL.md
│   └── adr/                  ← architectural decision records (created lazily)
│       ├── 0001-....md
│       └── 0002-....md
├── backend/
└── frontend/
```

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `docs/CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids. In particular, the four `transfer_policy` modes (`platform_sale`, `official_only`, `connect_only`, `unknown`) are the canonical names - use them verbatim.

If the concept you need isn't in the glossary yet, that's a signal - either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/grill-with-docs`).

## Flag ADR conflicts

If your output contradicts an existing ADR or a decision recorded in `docs/CONTEXT.md`, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 (or CONTEXT decision Dx) - but worth reopening because..._
