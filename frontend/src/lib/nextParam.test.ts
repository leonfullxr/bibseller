import { describe, expect, it } from 'vitest';
import { safeNext } from './nextParam';

describe('safeNext', () => {
	it('accepts internal paths', () => {
		expect(safeNext('/listings/abc')).toBe('/listings/abc');
		expect(safeNext('/es/races?country=ES')).toBe('/es/races?country=ES');
	});

	it('falls back to / for anything else', () => {
		expect(safeNext(null)).toBe('/');
		expect(safeNext('')).toBe('/');
		expect(safeNext('https://evil.example')).toBe('/');
		expect(safeNext('//evil.example')).toBe('/');
		expect(safeNext('/\\evil.example')).toBe('/');
		expect(safeNext('javascript:alert(1)')).toBe('/');
	});
});
