import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { formatDate, formatPrice, formatTime, todayISO } from './format';

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

describe('formatTime', () => {
	const originalTZ = process.env.TZ;
	beforeEach(() => {
		process.env.TZ = 'America/New_York'; // fixed non-UTC zone so this isn't UTC by coincidence
	});
	afterEach(() => {
		// TZ='' + reassigning undefined would set the literal string "undefined";
		// delete it instead when it was unset, so a bad zone can't leak to later tests.
		if (originalTZ === undefined) delete process.env.TZ;
		else process.env.TZ = originalTZ;
	});

	it("formats in the process's local timezone, not UTC", () => {
		// 2026-12-06T18:30:00Z is 13:30 in America/New_York (UTC-5 in December).
		// Pin the locale and match \s (covers ICU's narrow no-break space before
		// the meridiem, which varies by ICU version) so this isn't CI-fragile.
		expect(formatTime('2026-12-06T18:30:00Z', 'en-US')).toMatch(/^1:30\sPM$/);
	});
});

describe('todayISO', () => {
	it('is a YYYY-MM-DD string', () => {
		expect(todayISO()).toMatch(/^\d{4}-\d{2}-\d{2}$/);
	});
});
