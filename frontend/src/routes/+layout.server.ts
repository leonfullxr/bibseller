import type { LayoutServerLoad } from './$types';

/** Surfaces the signed-in user, active locale, and the Spanish-suggestion flag (all set by hooks.server.ts) to every page + the nav. */
export const load: LayoutServerLoad = ({ locals }) => {
	return { user: locals.user, locale: locals.locale, suggestLocale: locals.suggestLocale };
};
