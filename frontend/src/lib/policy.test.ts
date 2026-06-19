import { describe, expect, it } from 'vitest';
import { policyView, requiresAck, transferPolicies } from './policy';

describe('policy view', () => {
	it('defines a view for every mode', () => {
		// Labels/disclaimers moved to the i18n dictionary (coverage in i18n/locale.test.ts).
		for (const p of transferPolicies) {
			expect(policyView[p]).toBeDefined();
		}
	});

	it('offers the buy affordance only for platform_sale', () => {
		// Mirrors the server-side invariant: no buy path exists off platform_sale.
		for (const p of transferPolicies) {
			expect(policyView[p].primaryAction === 'buy').toBe(p === 'platform_sale');
		}
	});

	it('maps each mode to its CTA affordance', () => {
		expect(policyView.platform_sale.primaryAction).toBe('buy');
		expect(policyView.official_only.primaryAction).toBe('official');
		expect(policyView.connect_only.primaryAction).toBeNull();
		expect(policyView.unknown.primaryAction).toBeNull();
	});

	it('collapses the four modes into three disclaimer tones', () => {
		expect(policyView.platform_sale.tone).toBe('sale');
		expect(policyView.official_only.tone).toBe('official');
		expect(policyView.connect_only.tone).toBe('restricted');
		expect(policyView.unknown.tone).toBe('restricted');
	});

	it('requires an acknowledgment only for the restricted modes', () => {
		// Mirrors the server-side ack gate: connect_only/unknown need it, the rest do not.
		expect(requiresAck('platform_sale')).toBe(false);
		expect(requiresAck('official_only')).toBe(false);
		expect(requiresAck('connect_only')).toBe(true);
		expect(requiresAck('unknown')).toBe(true);
	});
});
