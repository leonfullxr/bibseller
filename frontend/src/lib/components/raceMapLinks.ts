// Pure helpers for the /races map's links and dot visibility. Extracted from
// RaceMap.svelte so they can be unit-tested without a DOM (the project has no
// component-test runner; see vite.config.ts).

export type MapFilters = { sport: string; policy: string; q: string };

// Build the query-string suffix ("?sport=...&country=FR", or "" when empty) for a
// /races map link, preserving the active sport/policy/q filters and applying the
// given overrides: country on a country click, country + q on a city click, or
// country: '' to clear it for the "all of Europe" link. An empty override value
// clears that param. Returned as a suffix to append to the resolved base href so
// the route still flows through resolve()/link() (svelte/no-navigation rule, i18n
// locale prefix). Fixes #90 (map links used to rebuild the URL from scratch,
// silently discarding sport/policy/q).
export function mapQuery(filters: MapFilters, overrides: Record<string, string>): string {
	const p = new URLSearchParams();
	if (filters.sport) p.set('sport', filters.sport);
	if (filters.policy) p.set('policy', filters.policy);
	if (filters.q) p.set('q', filters.q);
	for (const [k, v] of Object.entries(overrides)) {
		if (v) p.set(k, v);
		else p.delete(k);
	}
	const qs = p.toString();
	return qs ? `?${qs}` : '';
}

// A city dot is shown when no country filter is active, or it belongs to the
// active country. Keyed on the filter being active (not on the map having a
// viewBox to zoom to) so an off-map ISO code (e.g. GB) still shows only that
// country's dots instead of all of Europe's. Fixes #91.
export function mapCityVisible(cityCountry: string, activeCountry: string): boolean {
	return !activeCountry || cityCountry === activeCountry;
}
