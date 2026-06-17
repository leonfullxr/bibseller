import { fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { clearSessionCookie, sessionHeader } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	if (!locals.user) redirect(303, '/login');
	return { user: locals.user };
};

export const actions: Actions = {
	profile: async ({ request, cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');

		const data = await request.formData();
		const displayName = String(data.get('display_name') ?? '').trim();

		// Server-side mirror of the HTML5 constraints - the browser attributes
		// are UX, not security; never trust the client.
		if (displayName.length < 2 || displayName.length > 50) {
			return fail(400, {
				value: displayName,
				error: 'Display name must be between 2 and 50 characters.'
			});
		}

		let res: Response;
		try {
			res = await apiFetch(`/api/v1/users/${locals.user.id}`, {
				method: 'PATCH',
				// Forward the session so the API can confirm we own this id (403 otherwise).
				headers: { 'Content-Type': 'application/json', ...sessionHeader(cookies) },
				body: JSON.stringify({ display_name: displayName })
			});
		} catch {
			return fail(502, { value: displayName, error: 'The API is unreachable.' });
		}

		if (!res.ok) {
			// The Go API's error envelope carries a human message; surface it.
			const body = (await res.json().catch(() => null)) as {
				error?: { message?: string };
			} | null;
			return fail(res.status >= 500 ? 502 : res.status, {
				value: displayName,
				error: body?.error?.message ?? 'The API rejected the update.'
			});
		}

		const user = (await res.json()) as { id: string; display_name: string };
		// Keep this request's locals in sync so the re-run layout load (and thus
		// the nav) shows the new name immediately, not on the next navigation.
		locals.user.display_name = user.display_name;
		return { success: true, value: user.display_name };
	},

	changePassword: async ({ request, cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');

		const data = await request.formData();
		const current = String(data.get('current_password') ?? '');
		const next = String(data.get('new_password') ?? '');
		const confirm = String(data.get('confirm_password') ?? '');

		// Server-side mirror of the HTML5 constraints; the 8-char floor matches
		// the API. The current password is checked there, not here.
		if (next.length < 8) {
			return fail(400, { pwError: 'New password must be at least 8 characters.' });
		}
		if (next !== confirm) {
			return fail(400, { pwError: 'The two new passwords do not match.' });
		}

		let res: Response;
		try {
			res = await apiFetch('/api/v1/auth/password', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json', ...sessionHeader(cookies) },
				body: JSON.stringify({ current_password: current, new_password: next })
			});
		} catch {
			return fail(502, { pwError: 'The API is unreachable.' });
		}

		if (res.status === 401) {
			return fail(400, { pwError: 'Your current password is incorrect.' });
		}
		if (!res.ok) {
			return fail(res.status >= 500 ? 502 : res.status, {
				pwError: 'Could not change your password.'
			});
		}

		return { pwSuccess: true };
	},

	logoutAll: async ({ cookies, locals }) => {
		if (!locals.user) redirect(303, '/login');

		// Idempotent on the API side (a dead token is still a 204). This revokes
		// every session including the current one, so drop the cookie and send
		// the user back to log in regardless of the API's answer.
		await apiFetch('/api/v1/auth/logout/all', {
			method: 'POST',
			headers: sessionHeader(cookies)
		}).catch(() => {});
		clearSessionCookie(cookies);
		redirect(303, '/login');
	}
};
