import { apiGet } from '$lib/api/server';
import type { ListingDetail } from '$lib/api/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch, setHeaders }) => {
	const listing = await apiGet<ListingDetail>(`/api/v1/listings/${params.id}`, fetch);
	setHeaders({ 'cache-control': 'public, max-age=60' });
	return { listing };
};
