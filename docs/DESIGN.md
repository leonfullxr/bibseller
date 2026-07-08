# Design language: The Runner's Journal

The frontend's visual system, shipped 2026-07-08 after a one-day trial of its
predecessor ("race day", D27 - see git history of this file). One idea
carries everything: **a calm, premium sports journal** - ivory paper, warm
charcoal ink, a single bordeaux accent, serif display type, hairline rules.
This file is the reference for anyone styling new UI; the single source of
truth for values is `frontend/src/routes/layout.css`.

## Voice

Quiet, editorial, trustworthy. Surfaces feel like a printed journal: white
plates on ivory, hairline rules, generous whitespace. Emphasis comes from the
serif display face and one deep accent color, never from blocks or shadows.

## Tokens (layout.css `:root`)

| Token | Value | Role |
| --- | --- | --- |
| `--paper` | `#f7f3ea` | Page background - ivory |
| `--paper-2` | `#efe9db` | Mat tints (chat canvas, wells, closing panel) |
| `--ink` | `#26221c` | Warm charcoal: text and hairline rules |
| `--ink-2` | `#3d372e` | Secondary ink |
| `--brand-50..800` | Bordeaux ramp | Brand + every primary action. White text needs `brand-600`+ (7:1+) |
| `--slate-*` | Warm stone values | Neutrals (historical names, deliberate - see D27 note) |
| `--emerald-*` / `--sky-*` / `--amber-*` | unchanged | **Semantic only**: success/allowed, official process, caution/connect-only. Identical across designs |
| `--font-display` | Fraunces Variable | Headings (sentence case), wordmark, prices, folio numbers |
| `--font-body` | Barlow | Everything else; `tabular-nums` app-wide |
| `--shadow-hard(-sm)` | Soft layered shadows | "Plate" shadows. Names kept from the race-day system so components restyle through the variables |

## Rules

- **Brand vs semantic**: bordeaux means "act here"; green/sky/amber/slate keep
  their policy meanings. Never a bordeaux status, never a green button.
- Headings are serif, sentence case, `-0.01em`; the italic serif in bordeaux
  is the accent voice (hero highlight, wordmark suffix).
- Eyebrow style for small labels: `0.6875-0.75rem`, `0.08-0.1em` tracking,
  uppercase (nav, card dates, lane tags).
- Cards are **plates**: white, 1px `slate-200` hairline, 0.25rem radius, soft
  shadow; hover darkens the hairline and lifts the shadow. No hard offsets,
  no thick borders, no diagonals - those belong to the race-day system.
- Buttons: body face, sentence case, weight 600, near-square radius;
  `.btn-primary` = `brand-700`, hover `brand-800`.
- Section headings sit over a short 1px hairline (the shared device).
- Motion budget: 150ms color transitions only.
- Keep `overflow-x: clip` on html; never suppress focus outlines
  (`brand-600` ring); the four-policy-mode question applies to anything
  policy-adjacent.

## Fonts

Self-hosted fontsource woff2 (strict CSP, same-origin): Fraunces Variable
(wght axis; use 500-650 - it gets heavy above 700) + Barlow 400-700.
Imported in `+layout.svelte`; Barlow Condensed was retired with race-day.
