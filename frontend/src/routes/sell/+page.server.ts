import { redirect } from '@sveltejs/kit';
import { apiGet } from '$lib/api/server';
import type { Page, RaceSummary } from '$lib/api/types';
import { todayISO } from '$lib/format';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals, fetch }) => {
	if (!locals.user) redirect(303, '/login');

	const q = url.searchParams.get('q') ?? '';
	const params = new URLSearchParams();
	if (q) params.set('q', q);
	// You can only list a bib for a race that hasn't happened yet.
	params.set('date_from', todayISO());
	params.set('limit', '24');

	const data = await apiGet<Page<RaceSummary>>(`/api/v1/races?${params}`, fetch);
	return { races: data.items, q, verified: locals.user.email_verified };
};
