import { apiFetch, apiGet } from '$lib/api/server';
import type { Page, RaceSummary } from '$lib/api/types';
import { todayISO } from '$lib/format';
import type { PageServerLoad } from './$types';

// The landing is the search: the same catalog passthrough as /races, trimmed
// to the quick filters the landing exposes.
const PASSTHROUGH = ['q', 'country', 'sport'] as const;

export const load: PageServerLoad = async ({ url, fetch }) => {
	const filters = {
		q: url.searchParams.get('q') ?? '',
		country: url.searchParams.get('country') ?? '',
		sport: url.searchParams.get('sport') ?? ''
	};
	let apiStatus: 'ok' | 'down' = 'down';
	let races: RaceSummary[] = [];
	let countryCounts: Record<string, number> = {};
	try {
		const res = await apiFetch('/api/healthz', { signal: AbortSignal.timeout(1500) });
		apiStatus = res.ok ? 'ok' : 'down';
		if (res.ok) {
			const params = new URLSearchParams();
			for (const key of PASSTHROUGH) {
				if (filters[key]) params.set(key, filters[key]);
			}
			params.set('date_from', todayISO());
			params.set('limit', '12');
			const [data, map] = await Promise.all([
				apiGet<Page<RaceSummary>>(`/api/v1/races?${params}`, fetch),
				// Country counts feed the quick-filter pills; they are decorative
				// enough that their fetch must never fail the landing.
				apiGet<{ countries: Record<string, number> }>('/api/v1/races/map-counts', fetch).catch(
					() => null
				)
			]);
			races = data.items;
			countryCounts = map?.countries ?? {};
		}
	} catch {
		// The landing page must render even with the API down.
	}
	return { apiStatus, races, countryCounts, filters };
};
