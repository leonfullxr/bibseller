import { describe, expect, it } from 'vitest';
import {
	cityCoords,
	CITY_COORDS,
	COUNTRY_VIEWBOX,
	EUROPE_VIEWBOX,
	fitViewBox,
	project
} from './cities';

// races.city -> the country whose viewBox the marker must land inside. Guards
// both the projection constants and the per-city coordinates: a wrong dot lands
// outside its country's box and fails here.
const CITY_COUNTRY: Record<string, string> = {
	Munich: 'DE',
	Granada: 'ES',
	Brussels: 'BE',
	'Riva del Garda': 'IT',
	Alcúdia: 'ES',
	Valencia: 'ES',
	Paris: 'FR',
	Amsterdam: 'NL',
	Vienna: 'AT',
	Berlin: 'DE',
	Porto: 'PT',
	Frankfurt: 'DE',
	Milan: 'IT',
	Rotterdam: 'NL',
	Sevilla: 'ES',
	Lisbon: 'PT',
	Kraków: 'PL',
	Bilbao: 'ES',
	Madrid: 'ES'
};

describe('city projection', () => {
	it('projects every seed city inside its country frame', () => {
		for (const [city, [lat, lng]] of Object.entries(CITY_COORDS)) {
			const cc = CITY_COUNTRY[city];
			const [x0, y0, w, h] = COUNTRY_VIEWBOX[cc].split(' ').map(Number);
			const [x, y] = project(lat, lng);
			expect(x, `${city} x`).toBeGreaterThanOrEqual(x0);
			expect(x, `${city} x`).toBeLessThanOrEqual(x0 + w);
			expect(y, `${city} y`).toBeGreaterThanOrEqual(y0);
			expect(y, `${city} y`).toBeLessThanOrEqual(y0 + h);
		}
	});

	it('pads every country frame to a constant aspect ratio, centered', () => {
		// Keeps the map a constant on-page size; the original frame stays enclosed.
		const [, , ew, eh] = EUROPE_VIEWBOX.split(' ').map(Number);
		const target = ew / eh;
		for (const vb of [EUROPE_VIEWBOX, ...Object.values(COUNTRY_VIEWBOX)]) {
			const [x, y, w, h] = fitViewBox(vb);
			expect(w / h).toBeCloseTo(target, 5);
			const [ox, oy, ow, oh] = vb.split(' ').map(Number);
			expect(x).toBeLessThanOrEqual(ox + 1e-6);
			expect(y).toBeLessThanOrEqual(oy + 1e-6);
			expect(x + w).toBeGreaterThanOrEqual(ox + ow - 1e-6);
			expect(y + h).toBeGreaterThanOrEqual(oy + oh - 1e-6);
		}
	});

	it('resolves coordinates ignoring accents and case', () => {
		// Seed files spell cities inconsistently (e.g. "Alcudia" vs "Alcúdia").
		expect(cityCoords('Alcudia')).toEqual(CITY_COORDS['Alcúdia']);
		expect(cityCoords('krakow')).toEqual(CITY_COORDS['Kraków']);
		expect(cityCoords('Nowhere')).toBeUndefined();
	});
});
