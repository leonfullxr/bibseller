import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { clientIPHeader } from '$lib/server/clientip';
import { sessionHeader } from '$lib/server/session';
import type { Actions } from './$types';

/** Target of the "resend verification email" banner button. */
export const actions: Actions = {
	default: async ({ cookies, request }) => {
		// Best-effort and idempotent server-side (204 even if already verified).
		// clientIPHeader: the API's per-IP limiter must see the user, not this
		// server (#133).
		await apiFetch('/api/v1/auth/verify/resend', {
			method: 'POST',
			headers: { ...sessionHeader(cookies), ...clientIPHeader(request) }
		}).catch(() => {});
		redirect(303, '/?verification=sent');
	}
};
