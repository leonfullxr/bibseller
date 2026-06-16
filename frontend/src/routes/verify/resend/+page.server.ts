import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { sessionHeader } from '$lib/server/session';
import type { Actions } from './$types';

/** Target of the "resend verification email" banner button. */
export const actions: Actions = {
	default: async ({ cookies }) => {
		// Best-effort and idempotent server-side (204 even if already verified).
		await apiFetch('/api/v1/auth/verify/resend', {
			method: 'POST',
			headers: sessionHeader(cookies)
		}).catch(() => {});
		redirect(303, '/?verification=sent');
	}
};
