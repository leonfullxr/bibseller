import { fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { SessionResponse } from '$lib/api/types';
import { setSessionCookie } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	if (locals.user) redirect(303, '/');
};

export const actions: Actions = {
	default: async ({ request, cookies, locals }) => {
		const data = await request.formData();
		const email = String(data.get('email') ?? '').trim();
		const displayName = String(data.get('display_name') ?? '').trim();
		const password = String(data.get('password') ?? '');

		// Echo email + display name back on failure so the form stays filled.
		// The password is deliberately NEVER echoed: it would end up in the
		// server-rendered HTML of the response.
		const values = { email, display_name: displayName };

		// Server-side mirror of the HTML5 constraints (the attributes are UX,
		// not enforcement). The Go API re-validates again - it can't trust us.
		if (!email.includes('@')) {
			return fail(400, { ...values, error: 'Enter a valid email address.' });
		}
		if (displayName.length < 2 || displayName.length > 50) {
			return fail(400, { ...values, error: 'Display name must be between 2 and 50 characters.' });
		}
		if (password.length < 8) {
			return fail(400, { ...values, error: 'Password must be at least 8 characters.' });
		}

		let res: Response;
		try {
			res = await apiFetch('/api/v1/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				// Forward the locale the signup happened in so the account + its
				// verification email default to it (the API re-validates).
				body: JSON.stringify({ email, password, display_name: displayName, locale: locals.locale })
			});
		} catch {
			return fail(502, { ...values, error: 'The API is unreachable.' });
		}

		if (!res.ok) {
			// Surface the API's error-envelope message ("email_taken" -> 409,
			// validation -> 400) into the form.
			const body = (await res.json().catch(() => null)) as {
				error?: { message?: string };
			} | null;
			return fail(res.status >= 500 ? 502 : res.status, {
				...values,
				error: body?.error?.message ?? 'Could not create the account.'
			});
		}

		// The raw session token arrives here once, server-to-server; turning
		// it into the __Host-session cookie is this layer's job.
		const session = (await res.json()) as SessionResponse;
		setSessionCookie(cookies, session.token, session.expires_at);

		// Post/Redirect/Get: the browser lands on a fresh GET, so refreshing
		// never re-submits the registration.
		redirect(303, '/');
	}
};
