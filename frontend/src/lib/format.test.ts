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
	// Expected strings are derived via the same Intl.DateTimeFormat call
	// formatWhen uses, so these assert the today/older branching and options,
	// not a hard-coded ICU rendering (fragile across Node/ICU versions).
	it('shows a time for a timestamp from today (UTC)', () => {
		const iso = `${todayISO()}T09:30:00Z`;
		const expected = new Intl.DateTimeFormat('en', { timeStyle: 'short', timeZone: 'UTC' }).format(
			new Date(iso)
		);
		expect(formatWhen(iso)).toBe(expected);
	});
	it('shows a medium date for older timestamps', () => {
		const iso = '2024-12-06T09:30:00Z';
		const expected = new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeZone: 'UTC' }).format(
			new Date(iso)
		);
		expect(formatWhen(iso)).toBe(expected);
	});
	it('localizes', () => {
		const iso = '2024-12-06T09:30:00Z';
		const expected = new Intl.DateTimeFormat('es', { dateStyle: 'medium', timeZone: 'UTC' }).format(
			new Date(iso)
		);
		expect(formatWhen(iso, 'es')).toBe(expected);
	});
});

describe('todayISO', () => {
	it('is a YYYY-MM-DD string', () => {
		expect(todayISO()).toMatch(/^\d{4}-\d{2}-\d{2}$/);
	});
});
