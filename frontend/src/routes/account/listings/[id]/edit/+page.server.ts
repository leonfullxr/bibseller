import { fail, redirect } from '@sveltejs/kit';
import { apiFetch, apiGet } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import type { ListingDetail } from '$lib/api/types';
import { parseListingPrice } from '$lib/listing';
import { sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, fetch }) => {
	if (!locals.user) redirect(303, '/login');
	const listing = await apiGet<ListingDetail>(`/api/v1/listings/${params.id}`, fetch);
	return { listing };
};

export const actions: Actions = {
	default: async ({ params, request, cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');
		const t = createTranslator(locals.locale);

		const form = await request.formData();
		const description = String(form.get('description') ?? '').trim();
		const parsed = parseListingPrice(
			String(form.get('price') ?? ''),
			String(form.get('original_price') ?? '')
		);
		if (!parsed.ok) {
			return fail(400, { error: t(parsed.key), values: snapshot(form) });
		}

		let res: Response;
		try {
			res = await apiFetch(`/api/v1/listings/${params.id}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json', ...sessionHeader(cookies) },
				body: JSON.stringify({
					price_cents: parsed.value.priceCents,
					original_price_cents: parsed.value.originalCents,
					description: description || null
				})
			});
		} catch {
			return fail(502, { error: t('apiError.unreachable'), values: snapshot(form) });
		}

		if (res.status === 403) {
			return fail(403, { error: t('formError.editOwnOnly'), values: snapshot(form) });
		}
		if (res.status === 409) {
			return fail(409, { error: t('formError.editNotActive'), values: snapshot(form) });
		}
		if (!res.ok) {
			return fail(res.status >= 500 ? 502 : res.status, {
				error: t('formError.editFailed'),
				values: snapshot(form)
			});
		}

		redirect(303, '/account/listings');
	}
};

function snapshot(form: FormData) {
	return {
		price: String(form.get('price') ?? ''),
		original_price: String(form.get('original_price') ?? ''),
		description: String(form.get('description') ?? '')
	};
}
