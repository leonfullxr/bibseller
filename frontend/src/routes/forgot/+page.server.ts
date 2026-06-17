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

		try {
			// The API always answers 204 (no account-enumeration oracle), so we
			// never branch on whether the email exists - the UI says the same
			// thing either way.
			await apiFetch('/api/v1/auth/password/reset/request', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email })
			});
		} catch {
			return fail(502, { email, error: 'The API is unreachable.' });
		}

		return { sent: true };
	}
};
