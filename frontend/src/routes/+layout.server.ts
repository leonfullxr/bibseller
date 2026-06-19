import type { LayoutServerLoad } from './$types';

/** Surfaces the signed-in user and active locale (both set by hooks.server.ts) to every page + the nav. */
export const load: LayoutServerLoad = ({ locals }) => {
	return { user: locals.user, locale: locals.locale };
};
