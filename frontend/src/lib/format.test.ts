import { describe, expect, it } from 'vitest';
import { formatDate, formatPrice, formatWhen, todayISO } from './format';

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

describe('formatWhen', () => {
	it('shows a time for a timestamp from today (UTC)', () => {
		expect(formatWhen(`${todayISO()}T09:30:00Z`)).toBe('9:30 AM');
	});
	it('shows a medium date for older timestamps', () => {
		expect(formatWhen('2024-12-06T09:30:00Z')).toBe('Dec 6, 2024');
	});
	it('localizes', () => {
		expect(formatWhen('2024-12-06T09:30:00Z', 'es')).toBe('6 dic 2024');
	});
});

describe('todayISO', () => {
	it('is a YYYY-MM-DD string', () => {
		expect(todayISO()).toMatch(/^\d{4}-\d{2}-\d{2}$/);
	});
});
