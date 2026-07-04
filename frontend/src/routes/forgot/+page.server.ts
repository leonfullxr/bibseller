import { fail } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { createTranslator } from '$lib/i18n';
import { clientIPHeader } from '$lib/server/clientip';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request, locals }) => {
		const t = createTranslator(locals.locale);
		const data = await request.formData();
		const email = String(data.get('email') ?? '').trim();

		if (!email.includes('@')) {
			return fail(400, { email, error: t('formError.invalidEmail') });
		}

		let res: Response;
		try {
			// On success the API answers 204 whether or not the email exists (no
			// account-enumeration oracle), so we never branch on that. We do
			// surface a 429 (per-IP limiter) or a server error, neither of which
			// reveals whether the account exists.
			res = await apiFetch('/api/v1/auth/password/reset/request', {
				method: 'POST',
				// clientIPHeader: the per-IP limiter must see the user, not this
				// server (#133).
				headers: { 'Content-Type': 'application/json', ...clientIPHeader(request) },
				body: JSON.stringify({ email })
			});
		} catch {
			return fail(502, { email, error: t('apiError.unreachable') });
		}

		if (res.status === 429) {
			return fail(429, { email, error: t('formError.tooManyRequests') });
		}
		if (!res.ok) {
			return fail(502, { email, error: t('formError.resetEmailFailed') });
		}

		return { sent: true };
	}
};
