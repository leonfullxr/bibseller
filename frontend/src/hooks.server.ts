/**
 * Populates `locals.user` on every request by resolving the __Host-session
 * cookie against the Go API's GET /auth/me. This is the server-side half the
 * auth design anticipates (docs/ARCHITECTURE.md -> Auth & sessions): the cookie
 * is an opaque token, so "who is this?" can only be answered by the API.
 *
 * It also resolves the active locale (D17). The URL is authoritative for what
 * renders - `/es...` is Spanish, everything else English. A *settled* choice
 * (signed-in `users.locale`, else the `locale` cookie) redirects the URL to
 * match it; soft signals (geo IP, `Accept-Language`) never redirect - they only
 * raise the dismissible "switch to Spanish" banner on the English pages.
 */
import type { Cookies, Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';
import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import { apiFetch } from '$lib/api/server';
import type { SessionUser } from '$lib/api/types';
import { SESSION_COOKIE, sessionHeader } from '$lib/server/session';
import {
	LOCALE_COOKIE,
	type Locale,
	isBot,
	localeFromPath,
	pathForLocale,
	stripLocale,
	suggestsSpanish
} from '$lib/i18n/locale';

export const handle: Handle = async ({ event, resolve }) => {
	const user = await resolveUser(event.cookies);
	event.locals.user = user;

	const urlLocale = localeFromPath(event.url.pathname);
	event.locals.locale = urlLocale;
	event.locals.suggestLocale = null;

	const crawler = isBot(event.request.headers.get('user-agent'));

	// A settled choice is authoritative: redirect so the URL matches it. Crawlers
	// are never redirected, keeping both language URLs crawlable. No settled
	// choice => no redirect (the visitor lands on the URL they asked for, English
	// at the root).
	const settled = settledLocale(user, event.cookies.get(LOCALE_COOKIE));
	if (settled && settled !== urlLocale && !crawler) {
		redirect(302, pathForLocale(settled, stripLocale(event.url.pathname)) + event.url.search);
	}

	// No settled choice, on an English page, not a crawler: suggest Spanish via a
	// dismissible banner when the visitor looks Spanish by location. cf-ipcountry
	// is Cloudflare's geo header in prod; in dev a DEV_IP_COUNTRY env var (or a
	// `curl -H "cf-ipcountry: ES"`) stands in. The banner persists across pages
	// until the visitor accepts (cookie=es + /es) or dismisses (cookie=en).
	if (!settled && !crawler && urlLocale === 'en') {
		const country =
			event.request.headers.get('cf-ipcountry') ?? (dev ? (env.DEV_IP_COUNTRY ?? null) : null);
		if (suggestsSpanish(country, event.request.headers.get('accept-language'))) {
			event.locals.suggestLocale = 'es';
		}
	}

	return resolve(event, { transformPageChunk: ({ html }) => html.replace('%lang%', urlLocale) });
};

// The settled locale: a signed-in account preference, else the cookie set by the
// switcher/banner. Soft signals (IP, Accept-Language) are deliberately excluded.
function settledLocale(user: SessionUser | null, cookie: string | undefined): Locale | null {
	if (user?.locale === 'en' || user?.locale === 'es') return user.locale;
	if (cookie === 'en' || cookie === 'es') return cookie;
	return null;
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
