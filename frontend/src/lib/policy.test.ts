import { describe, expect, it } from 'vitest';
import type { TransferPolicy } from '$lib/api/types';
import { policyDisclaimer, policyLabel, policyView, requiresAck } from './policy';

const policies: TransferPolicy[] = ['platform_sale', 'official_only', 'connect_only', 'unknown'];

describe('policy view', () => {
	it('defines a view, label, and disclaimer for every mode', () => {
		for (const p of policies) {
			expect(policyView[p]).toBeDefined();
			expect(policyLabel[p]).toBeTruthy();
			expect(policyDisclaimer[p].title).toBeTruthy();
			expect(policyDisclaimer[p].body).toBeTruthy();
		}
	});

	it('offers the buy affordance only for platform_sale', () => {
		// Mirrors the server-side invariant: no buy path exists off platform_sale.
		for (const p of policies) {
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
