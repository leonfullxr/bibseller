# Design language: Negative Split

The frontend's visual system, an exploration branch alternative to "The
Runner's Journal" (D28 - see git history of this file). One idea carries
everything: **run the second half faster** - a contemporary, athletic-modern
surface: cool porcelain, night-navy ink, a single violet accent, a chunky
grotesk display face, pill buttons and generous radii. This file is the
reference for anyone styling new UI; the single source of truth for values
is `frontend/src/routes/layout.css`.

## Voice

Fast, friendly, trustworthy. Surfaces feel like a well-built product: white
cards on cool porcelain, soft ambient shadows, generous whitespace. Emphasis
comes from the display grotesk and one saturated accent color, never from
hairlines or blocks.

## Tokens (layout.css `:root`)

| Token | Value | Role |
| --- | --- | --- |
| `--paper` | `#f7f8fb` | Page background - cool porcelain |
| `--paper-2` | `#eef0f6` | Tinted wells (chat canvas, mats) |
| `--ink` | `#131629` | Night navy: text; the footer block |
| `--ink-2` | `#363b52` | Secondary ink |
| `--brand-50..800` | Violet ramp | Brand + every primary action. White text needs `brand-600`+ (AA, 5:1+; `brand-700`+ for 7:1) |
| `--slate-*` | True cool slate | Neutrals (the historical names carry real slate again on this design) |
| `--emerald-*` / `--sky-*` / `--amber-*` | unchanged | **Semantic only**: success/allowed, official process, caution/connect-only. Identical across designs |
| `--font-display` | Bricolage Grotesque Variable | Headings, wordmark, prices, step numerals (weights 700-800) |
| `--font-body` | Manrope Variable | Everything else; `tabular-nums` app-wide |
| `--radius-sm / --radius / --radius-lg` | 0.5 / 0.75 / 1rem | Shape scale; buttons and chips are full pills (`9999px`) |
| `--shadow-hard(-sm)` | Soft ambient shadows | Card depth. Names kept from earlier systems so components restyle through the variables |

## Rules

- **Brand vs semantic**: violet means "act here"; green/sky/amber/slate keep
  their policy meanings. Never a violet status, never a green button.
- Headings are the display grotesk, sentence case, `-0.02em`; weight 700-800
  (chunky, never condensed or serif).
- Buttons and chips are **pills** (`border-radius: 9999px`); `.btn-primary` =
  `brand-600`, hover `brand-700`. Body face, weight 700.
- Cards are soft plates: white, 1px `slate-200`, `--radius-lg`, ambient
  shadow; hover darkens the border, lifts the shadow, and floats the card
  `-2px`. No hairline rules on `--ink`, no hard offsets, no serifs - those
  belong to earlier systems.
- The **dawn gradient** is the signature and appears exactly once: the home
  hero (violet first light over porcelain, radial + linear). Nothing else on
  the site gets a gradient.
- The header is a **floating capsule**: sticky, frosted white
  (`rgb(255 255 255 / 0.85)` + `backdrop-filter: blur(12px)`), fully rounded,
  detached from the page top. The footer is the counterweight: a night-navy
  (`--ink`) block with the wordmark in white.
- Home devices: the four how-it-works steps are numbered cards (display-face
  numerals in violet); the buyer/seller journey keeps the paired-columns
  "duet" - seller lane violet, buyer lane `sky-700`; the three policy-mode
  cards wear their semantic tone as a 3px top border (emerald / sky / amber).
- Header actions (log out, language) are small bordered pill chips at nav
  type size - actions get frames, links stay bare.
- Motion budget: 150ms color/translate transitions; the nav progress bar
  respects `prefers-reduced-motion`.
- Keep `overflow-x: clip` on html; never suppress focus outlines
  (`brand-600` ring); the four-policy-mode question applies to anything
  policy-adjacent.

## Fonts

Self-hosted fontsource woff2 (strict CSP, same-origin): Bricolage Grotesque
Variable (display; use 700-800 - it gets soft below 600) + Manrope Variable
(body, 400-800 via the wght axis). Imported in `+layout.svelte`; Fraunces
and Barlow were retired with the journal.
