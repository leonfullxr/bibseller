import { describe, expect, it } from 'vitest';
import { activeSection } from './sections';

describe('activeSection', () => {
	it('defaults to profile when the param is absent', () => {
		expect(activeSection(null)).toBe('profile');
	});

	it('returns a known section verbatim', () => {
		expect(activeSection('security')).toBe('security');
		expect(activeSection('account')).toBe('account');
	});

	it('falls back to profile on unknown values', () => {
		expect(activeSection('nonsense')).toBe('profile');
		expect(activeSection('')).toBe('profile');
	});
});
