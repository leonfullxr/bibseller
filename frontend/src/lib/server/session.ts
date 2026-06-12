/**
 * Browser-side half of the session contract (docs/ARCHITECTURE.md → Auth &
 * sessions). The Go API mints the token and stores only its SHA-256; this
 * module is the ONLY place the raw token is turned into a cookie. It lives
 * under $lib/server so SvelteKit refuses to bundle it into client code.
 */
import type { Cookies } from '@sveltejs/kit';

export const SESSION_COOKIE = '__Host-session';

export function setSessionCookie(cookies: Cookies, token: string, expiresAt: string): void {
	cookies.set(SESSION_COOKIE, token, {
		// The __Host- prefix is a browser-enforced contract: the cookie is
		// rejected unless Secure and Path=/ are set and Domain is absent.
		// That pins it to exactly this host over HTTPS — no subdomain can
		// shadow it, no plaintext network attacker can overwrite it.
		path: '/',
		secure: true, // dev note: Chrome/Firefox accept Secure cookies on
		// localhost (it's a trustworthy origin); Safari does not — use
		// Chrome/Firefox for local auth work.
		httpOnly: true, // invisible to document.cookie: XSS can't exfiltrate it
		sameSite: 'lax', // sent on top-level navigations, withheld from
		// cross-site subrequests/POSTs — first CSRF layer, complemented by
		// the API's Sec-Fetch-Site guard.

		// The cookie dies when the session would (30 days). The server slides
		// its idle window on activity but the cookie keeps the original
		// stamp; login refreshes both. Good enough until a session-refresh
		// hook exists.
		expires: new Date(expiresAt)
	});
}

/** Forward the browser's session to the Go API on a server-side fetch. */
export function sessionHeader(cookies: Cookies): Record<string, string> {
	const token = cookies.get(SESSION_COOKIE);
	return token ? { cookie: `${SESSION_COOKIE}=${token}` } : {};
}
