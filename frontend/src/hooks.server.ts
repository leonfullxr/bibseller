/**
 * Populates `locals.user` on every request by resolving the __Host-session
 * cookie against the Go API's GET /auth/me. This is the server-side half the
 * auth design anticipates (docs/ARCHITECTURE.md -> Auth & sessions): the cookie
 * is an opaque token, so "who is this?" can only be answered by the API.
 */
import type { Cookies, Handle } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { SessionUser } from '$lib/api/types';
import { SESSION_COOKIE, sessionHeader } from '$lib/server/session';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.user = await resolveUser(event.cookies);
	return resolve(event);
};

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
