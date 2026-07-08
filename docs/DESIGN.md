# Design language: race day

The frontend's visual system, shipped in #198/#199 (2026-07-07). One idea
carries everything: **the interface of a race, not a web app** - the bib, the
start-gun, the timing board, the finish tape. This file is the reference for
anyone (human or agent) styling new UI; the single source of truth for values
is `frontend/src/routes/layout.css`.

## Voice

Bold, warm, physical. Surfaces feel like paper and print, not glass and blur.
Energy comes from ink blocks, start-gun orange, and condensed uppercase type;
trust comes from generous spacing, quiet panels, and strictly semantic color.

## Tokens (layout.css `:root`)

| Token | Value | Role |
| --- | --- | --- |
| `--paper` | `#faf6ef` | Page background - warm paper |
| `--paper-2` | `#f3ede1` | Tinted wells (chat canvas, card strips) |
| `--ink` | `#171c26` | Text, color blocks (header/footer/hero/board), hard shadows |
| `--ink-2` | `#2a3140` | Secondary ink: hovers and wells on ink blocks |
| `--brand-50..800` | Tailwind orange ramp | Brand + every primary action ("start-gun orange") |
| `--slate-50..900` | Warm stone values | Neutrals. **Historical names**: the `slate-*` tokens deliberately hold stone values so ~200 scoped usages shifted warm without edits |
| `--emerald-*` | Tailwind emerald | **Semantic only**: success, "resale allowed", active listing counts. Not the brand (it was, pre-#198) |
| `--sky-*` / `--amber-*` | unchanged | Semantic: official-process / caution+connect-only |
| `--font-display` | Barlow Condensed | Headings (globally uppercase), buttons, wordmark, bib numbers |
| `--font-body` | Barlow | Everything else; `tabular-nums` app-wide |
| `--shadow-hard(-sm)` | `4px/3px offset, no blur, ink` | Poster shadows on cards and primary buttons |

## Color rules (contrast-checked, keep them)

- **Brand vs semantic is the core discipline.** Orange means "act here";
  green means "allowed/succeeded"; sky means "official process"; amber means
  "caution / connect-only"; slate means "inactive/unknown". Never use orange
  for a status or green for a button.
- White text sits only on `brand-700`/`brand-800` (5.2:1+) or ink.
  `brand-500/600` are accent-only: text/shapes ON ink (6.1:1), borders,
  underline bars, decorative fills - never behind white body text.
- Small gray text on white is `slate-500` minimum.
- Own chat bubbles: `brand-700` bg, white text, `brand-50` time.

## Signature elements (the identity - reuse, don't reinvent)

1. **Diagonal cuts** (`clip-path`): hero bottom edge, footer top edge. The
   finish-tape angle. Used exactly twice per page; don't sprinkle it.
2. **Bib-tag cards** (`RaceCard.svelte` is the reference): white card, 2px ink
   border, hard shadow, a `--paper-2` top strip with two punched holes
   (`::before/::after` rings) and the date as the "bib number"; hover
   translates 1-2px while the shadow shrinks (the card "presses down").
3. **Timing board** (home journey): ink block, lane-coded rows - 4px accent
   edge + faint wash + lane-colored position number (orange seller / sky
   buyer), icons in `ink-2` circular wells.
4. **Course checkpoints** (home how-it-works): a 3px ink track line through
   numbered orange markers; the last marker is a checkered flag
   (`repeating-conic-gradient`), no number.
5. **Poster buttons**: `.btn-primary` is `brand-700` with a hard shadow that
   collapses on hover/active (translate toward the shadow). Button labels are
   condensed uppercase.

## Primitives (closed set - #147 still applies)

`.btn/.btn-primary/.btn-outline`, `.field`, `.panel`, `.alert`, `.pill`,
`.empty` in `layout.css`. Panels stay **quiet** (1px stone border, no shadow);
the loud treatment (ink borders + hard shadows) is reserved for cards and
CTAs, so pages keep a foreground/background rhythm. Do not add new global
classes; pages add scoped margins/overrides only.

## Type rules

- `h1-h3` are globally display-face, uppercase, `letter-spacing: 0.015em`,
  `text-wrap: balance`. Sizes stay per-component.
- Buttons: display face, uppercase, `0.04em` tracking.
- Numerals are tabular app-wide; prices/dates render as data, ink-colored and
  display-weight where they are the point (listing price = bib-number style).
- Fonts are self-hosted fontsource woff2 (CSP: same-origin only), weights
  400-800 body / 600-800 condensed, imported in `+layout.svelte`.

## Motion

150ms color transitions, 100ms translate/shadow on press - that's the whole
budget. No entrance animations, no parallax. The only long-running animation
is the pre-existing home marquee, which pauses on hover/focus and collapses
under `prefers-reduced-motion`.

## Do / don't

- DO run the four-mode question on anything policy-adjacent: the tones above
  are load-bearing semantics, not decoration.
- DO keep `overflow-x: clip` on html: full-bleed blocks use `50vw` negative
  margins and overflow by the scrollbar width otherwise.
- DON'T suppress focus outlines; the global ring is `brand-600` and must
  survive local styling (Copilot caught one, #198).
- DON'T put white body text on `brand-600` or lighter - it fails AA.
- DON'T reintroduce emerald as an action color; that era ended at #198.
