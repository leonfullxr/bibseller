import { fail } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { Actions, PageServerLoad } from './$types';

/** The reset link lands here as /reset?token=... (mirrors the verify flow). */
export const load: PageServerLoad = ({ url }) => {
	return { token: url.searchParams.get('token') ?? '' };
};

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const token = String(data.get('token') ?? '');
		const password = String(data.get('password') ?? '');
		const confirm = String(data.get('confirm') ?? '');

		// Server-side mirror of the HTML5 constraints; the browser attributes
		// are UX, not trust. Length floor matches the API (8).
		if (!token) {
			return fail(400, { error: 'This reset link is missing its token.' });
		}
		if (password.length < 8) {
			return fail(400, { error: 'Password must be at least 8 characters.' });
		}
		if (password !== confirm) {
			return fail(400, { error: 'The two passwords do not match.' });
		}

		let res: Response;
		try {
			res = await apiFetch('/api/v1/auth/password/reset', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ token, password })
			});
		} catch {
			return fail(502, { error: 'The API is unreachable.' });
		}

		if (res.status === 400) {
			return fail(400, { error: 'This reset link is invalid or has expired.' });
		}
		if (!res.ok) {
			return fail(502, { error: 'Could not reset your password.' });
		}

		return { done: true };
	}
};
