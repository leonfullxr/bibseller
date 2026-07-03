import { describe, expect, it } from 'vitest';
import { pendingForm } from './forms.svelte';

describe('pendingForm', () => {
	it('flips busy around update()', async () => {
		const { busy, submit } = pendingForm();
		expect(busy.value).toBe(false);

		const after = submit({} as never);
		expect(busy.value).toBe(true);
		if (typeof after !== 'function') throw new Error('expected an enhance callback');

		let updated = false;
		await after({
			update: async () => {
				updated = true;
				expect(busy.value).toBe(true); // still pending while update runs
			}
		} as never);

		expect(updated).toBe(true);
		expect(busy.value).toBe(false);
	});
});
