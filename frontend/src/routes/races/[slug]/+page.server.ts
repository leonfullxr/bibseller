import { apiGet } from '$lib/api/server';
import type { ListingSummary, Page, RaceDetail } from '$lib/api/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch, setHeaders }) => {
	const [race, listings] = await Promise.all([
		apiGet<RaceDetail>(`/api/v1/races/${params.slug}`, fetch),
		apiGet<Page<ListingSummary>>(`/api/v1/races/${params.slug}/listings`, fetch)
	]);

	setHeaders({ 'cache-control': 'public, max-age=60' });
	return { race, listings: listings.items };
};
