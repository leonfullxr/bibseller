import { apiGet } from '$lib/api/server';
import type { Page, RaceSummary } from '$lib/api/types';
import { todayISO } from '$lib/format';
import type { PageServerLoad } from './$types';

const PASSTHROUGH = ['country', 'sport', 'policy', 'q', 'distance', 'cursor'] as const;

// Date params are user-editable URL state; drop malformed values here so a
// hand-mangled query string can never 400 the whole page at the API.
const ISO_DATE = /^\d{4}-\d{2}-\d{2}$/;

type MapCounts = {
	countries: Record<string, number>;
	cities: {
		city: string;
		country: string;
		count: number;
		races: { name: string; slug: string }[];
	}[];
};

export const load: PageServerLoad = async ({ url, fetch, setHeaders, locals }) => {
	const params = new URLSearchParams();
	for (const key of PASSTHROUGH) {
		const v = url.searchParams.get(key);
		if (v) params.set(key, v);
	}
	// Hide past races from the browse view: the user's date_from is honored
	// only from today forward (ISO date strings compare lexicographically), so
	// the server-side floor from the pre-filter days still holds.
	const today = todayISO();
	const fromRaw = url.searchParams.get('date_from') ?? '';
	const dateFrom = ISO_DATE.test(fromRaw) && fromRaw > today ? fromRaw : '';
	params.set('date_from', dateFrom || today);
	const toRaw = url.searchParams.get('date_to') ?? '';
	const dateTo = ISO_DATE.test(toRaw) ? toRaw : '';
	if (dateTo) params.set('date_to', dateTo);
	params.set('limit', '24');

	const [data, map] = await Promise.all([
		apiGet<Page<RaceSummary>>(`/api/v1/races?${params}`, fetch),
		// Per-country / per-city counts for the decorative map, aggregated
		// server-side (one GROUP BY, not bounded by a page of races). The map is
		// decorative, so never let its fetch fail the page: degrade to no map
		// (+page.svelte only renders the map when countryCounts is non-empty).
		apiGet<MapCounts>(`/api/v1/races/map-counts`, fetch).catch(() => null)
	]);

	// Gate the cache header on auth: the page HTML embeds the layout nav (the
	// signed-in user's name, inbox, log out), so a signed-in response must never
	// sit in a shared cache. Anonymous browse stays publicly cacheable (CONTEXT
	// M3 cache note - mirrors listings/[id]).
	setHeaders({ 'cache-control': locals.user ? 'private, no-store' : 'public, max-age=60' });
	return {
		races: data.items,
		nextCursor: data.next_cursor,
		countryCounts: map?.countries ?? {},
		cities: map?.cities ?? [],
		filters: {
			country: url.searchParams.get('country') ?? '',
			sport: url.searchParams.get('sport') ?? '',
			policy: url.searchParams.get('policy') ?? '',
			q: url.searchParams.get('q') ?? '',
			distance: url.searchParams.get('distance') ?? '',
			date_from: dateFrom,
			date_to: dateTo
		}
	};
};
