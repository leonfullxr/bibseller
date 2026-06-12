import { fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { sessionHeader, setSessionCookie } from '$lib/server/session';
import type { Actions } from './$types';

interface SessionResponse {
	token: string;
	expires_at: string;
	user: { id: string; email: string; display_name: string };
}

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const email = String(data.get('email') ?? '').trim();
		const password = String(data.get('password') ?? '');

		if (!email.includes('@') || password.length === 0) {
			return fail(400, { email, error: 'Enter your email and password.' });
		}

		let res: Response;
		try {
			res = await apiFetch('/api/v1/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					// Forward any existing session so the API rotates it out:
					// a token minted before this authentication must not
					// survive it (session fixation defense).
					...sessionHeader(cookies)
				},
				body: JSON.stringify({ email, password })
			});
		} catch {
			return fail(502, { email, error: 'The API is unreachable.' });
		}

		if (res.status === 401) {
			// Deliberately vague, mirroring the API: which half was wrong is
			// exactly what an account-enumeration attacker wants to know.
			return fail(401, { email, error: 'Invalid email or password.' });
		}
		if (!res.ok) {
			return fail(502, { email, error: 'Could not log in.' });
		}

		const session = (await res.json()) as SessionResponse;
		setSessionCookie(cookies, session.token, session.expires_at);

		redirect(303, '/');
	}
};
