import { fail, redirect } from '@sveltejs/kit';
import { apiFetch, apiGet } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import type { OwnedListing } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies, fetch }) => {
	if (!locals.user) redirect(303, '/login');

	const data = await apiGet<{ items: OwnedListing[] }>('/api/v1/me/listings', fetch, {
		headers: sessionHeader(cookies)
	});
	return { listings: data.items };
};

export const actions: Actions = {
	cancel: async ({ request, cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');
		const t = createTranslator(locals.locale);

		const id = String((await request.formData()).get('id') ?? '');
		if (id === '') {
			return fail(400, { error: t('formError.missingListingId') });
		}

		let res: Response;
		try {
			res = await apiFetch(`/api/v1/listings/${id}/cancel`, {
				method: 'POST',
				headers: sessionHeader(cookies)
			});
		} catch {
			return fail(502, { error: t('apiError.unreachable') });
		}
		// 409 means it was already not active - the end state the user wanted.
		if (!res.ok && res.status !== 409) {
			return fail(res.status >= 500 ? 502 : res.status, { error: t('formError.cancelFailed') });
		}
		return { cancelled: id };
	}
};
