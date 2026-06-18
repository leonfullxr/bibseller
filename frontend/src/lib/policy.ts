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

// Words live apart from facts (the seed of the M8 i18n dictionary). Labels feed
// the badge and the races filter; one canonical string per mode kills drift.
export const policyLabel: Record<TransferPolicy, string> = {
	platform_sale: 'Resale allowed',
	official_only: 'Official transfer',
	connect_only: 'Chat only',
	unknown: 'Policy unverified'
};

export interface PolicyDisclaimer {
	title: string;
	body: string;
}

// Keyed by policy (not tone) so connect_only and unknown can diverge later
// without reshaping (docs/CONTEXT.md D10) - their bodies start identical.
export const policyDisclaimer: Record<TransferPolicy, PolicyDisclaimer> = {
	platform_sale: {
		title: 'This race allows bib resale.',
		body: 'Agree with the seller in chat, then pay securely through the platform - funds are held until the transfer is confirmed. Zero commission.'
	},
	official_only: {
		title: 'This race runs its own official name-change process.',
		body: 'Find each other and agree on the details here - the transfer itself (and any official fee) goes through the race organizer. The platform never handles money for this race.'
	},
	connect_only: {
		title: 'This race restricts bib transfers.',
		body: "The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything."
	},
	unknown: {
		title: 'Transfer policy not verified yet - treat this race as chat-only.',
		body: "The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything."
	}
};
