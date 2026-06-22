import { apiGet } from '$lib/api/server';
import type { Page, RaceSummary } from '$lib/api/types';
import { todayISO } from '$lib/format';
import type { PageServerLoad } from './$types';

const PASSTHROUGH = ['country', 'sport', 'policy', 'q', 'cursor'] as const;

export const load: PageServerLoad = async ({ url, fetch, setHeaders, locals }) => {
	const params = new URLSearchParams();
	for (const key of PASSTHROUGH) {
		const v = url.searchParams.get(key);
		if (v) params.set(key, v);
	}
	// Hide past races from the default browse view.
	params.set('date_from', todayISO());
	params.set('limit', '24');

	const data = await apiGet<Page<RaceSummary>>(`/api/v1/races?${params}`, fetch);

	// Gate the cache header on auth: the page HTML embeds the layout nav (the
	// signed-in user's name, inbox, log out), so a signed-in response must never
	// sit in a shared cache. Anonymous browse stays publicly cacheable (CONTEXT
	// M3 cache note - mirrors listings/[id]).
	setHeaders({ 'cache-control': locals.user ? 'private, no-store' : 'public, max-age=60' });
	return {
		races: data.items,
		nextCursor: data.next_cursor,
		filters: {
			country: url.searchParams.get('country') ?? '',
			sport: url.searchParams.get('sport') ?? '',
			policy: url.searchParams.get('policy') ?? '',
			q: url.searchParams.get('q') ?? ''
		}
	};
};
