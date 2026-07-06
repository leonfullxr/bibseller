import { error, fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import type { ListingDetail } from '$lib/api/types';
import { listingFormSnapshot, parseListingPrice } from '$lib/listing';
import { sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, cookies }) => {
	if (!locals.user) redirect(303, '/login');
	// Fetch authed so the API computes is_own_listing for this viewer; a 404 for
	// someone else's listing avoids confirming it exists behind an edit URL.
	let res: Response;
	try {
		res = await apiFetch(`/api/v1/listings/${params.id}`, { headers: sessionHeader(cookies) });
	} catch {
		error(502, { message: 'The API is unreachable.', key: 'apiError.unreachable' });
	}
	if (res.status === 404) error(404, { message: 'Not found', key: 'apiError.not_found' });
	if (!res.ok)
		error(502, { message: 'The API returned an unexpected error.', key: 'apiError.unknown' });
	const listing = (await res.json()) as ListingDetail;
	if (!listing.is_own_listing) error(404, { message: 'Not found', key: 'apiError.not_found' });
	if (listing.status !== 'active') redirect(303, '/account/listings');
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
			return fail(400, { error: t(parsed.key), values: listingFormSnapshot(form) });
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
			return fail(502, { error: t('apiError.unreachable'), values: listingFormSnapshot(form) });
		}

		if (res.status === 403) {
			return fail(403, { error: t('formError.editOwnOnly'), values: listingFormSnapshot(form) });
		}
		if (res.status === 409) {
			return fail(409, { error: t('formError.editNotActive'), values: listingFormSnapshot(form) });
		}
		if (!res.ok) {
			return fail(res.status >= 500 ? 502 : res.status, {
				error: t('formError.editFailed'),
				values: listingFormSnapshot(form)
			});
		}

		redirect(303, '/account/listings');
	}
};
