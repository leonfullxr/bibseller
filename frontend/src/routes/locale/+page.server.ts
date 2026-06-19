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

		// `next` is the path (with optional query) to return to. Split off the
		// query, then on the pathname only: require a same-origin absolute path
		// (reject protocol-relative `//host` and backslash tricks) and strip any
		// locale prefix so a `/es/...` value is not double-prefixed below. The
		// query is re-appended after pathForLocale so the result is canonical
		// (`/es/races?q=1`, never `/es/?q=1` or `/es/es?q=1`).
		const nextRaw = String(data.get('next') ?? '/');
		const queryAt = nextRaw.indexOf('?');
		const rawPath = queryAt === -1 ? nextRaw : nextRaw.slice(0, queryAt);
		let path = '/';
		let search = '';
		if (/^\/(?![/\\])/.test(rawPath)) {
			path = stripLocale(rawPath);
			search = queryAt === -1 ? '' : nextRaw.slice(queryAt);
		}

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

		redirect(303, pathForLocale(to, path) + search);
	}
};
