import { describe, expect, it } from 'vitest';
import { mapQuery, mapCityVisible } from './raceMapLinks';

const none = { sport: '', policy: '', q: '' };

describe('mapQuery', () => {
	it('keeps the active sport/policy/q when setting a country (#90)', () => {
		const f = { sport: 'triathlon', policy: 'platform_sale', q: 'mara' };
		expect(mapQuery(f, { country: 'FR' })).toBe(
			'?sport=triathlon&policy=platform_sale&q=mara&country=FR'
		);
	});

	it('overrides q with the city on a city click but keeps sport/policy (#90)', () => {
		const f = { sport: 'running', policy: '', q: 'old' };
		expect(mapQuery(f, { country: 'ES', q: 'Granada' })).toBe(
			'?sport=running&q=Granada&country=ES'
		);
	});

	it('url-encodes override values', () => {
		expect(mapQuery(none, { country: 'IT', q: 'San Remo' })).toBe('?country=IT&q=San+Remo');
	});

	it('clears a param when the override value is empty (the "all of Europe" link)', () => {
		const f = { sport: 'cycling', policy: '', q: '' };
		expect(mapQuery(f, { country: '' })).toBe('?sport=cycling');
	});

	it('returns an empty suffix when nothing is set', () => {
		expect(mapQuery(none, { country: '' })).toBe('');
	});
});

describe('mapCityVisible', () => {
	it('shows every city when no country filter is active', () => {
		expect(mapCityVisible('FR', '')).toBe(true);
		expect(mapCityVisible('DE', '')).toBe(true);
	});

	it('shows only the active country, in-set or off-map (#91)', () => {
		expect(mapCityVisible('GB', 'GB')).toBe(true);
		expect(mapCityVisible('FR', 'GB')).toBe(false);
		expect(mapCityVisible('FR', 'FR')).toBe(true);
		expect(mapCityVisible('DE', 'FR')).toBe(false);
	});
});
