// The policy view: the frontend's single derivation of a race's transfer_policy
// into presentation facts (docs/CONTEXT.md -> Language). Components read this
// instead of branching on the policy string; the server still enforces what
// money may flow (invariant 7) - this only governs what renders.
import type { TransferPolicy } from '$lib/api/types';

// tone collapses the four modes into three visual/disclaimer groupings
// (connect_only and unknown share "restricted"). primaryAction is the listing
// CTA affordance; "buy" exists only for platform_sale, mirroring the API.
export type PolicyTone = 'sale' | 'official' | 'restricted';
export type PrimaryAction = 'buy' | 'official' | null;

export interface PolicyView {
	tone: PolicyTone;
	primaryAction: PrimaryAction;
}

export const policyView: Record<TransferPolicy, PolicyView> = {
	platform_sale: { tone: 'sale', primaryAction: 'buy' },
	official_only: { tone: 'official', primaryAction: 'official' },
	connect_only: { tone: 'restricted', primaryAction: null },
	unknown: { tone: 'restricted', primaryAction: null }
};

// Whether contacting a seller for this race needs the buyer to acknowledge the
// venue-only terms first: connect_only and unknown. An explicit policy check,
// not derived from the presentational `tone`, so a UI tweak can never shift this
// gate. The server enforces it too - the chat API rejects a first message
// without the recorded ack - this only drives whether the UI shows the gate.
export function requiresAck(policy: TransferPolicy): boolean {
	return policy === 'connect_only' || policy === 'unknown';
}

// The four transfer policies, for iterating the modes (the races filter, tests).
// The display words now live in the i18n dictionary (docs/CONTEXT.md -> Language,
// D17): $lib/i18n/en.ts keys policy.label.* and policy.disclaimer.*. This module
// keeps only the policy facts, read via the typed keys in PolicyBadge/PolicyCallout.
export const transferPolicies: TransferPolicy[] = [
	'platform_sale',
	'official_only',
	'connect_only',
	'unknown'
];
