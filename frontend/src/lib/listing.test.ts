import { describe, it, expect } from 'vitest';
import { parseListingPrice } from './listing';

describe('parseListingPrice', () => {
	it('parses whole and decimal euros to cents', () => {
		expect(parseListingPrice('45', '')).toEqual({
			ok: true,
			value: { priceCents: 4500, originalCents: null }
		});
		expect(parseListingPrice('45.50', '60')).toEqual({
			ok: true,
			value: { priceCents: 4550, originalCents: 6000 }
		});
	});

	it('treats empty fields as null (both optional)', () => {
		expect(parseListingPrice('', '')).toEqual({
			ok: true,
			value: { priceCents: null, originalCents: null }
		});
	});

	it('rejects non-numeric or negative amounts', () => {
		expect(parseListingPrice('abc', '').ok).toBe(false);
		expect(parseListingPrice('-5', '').ok).toBe(false);
	});

	it('enforces the D2 cap: price cannot exceed face value', () => {
		const r = parseListingPrice('70', '60');
		expect(r.ok).toBe(false);
		if (!r.ok) expect(r.error).toMatch(/face value/);
	});

	it('allows price equal to face value', () => {
		expect(parseListingPrice('60', '60')).toEqual({
			ok: true,
			value: { priceCents: 6000, originalCents: 6000 }
		});
	});
});
