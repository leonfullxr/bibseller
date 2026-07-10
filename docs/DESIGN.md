# Design language: Sunday Run

The frontend's visual system, an exploration branch alternative to "The
Runner's Journal" (D28 - see git history of this file). One idea carries
everything: **the warmth of a club run** - chalk paper, espresso ink, a
single poppy-red accent (race-poster red, the track-clay family), a slab
display face, pill buttons and generous radii. Modern structure, human
surface. This file is the reference for anyone styling new UI; the single
source of truth for values is `frontend/src/routes/layout.css`.

## Voice

Warm, direct, trustworthy. Surfaces feel like a well-loved clubhouse
noticeboard rebuilt as a product: white cards on warm chalk, soft ambient
shadows, generous whitespace. Emphasis comes from the slab display face and
one warm accent color - never from gradients, frosted glass, or hairline
rules (those belong to earlier systems).

## Tokens (layout.css `:root`)

| Token | Value | Role |
| --- | --- | --- |
| `--paper` | `#faf7f2` | Page background - warm chalk |
| `--paper-2` | `#f1ece3` | Tinted wells (chat canvas, mats, date chips) |
| `--ink` | `#2b2521` | Warm espresso: text; the footer block |
| `--ink-2` | `#4a423b` | Secondary ink |
| `--brand-50..800` | Poppy ramp | Brand + every primary action. White text needs `brand-600`+ (AA, 5:1+; `brand-700`+ for ~7:1) |
| `--slate-*` | Warm stone | Neutrals (historical names, warm values - same convention as the journal) |
| `--emerald-*` / `--sky-*` / `--amber-*` | unchanged | **Semantic only**: success/allowed, official process, caution/connect-only. Identical across designs |
| `--font-display` | Bitter Variable | Headings, wordmark, prices, step numerals (weights 700-800, slab) |
| `--font-body` | Source Sans 3 Variable | Everything else; `tabular-nums` app-wide |
| `--radius-sm / --radius / --radius-lg` | 0.5 / 0.75 / 1rem | Shape scale; buttons and chips are full pills (`9999px`) |
| `--shadow-hard(-sm)` | Soft warm shadows | Card depth. Names kept from earlier systems so components restyle through the variables |

## Rules

- **Brand vs semantic**: poppy means "act here"; green/sky/amber/slate keep
  their policy meanings. Never a poppy status, never a green button.
- Headings are the slab face, sentence case, `-0.01em`, weight 700-800.
  No condensed type, no grotesk-at-800, no serifs-as-decoration.
- Buttons and chips are **pills** (`border-radius: 9999px`); `.btn-primary` =
  `brand-600`, hover `brand-700`. Body face, weight 700.
- Cards are soft plates: white, 1px `slate-200`, `--radius-lg`, ambient
  shadow; hover darkens the border, lifts the shadow, floats the card `-2px`.
- **No gradients anywhere.** The one hand-made device is the hero's marker
  underline (`text-decoration` in `brand-300` under the accent words).
  Date chips on cards are neutral (`paper-2`/`ink-2`) so the policy badge is
  the only color on the card strip.
- The header is a **solid white floating capsule**: sticky, fully rounded,
  1px stone border, soft shadow - no backdrop blur. The footer is the
  counterweight: a warm espresso (`--ink`) block with the wordmark in white.
- Home devices: the four how-it-works steps are numbered cards (slab
  numerals in poppy); the buyer/seller journey keeps the paired-columns
  "duet" - seller lane poppy, buyer lane `sky-700`; the three policy-mode
  cards wear their semantic tone as a 3px top border (emerald / sky / amber).
- Header actions (log out, language) are small bordered pill chips at nav
  type size - actions get frames, links stay bare.
- Motion budget: 150ms color/translate transitions; the nav progress bar
  respects `prefers-reduced-motion`.
- Keep `overflow-x: clip` on html; never suppress focus outlines
  (`brand-600` ring); the four-policy-mode question applies to anything
  policy-adjacent.

## Fonts

Self-hosted fontsource woff2 (strict CSP, same-origin): Bitter Variable
(display slab; use 700-800) + Source Sans 3 Variable (body, wght axis).
Imported in `+layout.svelte`; Bricolage Grotesque and Manrope were retired
with the "Negative Split" iteration of this branch (git history).
