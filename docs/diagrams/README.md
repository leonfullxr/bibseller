# Diagrams

draw.io sources (open with [app.diagrams.net](https://app.diagrams.net) or the
VS Code extension) + committed PNG exports in [`exports/`](exports/) so they
render on GitHub and drop straight into decks. After editing a `.drawio`,
re-export: `drawio -x -f png -s 2 -o exports/ <file>.drawio`.

| Diagram | Shows |
| --- | --- |
| [architecture-overview](exports/architecture-overview.png) | The whole system: Cloudflare -> tunnel -> Caddy -> web/api, Postgres/MinIO, email relay, Stripe, the batchjob - and the transfer_policy gate that decides where money may flow |
| [payments-money-flow](exports/payments-money-flow.png) | Checkout sequence (separate charges & transfers): reserve -> pay -> held -> transfer confirmed -> release, plus the failure overlays (declined card, TTL race, dispute, refund) |
| [order-state-machine](exports/order-state-machine.png) | The seven order states, who moves them, and what each transition does with the money |
| [payments-idempotency](exports/payments-idempotency.png) | The three idempotency layers (browser->API, API->Stripe, Stripe->API) and the reconciler that backstops every crash window (D32) |
| [seller-onboarding](exports/seller-onboarding.png) | Stripe Connect Express lifecycle: none -> started -> active <-> restricted, what each state gates, and what happens to in-flight money (D34) |
| [checkout-policy-modes](exports/checkout-policy-modes.png) | What a buyer can do in each of the four transfer_policy modes - the safety architecture in one picture |
| [dispute-runbook](exports/dispute-runbook.png) | Dispute intake -> evidence -> resolution paths, with the buyer-favoring default and the 48h/7d SLA (D33) |

Design narrative: [PAYMENTS_AND_COMPLIANCE.md](../PAYMENTS_AND_COMPLIANCE.md)
(money flow, idempotency & failure handling) and the order state machine in
[DATA_MODEL.md](../DATA_MODEL.md#order-state-machine).
