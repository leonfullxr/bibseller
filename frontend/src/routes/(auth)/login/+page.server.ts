import { fail, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import type { SessionResponse } from '$lib/api/types';
import { safeNext } from '$lib/nextParam';
import { clientIPHeader } from '$lib/server/clientip';
import { sessionHeader, setSessionCookie } from '$lib/server/session';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals, url }) => {
	if (locals.user) redirect(303, safeNext(url.searchParams.get('next')));
};

export const actions: Actions = {
	default: async ({ request, cookies, locals, url }) => {
		const t = createTranslator(locals.locale);
		const data = await request.formData();
		const email = String(data.get('email') ?? '').trim();
		const password = String(data.get('password') ?? '');

		if (!email.includes('@') || password.length === 0) {
			return fail(400, { email, error: t('formError.loginRequired') });
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
					...sessionHeader(cookies),
					// Forward the client address so the API's per-IP limiter and
					// session audit see the user, not this server (#133).
					...clientIPHeader(request)
				},
				body: JSON.stringify({ email, password })
			});
		} catch {
			return fail(502, { email, error: t('apiError.unreachable') });
		}

		if (res.status === 401) {
			// Deliberately vague, mirroring the API: which half was wrong is
			// exactly what an account-enumeration attacker wants to know.
			return fail(401, { email, error: t('formError.invalidCredentials') });
		}
		if (!res.ok) {
			return fail(502, { email, error: t('formError.loginFailed') });
		}

		const session = (await res.json()) as SessionResponse;
		setSessionCookie(cookies, session.token, session.expires_at);

		// ?next= survives the POST (the form submits to its own URL). Validated
		// server-side: only same-site paths, never an absolute/protocol-relative URL.
		redirect(303, safeNext(url.searchParams.get('next')));
	}
};
