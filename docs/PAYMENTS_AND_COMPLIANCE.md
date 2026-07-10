# Payments & EU compliance

## Principles

1. Zero commission. The platform earns nothing on any transaction. Payment-processor fees are passed through to the buyer at cost, itemized at checkout (decided - [D1](CONTEXT.md#founder-decisions-log)).
2. Payments are optional and narrow. They exist only for `platform_sale` races, only when the seller has completed Stripe onboarding, and only if the pair chooses them. Every other mode is chat-only - the platform never touches money there, structurally.
3. Held funds, not "escrow". We use Stripe Connect with delayed transfers - buyer's money sits on the platform account until the bib transfer is confirmed. We avoid the word escrow in legal copy: regulated term, different thing.

## Money flow (Stripe Connect, separate charges & transfers)

Why separate charges & transfers (vs destination charges): full control of when the seller is paid (after confirmation), trivial refunds while funds are still on the platform balance, and the buyer charge is decoupled from the seller's account state. Both sides are EU, so the same-region constraint is satisfied.

- Sellers onboard as Connect Express accounts (Account Links) with the `transfers` capability. Stripe runs KYC - we never hold identity documents, only the account id and its status (from `account.updated` webhooks).
- Buyers pay a PaymentIntent on the platform account (with a `transfer_group`); Stripe.js collects the card, SCA/3DS happens automatically.

```mermaid
sequenceDiagram
    participant B as Buyer
    participant API as Go API
    participant S as Stripe
    participant SL as Seller
    B->>API: POST /orders (Idempotency-Key)
    API->>API: race=platform_sale? seller onboarded? listing->reserved (30min TTL)
    API->>S: Create PaymentIntent (transfer_group)
    B->>S: confirm card (Stripe.js, SCA)
    S-->>API: webhook payment_intent.succeeded
    API->>API: order -> paid_held  (funds on platform balance)
    SL->>API: mark bib transferred
    B->>API: confirm received (or auto-release after N days)
    API->>S: Create Transfer -> seller account
    S->>SL: payout (standard schedule)
```

Failure paths: TTL/abort -> order `cancelled`, listing back to `active`. Problems before release -> `refunded` (Stripe Refund of the item amount - easy, funds never left the platform; the processing fee is non-refundable, D31). Buyer dispute after seller marks transferred -> `disputed`, manual resolution. Chargebacks while funds are held are low-risk (nothing to claw back from the seller); post-release chargebacks -> transfer reversal attempt + ToS clause. Full state machine: [DATA_MODEL.md](DATA_MODEL.md#order-state-machine). Idempotency and every failure path in depth: [Idempotency & failure handling](#idempotency--failure-handling-m6-design) below + the diagrams in [`diagrams/`](diagrams/).

## Fee reality & pass-through

Stripe is not free even when we are: EU consumer cards run ≈ 1.5% + €0.25 per charge, plus Connect payout-leg fees (≈ 0.25% + €0.10 per payout) and a possible per-active-account fee - verify current pricing at implementation time; treat all rates as config constants.

With pass-through ON (decided - D1), the buyer sees an itemized "payment processing" line and the gross-up formula keeps the platform at net ≈ €0:

```
total = (item_price + fixed_fees) / (1 − percentage_fees)
```

Example at 1.5% + €0.25: €50.00 bib -> buyer pays €51.02, seller receives €50.00, platform nets ~€0 (± cents; payout-leg fees either folded into the formula or absorbed - [Decision 1](https://github.com/leonfullxr/bibseller/issues/3)).

Refunds (decided - D31, resolves issue #3's Decision 4): Stripe does not return
the original processing fee on refunds, and the buyer bears that - a refund
returns `item_amount_cents`, never `total_amount_cents`. Consequences we
accept and handle explicitly:

- Checkout must disclose it where the fee is itemized ("payment processing,
  non-refundable") - the buyer learns this before paying, not at refund time.
- Chargebacks are the exception by construction: Stripe returns the full
  charge to the card and debits us, so a disputed charge cannot withhold the
  fee. That asymmetry is another reason the dispute flow routes people to
  refunds first.
- An admin can still issue a full refund manually (Stripe dashboard) as a
  goodwill override; the product default stays item-amount.
- Counsel review flag (M7 #11): EU consumer-rights rules on withheld fees can
  depend on who cancels and why; the ToS wording needs counsel sign-off
  before payments GA.

PCI: Stripe.js/Elements only - card data never touches our servers (SAQ-A). Webhooks: signature-verified, idempotent via the `stripe_events` table, replay-safe.

## Idempotency & failure handling (M6 design)

Decided 2026-07-10 (D32). Diagrams: [`diagrams/payments-idempotency.drawio`](diagrams/payments-idempotency.drawio) (the three layers), [`diagrams/payments-money-flow.drawio`](diagrams/payments-money-flow.drawio) (happy path + failure overlays).

### The three idempotency layers

Money moves across three unreliable boundaries; each gets its own mechanism.
No generic idempotency-key table - the domain provides natural keys.

**1. Browser -> API.** `POST /orders` is idempotent through the
`orders_one_live_per_listing` partial unique index plus create semantics:
a repeat create for the same listing returns the *existing* live order when
the caller is its buyer (double-click and network-retry safe) and 409s for
anyone else. State transitions (`mark transferred`, `confirm received`,
`cancel`) are idempotent through the guarded update
(`UPDATE orders SET state=$to WHERE id=$1 AND state=$from`): a replay whose
target state already holds returns **200 with the current state** - "already
done" is success; only a *conflicting* current state is a 409.

**2. API -> Stripe (the critical layer).** Every mutating Stripe call carries
a deterministic idempotency key derived from the order id - `pi:{order_id}`,
`transfer:{order_id}`, `refund:{order_id}`. One PaymentIntent per order:
Stripe PIs retry failed card attempts natively, so a decline never mints a
new PI. Combined with invariant 4 (state row commits first, Stripe call
second) this makes every crash window safe:

- crash after commit, before the Stripe call -> the reconciler retries with
  the same key -> the object is created exactly once;
- crash after the call, before we store `stripe_transfer_id` -> the retry
  returns the *same* object instead of paying the seller twice.

Caveat: Stripe idempotency keys expire after ~24h, so keys are the fast path
only. Every object we create carries `metadata.order_id` (+ the
`transfer_group`), and reconciliation of older stragglers deduplicates by
metadata lookup, never by key alone.

Never call Stripe before committing the state row. With commit-first the
worst case is "state says completed, transfer missing" - the reconciler
creates it, the key prevents doubles. The reverse order can move money with
no state row demanding it.

**3. Stripe -> API (webhooks).** The handler's first statement is
`INSERT INTO stripe_events ... ON CONFLICT (id) DO NOTHING`; zero rows means
replay -> 200 immediately. Otherwise the event is processed **in the same
transaction** as that insert - a processing failure rolls back the event row
too, and Stripe's retry gets a clean second attempt. At-least-once delivery
becomes exactly-once effect. Two hard rules: verify `amount`/`currency`
against the order row before honoring a success event (signatures prove the
sender, not our expectations), and treat event *order* as untrusted -
Stripe does not guarantee ordering; the from-state guards make wrong-order
events harmless no-ops.

### Failure taxonomy

| # | Failure | Handling |
|---|---|---|
| 1 | Card declined / SCA abandoned | Order stays `pending_payment`; same PI retryable in the UI until the reservation TTL; each attempt appends `order_events` (`payment_failed` + decline reason). `payment_intent.payment_failed` changes no state. |
| 2 | TTL vs. late success race | The TTL job **cancels the PaymentIntent at Stripe first**, then flips the order to `cancelled`. A canceled PI cannot succeed, so the race collapses: cancel wins -> clean cancellation; Stripe answers "already succeeded" -> transition to `paid_held` instead (the money is real). True late success is reachable only through a crash mid-job -> auto-refund (item amount), `order_events` trail. Never a silently kept charge. |
| 3 | Transfer fails at release (seller account restricted since onboarding) | Funds are still on the platform balance - nothing is lost. No new state: `completed` commits, the reconciler retries the Transfer (same key), and alerts after N attempts for manual action. Revisit with a `release_failed` state only if it turns out common. |
| 4 | Chargeback while held | `charge.dispute.created` -> `disputed`. Low risk: nothing was paid out. Evidence pack = `order_events` + the chat confirmation trail. |
| 5 | Chargeback after release | Transfer reversal attempt + ToS clause; what reversal does not recover, the platform eats at beta volumes. Fee cannot be withheld on chargebacks (Stripe returns the full charge). |
| 6 | Webhook endpoint down | Stripe retries for ~3 days; the reconciler is the backstop; alert if `stripe_events` stays silent for hours while orders sit non-terminal. |
| 7 | Any crash window | Covered by the reconciler sweep (below). |

### The reconciler

One `internal/platform/batchjob` job is the backstop for every partial
failure at once: sweep orders in non-terminal states older than a few
minutes, cross-check the PaymentIntent / Transfer / Refund at Stripe via
`metadata.order_id`, and drive the same guarded state transitions as
`actor: 'system:reconciler'`. It follows the identical code path as webhooks
(state machine + `order_events`), so reconciliation can never disagree with
the live flow. Deadline duties (reservation TTL, 3-day auto-release,
refund-if-race-date-passed) ride the same job schedule.

### Platform balance float

Zero commission means no fee revenue buffering refunds and chargebacks: the
platform account eats payout-leg fees, chargeback debits, and goodwill
refunds. Mitigation: keep a small manual float on the platform Stripe
balance (~ €200 at beta) and alert on low balance. Honest line item for
investor material.

## Liability posture by policy mode

| | `platform_sale` | `official_only` | `connect_only` / `unknown` |
|---|---|---|---|
| Platform role | Facilitates payment, holds funds | Venue only - transfer happens via the race's own process | Venue only |
| Money through us | yes | never | never |
| Evidence trail | Full `order_events` log | Link-out UI to official procedure | Recorded buyer acknowledgment (`policy_acks`) before first contact |
| Required source | `policy_source_url` (DB-enforced) | `official_transfer_url` (DB-enforced) | - |

This is the strongest structural defense available: for races that restrict transfers, the platform can demonstrate it cannot process a payment (no code path), warned both parties, and recorded their acknowledgment. ToS must still state the responsibility split explicitly ([M7 #11](https://github.com/leonfullxr/bibseller/issues/11)).

## GDPR

- Roles: we are the controller; processors include hosting, email provider, object storage (DPAs required - keep a list). Stripe acts as an independent controller for its KYC/fraud processing.
- Lawful bases: contract (accounts, listings, chat, orders), legal obligation (financial records), legitimate interest (policy acks, abuse prevention, audit log).
- Data minimization: no tracking cookies (session cookie only -> no consent banner), no analytics with PII in v1, no PII in logs, retention schedule enforced by jobs ([DATA_MODEL.md -> Retention](DATA_MODEL.md#retention-gdpr-minimization)).
- Rights: export (JSON) and delete (anonymization with documented carve-outs for financial records) - both shipped in [M7 #11](https://github.com/leonfullxr/bibseller/issues/11).
- Residency: EU hosting (Hetzner Falkenstein recommended); storage/email providers chosen with EU regions or SCCs ([Decisions 6 & 7](https://github.com/leonfullxr/bibseller/issues/3)).
- Records of processing: one honest markdown doc, updated as features land. Breach process: note who to notify and the 72h clock.

## DSA (Digital Services Act)

We are a hosting service / online platform. Micro-enterprise status (<50 staff, <€10M) exempts us from the heavier platform-tier obligations (transparency reports, internal complaint systems at scale), but hosting-tier duties still apply:

- Notice-and-action: anyone (including non-users) can report a listing/message; reports land in a moderation queue; action is taken and logged.
- Statement of reasons: content removal notifies the author with the reason (template + `audit_log` row).
- Single point of contact published (contact/imprint page).
- Clear ToS, including the per-policy-mode responsibility split and a repeat-infringer policy.
- Trader traceability (Art. 30): v1 sellers are consumers (C2C), so trader KYB duties shouldn't bite - but define a threshold (e.g., a seller with many listings/year) at which we re-evaluate. Flag for counsel.

## DAC7 (platform tax reporting)

DAC7 obligations attach to facilitating relevant sales - not to profiting from them. Zero commission does not exempt us. Once payments GA, the platform likely qualifies as a reporting platform operator for sellers above the de-minimis threshold (>=30 sales or >= €2,000/year per seller - below that, sellers are excluded).

- We already capture what reporting needs: seller country (`users.country`), per-seller yearly totals (via `orders`).
- Action (gates payments GA): confirm scope + registration with a tax advisor. Most beta sellers will sell 1-2 bibs/year, far under threshold - but the obligation assessment must be done, not assumed. Tracked in [M7 #11](https://github.com/leonfullxr/bibseller/issues/11).

## Compliance gates

| Gate | Must be true |
|---|---|
| Public beta (chat-only) | ToS + Privacy published - contact point page - report flow works - GDPR export/delete shipped - EU hosting |
| Payments GA | Everything above - Stripe flows + webhooks battle-tested in sandbox - refund/dispute runbook - DAC7 memo signed off - counsel has reviewed ToS responsibility split & disclaimers |
