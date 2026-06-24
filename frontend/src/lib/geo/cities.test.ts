import { describe, expect, it } from 'vitest';
import { cityCoords, CITY_COORDS, COUNTRY_VIEWBOX, project } from './cities';

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

	it('resolves coordinates ignoring accents and case', () => {
		// Seed files spell cities inconsistently (e.g. "Alcudia" vs "Alcúdia").
		expect(cityCoords('Alcudia')).toEqual(CITY_COORDS['Alcúdia']);
		expect(cityCoords('krakow')).toEqual(CITY_COORDS['Kraków']);
		expect(cityCoords('Nowhere')).toBeUndefined();
	});
});
