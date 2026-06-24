// City coordinates and the projection that places them on the static Europe SVG
// base map (flekschas/simple-world-map). No map library: project() is a linear
// fit calibrated against the base map's country positions. The fit holds across
// ~35-55N (max residual ~3px), which is the band every race city sits in.
//
// ponytail: static lookup for today's curated race catalog (~20 admin-seeded
// races). When races become user-submitted, move latitude/longitude onto the
// races table (DATA_MODEL) and feed real coords in instead of this map.

// lng -> svg x, lat -> svg y. Constants from the least-squares fit; see the
// self-check in cities.test.ts.
export function project(lat: number, lng: number): [number, number] {
	return [2.2688 * lng + 406.3966, -3.1451 * lat + 552.8604];
}

// Default frame: western/central Europe, wider than tall, North Africa and the
// far-flung islands (e.g. the Canaries) fall outside the window and so don't
// render - no SVG surgery needed.
export const EUROPE_VIEWBOX = '370 358 105 87';

// Click-to-zoom: padded bounding box of each filterable country's mainland.
export const COUNTRY_VIEWBOX: Record<string, string> = {
	AT: '427.4 398.3 18.6 9.6',
	BE: '412.1 390.7 10.3 7.3',
	DE: '417.0 377.0 25.5 31.6',
	ES: '382.2 411.4 41.7 29.6',
	FR: '392.4 389.1 36.8 33.6',
	IT: '418.2 400.8 34.2 36.3',
	NL: '413.0 383.0 10.7 11.6',
	PL: '434.1 376.4 27.5 24.7',
	PT: '384.4 419.2 8.1 19.1'
};

// Keys must match races.city exactly (accents included). A city with no entry
// simply renders no marker.
export const CITY_COORDS: Record<string, [number, number]> = {
	Munich: [48.14, 11.58],
	Granada: [37.18, -3.6],
	Brussels: [50.85, 4.35],
	'Riva del Garda': [45.88, 10.84],
	Alcúdia: [39.85, 3.12],
	Valencia: [39.47, -0.38],
	Paris: [48.86, 2.35],
	Amsterdam: [52.37, 4.9],
	Vienna: [48.21, 16.37],
	Berlin: [52.52, 13.4],
	Porto: [41.15, -8.61],
	Frankfurt: [50.11, 8.68],
	Milan: [45.46, 9.19],
	Rotterdam: [51.92, 4.48],
	Sevilla: [37.39, -5.99],
	Lisbon: [38.72, -9.14],
	Kraków: [50.06, 19.94],
	Bilbao: [43.26, -2.93],
	Madrid: [40.42, -3.7]
};
