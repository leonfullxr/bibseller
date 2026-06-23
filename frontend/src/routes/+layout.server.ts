import type { LayoutServerLoad } from './$types';

/** Surfaces the signed-in user, active locale, and the Spanish-suggestion flag (all set by hooks.server.ts) to every page + the nav. */
export const load: LayoutServerLoad = ({ locals, depends }) => {
	// Tagged so the verify page can refresh just the user (clearing the "verify
	// your email" banner) after confirmation, without re-running other loads.
	depends('app:user');
	return { user: locals.user, locale: locals.locale, suggestLocale: locals.suggestLocale };
};
