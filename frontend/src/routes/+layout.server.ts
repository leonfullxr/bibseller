import { apiFetch } from '$lib/api/server';
import { sessionHeader } from '$lib/server/session';
import type { LayoutServerLoad } from './$types';

/** Surfaces the signed-in user, active locale, and the Spanish-suggestion flag (all set by hooks.server.ts) to every page + the nav. */
export const load: LayoutServerLoad = async ({ locals, cookies, depends }) => {
	// Tagged so the verify page can refresh just the user (clearing the "verify
	// your email" banner) after confirmation, without re-running other loads.
	depends('app:user');

	// Header unread badge. Verified users only (the API 403s the rest), and a
	// failure just means no badge - the nav must never break on a chat hiccup.
	// ponytail: no polling - refreshes per navigation, deliberate (PR3 design).
	let unreadCount = 0;
	if (locals.user?.email_verified) {
		try {
			const res = await apiFetch('/api/v1/me/unread-count', { headers: sessionHeader(cookies) });
			if (res.ok) {
				unreadCount = ((await res.json()) as { unread_count: number }).unread_count;
			}
		} catch {
			// unreachable API -> badge shows nothing
		}
	}

	return {
		user: locals.user,
		locale: locals.locale,
		suggestLocale: locals.suggestLocale,
		unreadCount
	};
};
