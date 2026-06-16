# Product brief

## Mission

Race bibs go to waste every season: injured or busy runners eat the entry fee while others refresh sold-out registration pages. Bibseller connects the two - EU-wide, peer-to-peer, zero commission. The goal is not profit; it is making the match happen safely and within each race's rules.

Founding stance:

- We take no cut. Where payments exist they run at cost; the only money the platform passes through is the payment processor's own fee, itemized to the buyer (decided - [D1](CONTEXT.md#founder-decisions-log)).
- We respect race rules. Each race's transfer policy is researched, sourced, and encoded; the platform structurally refuses to move money for races that don't allow resale.
- We are a venue, not a party. For non-`platform_sale` races, the platform connects people and nothing more; this is stated in the UI, acknowledged by users, and recorded.

## Non-goals (v1)

- No revenue, no ads, no premium tiers.
- No race-organizer SaaS (bib management, official transfer tooling) - possible later.
- No scalping venue: prices are hard-capped at face value when the original price is known (decided - [D2](CONTEXT.md#founder-decisions-log)).
- No native mobile apps - the site is responsive; a PWA pass can come later.
- No anonymous contact: chatting requires a verified account (anti-spam, traceability).

## Personas

| Persona | Situation | Needs |
|---|---|---|
| Marta, seller | Registered for the Valencia Marathon, tore a ligament in October | Recover (some of) her entry fee; hand the bib to a real runner; not break the race's rules |
| Jonas, buyer | Berlin Half sold out in hours; he's trained all winter | Find a legitimate bib; not get scammed; know what the race officially allows |
| Race organizer (indirect) | Runs an official name-change window for €10 | Their rules respected; ideally fewer black-market bibs running under wrong names |

The organizer is not a v1 user, but the `official_only` mode is designed to keep races friendly toward the platform rather than hostile.

## The policy matrix (canonical UX definition)

`races.transfer_policy` is the single source of truth. Every surface derives from it:

| Capability | `platform_sale` | `official_only` | `connect_only` | `unknown` |
|---|---|---|---|---|
| Listing visible in catalog | yes | yes | yes | yes |
| Asking price shown | yes | yes (informational) | yes (informational) | yes (informational) |
| Chat between buyer & seller | yes | yes | yes after acknowledgment | yes after acknowledgment |
| Buy button / payment through platform | yes | no | no | no |
| Link-out to official transfer procedure | n/a | yes prominent | n/a | n/a |
| Disclaimer level | Standard marketplace terms | "Agree terms here; the transfer itself happens via the race's official process" | Strong: platform handles no money, takes no responsibility; the race's own rules may restrict transfers | Strong + "policy not yet verified" |
| Buyer acknowledgment recorded (`policy_acks`) | - | - | yes before first message | yes before first message |
| Platform legal posture | Facilitates payment, holds funds until confirmation | Venue only | Venue only | Venue only |

Two structural guarantees follow:

1. No payment path exists for non-`platform_sale` races - enforced in the checkout API and by DB constraints, not just hidden buttons.
2. Acknowledgments are evidence. Where the platform is a venue only, the buyer ticked a box saying exactly that, and we stored when.

## User journeys

### Seller (all modes)

1. Sign up, verify email.
2. "Sell a bib" -> search the race catalog (if missing -> "suggest a race", lands as `unknown`).
3. The form explains what this race's policy means for them before they invest effort ("this race = chat only - you arrange everything with the buyer").
4. Set asking price (capped per Decision 2), optional description + photo proof of registration.
5. Listing goes live; buyers open chat threads; seller manages everything from `/account/listings` + `/account/inbox`.
6. Close the loop: mark as sold (or it auto-expires after race day).

### Buyer - `platform_sale` race

1. Find race -> listing -> Buy (or chat first).
2. Checkout via Stripe (price + pass-through processing fee, itemized).
3. Funds are held by the platform - not the seller - while the bib/name transfer happens.
4. Seller marks "transferred"; buyer confirms (or auto-release after N days); seller is paid out.
5. Problems -> dispute -> manual resolution, refund path documented.

### Buyer - `official_only` race

1. Find race -> listing -> chat with the seller, agree on the handover.
2. Prominent panel: "This race runs an official name-change process" -> link out, deadlines and fees shown when known.
3. The actual transfer (and usually payment of the official fee) happens on the race's site. The platform's job ends at the introduction.

### Buyer - `connect_only` / `unknown` race

1. Find race -> listing -> "Contact seller" triggers a one-time interstitial: the race's rules apply, the platform moves no money and carries no responsibility -> must acknowledge (recorded).
2. Chat opens; everything beyond is between the two users.

## Trust & safety features (v1)

- Verified email required to list or chat.
- Policy acknowledgments recorded with timestamps.
- Report listing / report message -> moderation queue; block user.
- Hard price cap at face value when known, warning banner otherwise ([D2](CONTEXT.md#founder-decisions-log)) - removes most scalper incentive.
- For paid orders: funds held until buyer confirmation; itemized fees; full order event log.

## Launch strategy

Beta = chat-only marketplace (every race effectively connect/official mode, payments off). This validates the matching idea with the smallest possible legal and compliance surface. Payments GA follows once Stripe flows, compliance, and the verified `platform_sale` catalog are solid. Gates and ordering: [docs/ROADMAP.md](ROADMAP.md).

## Success metrics (beta)

- Listings created / week; threads started per listing.
- Self-reported successful handoffs ("mark as sold" with a buyer attached).
- % of catalog races with verified policies ([Ops #6](https://github.com/leonfullxr/bibseller/issues/6)).
- Report rate (low = healthy; zero = nobody's using it).
