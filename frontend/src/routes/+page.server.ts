import { apiFetch, apiGet } from '$lib/api/server';
import type { Page, RaceSummary } from '$lib/api/types';
import { todayISO } from '$lib/format';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	let apiStatus: 'ok' | 'down' = 'down';
	let upcoming: RaceSummary[] = [];
	try {
		const res = await apiFetch('/api/healthz', { signal: AbortSignal.timeout(1500) });
		apiStatus = res.ok ? 'ok' : 'down';
		if (res.ok) {
			const data = await apiGet<Page<RaceSummary>>(
				`/api/v1/races?date_from=${todayISO()}&limit=6`,
				fetch
			);
			upcoming = data.items;
		}
	} catch {
		// The landing page must render even with the API down.
	}
	return { apiStatus, upcoming };
};
