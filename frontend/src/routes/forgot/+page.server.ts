import { fail } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const email = String(data.get('email') ?? '').trim();

		if (!email.includes('@')) {
			return fail(400, { email, error: 'Enter a valid email address.' });
		}

		let res: Response;
		try {
			// On success the API answers 204 whether or not the email exists (no
			// account-enumeration oracle), so we never branch on that. We do
			// surface a 429 (per-IP limiter) or a server error, neither of which
			// reveals whether the account exists.
			res = await apiFetch('/api/v1/auth/password/reset/request', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email })
			});
		} catch {
			return fail(502, { email, error: 'The API is unreachable.' });
		}

		if (res.status === 429) {
			return fail(429, { email, error: 'Too many requests. Please wait a minute and try again.' });
		}
		if (!res.ok) {
			return fail(502, { email, error: 'Could not send the reset email. Please try again.' });
		}

		return { sent: true };
	}
};
