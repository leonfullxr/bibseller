import { describe, expect, it } from 'vitest';
import { apiErrorKey } from './errors';

describe('apiErrorKey', () => {
	it('maps a known code to its apiError key', () => {
		expect(apiErrorKey('email_taken')).toBe('apiError.email_taken');
		expect(apiErrorKey('race_past')).toBe('apiError.race_past');
	});

	it('falls back to apiError.unknown for unknown or empty codes', () => {
		expect(apiErrorKey('totally_made_up')).toBe('apiError.unknown');
		expect(apiErrorKey('')).toBe('apiError.unknown');
		expect(apiErrorKey(undefined)).toBe('apiError.unknown');
		expect(apiErrorKey(null)).toBe('apiError.unknown');
	});
});
