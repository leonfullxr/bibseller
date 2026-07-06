import { fail, redirect } from '@sveltejs/kit';
import { apiFetch, apiGet } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import type { ListingDetail } from '$lib/api/types';
import { requiresAck } from '$lib/policy';
import { sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch, setHeaders, locals, cookies }) => {
	// Signed-in viewers get the per-viewer is_own_listing and an uncacheable
	// response; anonymous viewers get the public, cacheable listing.
	if (locals.user) {
		const listing = await apiGet<ListingDetail>(`/api/v1/listings/${params.id}`, fetch, {
			headers: sessionHeader(cookies)
		});
		setHeaders({ 'cache-control': 'private, no-store' });
		return { listing };
	}
	const listing = await apiGet<ListingDetail>(`/api/v1/listings/${params.id}`, fetch);
	setHeaders({ 'cache-control': 'public, max-age=60' });
	return { listing };
};

export const actions: Actions = {
	// contact starts (or reuses) the buyer's thread with the seller, sending the
	// first message. In the restricted policy modes it first records the buyer's
	// acknowledgment of the venue-only terms - the chat API rejects the message
	// without it.
	contact: async ({ request, params, cookies, locals, fetch }) => {
		if (!locals.user) redirect(303, '/login');
		const t = createTranslator(locals.locale);
		if (!locals.user.email_verified) {
			return fail(403, { error: t('formError.verifyToContact'), body: '' });
		}

		const form = await request.formData();
		const body = String(form.get('body') ?? '').trim();
		if (body === '') {
			return fail(400, { error: t('formError.emptyMessage'), body: '' });
		}

		// Authoritative policy + slug from the API, never a tamperable hidden field.
		const listing = await apiGet<ListingDetail>(`/api/v1/listings/${params.id}`, fetch);
		if (requiresAck(listing.race.transfer_policy)) {
			if (form.get('ack') == null) {
				return fail(400, { error: t('formError.ackRequired'), body });
			}
			let ackRes: Response;
			try {
				ackRes = await apiFetch(`/api/v1/races/${listing.race.slug}/ack`, {
					method: 'POST',
					headers: sessionHeader(cookies)
				});
			} catch {
				return fail(502, { error: t('apiError.unreachable'), body });
			}
			if (!ackRes.ok) {
				return fail(502, { error: t('formError.ackFailed'), body });
			}
		}

		let res: Response;
		try {
			res = await apiFetch(`/api/v1/listings/${params.id}/threads`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json', ...sessionHeader(cookies) },
				body: JSON.stringify({ body })
			});
		} catch {
			return fail(502, { error: t('apiError.unreachable'), body });
		}
		if (res.status === 403) {
			return fail(403, { error: t('formError.cannotContact'), body });
		}
		if (res.status === 409) {
			return fail(409, { error: t('formError.listingUnavailable'), body });
		}
		if (!res.ok) {
			return fail(res.status >= 500 ? 502 : res.status, {
				error: t('formError.contactFailed'),
				body
			});
		}

		const { thread_id } = (await res.json()) as { thread_id: string };
		redirect(303, `/account/inbox/${thread_id}`);
	}
};
