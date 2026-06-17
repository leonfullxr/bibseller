import { error, fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { OwnedListing } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.user) redirect(303, '/login');

	let res: Response;
	try {
		res = await apiFetch('/api/v1/me/listings', { headers: sessionHeader(cookies) });
	} catch {
		error(502, 'The API is unreachable.');
	}
	if (!res.ok) error(502, 'Could not load your listings.');

	const data = (await res.json()) as { items: OwnedListing[] };
	return { listings: data.items };
};

export const actions: Actions = {
	cancel: async ({ request, cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');

		const id = String((await request.formData()).get('id') ?? '');
		let res: Response;
		try {
			res = await apiFetch(`/api/v1/listings/${id}/cancel`, {
				method: 'POST',
				headers: sessionHeader(cookies)
			});
		} catch {
			return fail(502, { error: 'The API is unreachable.' });
		}
		// 409 means it was already not active - the end state the user wanted.
		if (!res.ok && res.status !== 409) {
			return fail(res.status >= 500 ? 502 : res.status, { error: 'Could not cancel the listing.' });
		}
		return { cancelled: id };
	}
};
