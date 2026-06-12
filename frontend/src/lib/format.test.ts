import { describe, expect, it } from 'vitest';
import { formatDate, formatPrice, todayISO } from './format';

describe('formatPrice', () => {
	it('drops cents for whole amounts', () => {
		expect(formatPrice(8500)).toBe('€85');
	});
	it('keeps cents otherwise', () => {
		expect(formatPrice(8550)).toBe('€85.50');
	});
	it('returns null for missing prices', () => {
		expect(formatPrice(null)).toBeNull();
		expect(formatPrice(undefined)).toBeNull();
	});
});

describe('formatDate', () => {
	it('formats without timezone drift', () => {
		expect(formatDate('2026-12-06')).toBe('December 6, 2026');
	});
});

describe('todayISO', () => {
	it('is a YYYY-MM-DD string', () => {
		expect(todayISO()).toMatch(/^\d{4}-\d{2}-\d{2}$/);
	});
});
