import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import { sessionHeader } from '$lib/server/session';
import {
	LOCALE_COOKIE,
	LOCALE_COOKIE_MAX_AGE,
	type Locale,
	pathForLocale,
	stripLocale
} from '$lib/i18n/locale';
import type { Actions } from './$types';

// The nav language switcher posts here. Sets the locale cookie (which both
// persists the choice and, by its presence, stops the first-visit redirect from
// second-guessing it) and, when signed in, mirrors it to users.locale. A plain
// form POST means this is a full navigation, so the target page reloads in the
// new locale - no client-side locale reactivity needed.
export const actions: Actions = {
	default: async ({ request, cookies, locals }) => {
		const data = await request.formData();
		const to: Locale = data.get('to') === 'es' ? 'es' : 'en';

		// `next` is the path (with query) to return to. Only same-origin absolute
		// paths are allowed - reject protocol-relative (`//host`) and backslash
		// tricks so the switcher can't be turned into an open redirect. Strip any
		// locale prefix so a `/es/...` value is not double-prefixed by pathForLocale.
		const nextRaw = String(data.get('next') ?? '/');
		const next = /^\/(?![/\\])/.test(nextRaw) ? stripLocale(nextRaw) : '/';

		cookies.set(LOCALE_COOKIE, to, { path: '/', maxAge: LOCALE_COOKIE_MAX_AGE, sameSite: 'lax' });

		if (locals.user) {
			// Best-effort persistence; the cookie already covers this browser. The
			// API requires the whole profile, so resend the unchanged name/country.
			await apiFetch(`/api/v1/users/${locals.user.id}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json', ...sessionHeader(cookies) },
				body: JSON.stringify({
					display_name: locals.user.display_name,
					locale: to,
					country: locals.user.country
				})
			}).catch(() => {});
			locals.user.locale = to;
		}

		redirect(303, pathForLocale(to, next));
	}
};
