import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { clearSessionCookie, sessionHeader } from '$lib/server/session';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ cookies }) => {
		// Tell the API to delete the session row (idempotent - a dead token is
		// still a 204), then drop the cookie regardless of the API's answer.
		await apiFetch('/api/v1/auth/logout', {
			method: 'POST',
			headers: sessionHeader(cookies)
		}).catch(() => {});
		clearSessionCookie(cookies);
		redirect(303, '/');
	}
};
