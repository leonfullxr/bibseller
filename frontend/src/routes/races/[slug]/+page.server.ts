import { apiGet } from '$lib/api/server';
import type { ListingSummary, Page, RaceDetail } from '$lib/api/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch, setHeaders, locals }) => {
	const [race, listings] = await Promise.all([
		apiGet<RaceDetail>(`/api/v1/races/${params.slug}`, fetch),
		apiGet<Page<ListingSummary>>(`/api/v1/races/${params.slug}/listings`, fetch)
	]);

	// Signed-in responses embed the user's nav, so they must not share a cache;
	// anonymous race detail stays publicly cacheable (CONTEXT M3 cache note).
	setHeaders({ 'cache-control': locals.user ? 'private, no-store' : 'public, max-age=60' });
	return { race, listings: listings.items };
};
