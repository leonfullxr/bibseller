// Universal reroute: strip the `/es` prefix so Spanish URLs match the same route
// tree as English (D17). The original URL is preserved on `event.url`, so
// hooks.server.ts still sees `/es/...` to resolve the active locale, and
// `page.url` keeps the prefix for hreflang/canonical output.
import type { Reroute } from '@sveltejs/kit';
import { stripLocale } from '$lib/i18n/locale';

export const reroute: Reroute = ({ url }) => {
	if (url.pathname === '/es' || url.pathname.startsWith('/es/')) {
		return stripLocale(url.pathname);
	}
};
