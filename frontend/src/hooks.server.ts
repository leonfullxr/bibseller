/**
 * Populates `locals.user` on every request by resolving the __Host-session
 * cookie against the Go API's GET /auth/me. This is the server-side half the
 * auth design anticipates (docs/ARCHITECTURE.md -> Auth & sessions): the cookie
 * is an opaque token, so "who is this?" can only be answered by the API.
 *
 * It also resolves the active locale (D17). The URL is authoritative for what
 * renders - `/es...` is Spanish, everything else English - while the detection
 * chain (signed-in `users.locale` > `Accept-Language`) only decides the
 * first-visit redirect target.
 */
import type { Cookies, Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { SessionUser } from '$lib/api/types';
import { SESSION_COOKIE, sessionHeader } from '$lib/server/session';
import {
	LOCALE_COOKIE,
	LOCALE_COOKIE_MAX_AGE,
	type Locale,
	detectFromAcceptLanguage,
	isBot,
	localeFromPath,
	pathForLocale,
	stripLocale
} from '$lib/i18n/locale';

export const handle: Handle = async ({ event, resolve }) => {
	const user = await resolveUser(event.cookies);
	event.locals.user = user;

	const urlLocale = localeFromPath(event.url.pathname);
	event.locals.locale = urlLocale;

	// First-visit redirect (D17): the `locale` cookie marks "this visitor has a
	// locale", so once it is set the URL wins and this never fires again - which
	// also makes it loop-proof (the redirect target always matches its own
	// urlLocale). Crawlers are never redirected, keeping canonical URLs crawlable.
	if (!event.cookies.get(LOCALE_COOKIE) && !isBot(event.request.headers.get('user-agent'))) {
		const preferred = preferredLocale(user, event.request.headers.get('accept-language'));
		event.cookies.set(LOCALE_COOKIE, preferred, {
			path: '/',
			maxAge: LOCALE_COOKIE_MAX_AGE,
			sameSite: 'lax'
		});
		if (preferred !== urlLocale) {
			redirect(302, pathForLocale(preferred, stripLocale(event.url.pathname)) + event.url.search);
		}
	}

	return resolve(event, { transformPageChunk: ({ html }) => html.replace('%lang%', urlLocale) });
};

// Signed-in preference wins over the browser's Accept-Language (D17).
function preferredLocale(user: SessionUser | null, acceptLanguage: string | null): Locale {
	if (user?.locale === 'en' || user?.locale === 'es') return user.locale;
	return detectFromAcceptLanguage(acceptLanguage);
}

async function resolveUser(cookies: Cookies): Promise<SessionUser | null> {
	// Anonymous browsing (the whole public catalog) skips the API hop entirely.
	if (!cookies.get(SESSION_COOKIE)) return null;
	try {
		const res = await apiFetch('/api/v1/auth/me', { headers: sessionHeader(cookies) });
		// 401 here is normal (expired/rotated token), not an error - fall through
		// to signed-out rather than throwing a page error.
		return res.ok ? ((await res.json()) as SessionUser) : null;
	} catch {
		return null; // API unreachable: degrade to signed-out, never crash the page
	}
}
