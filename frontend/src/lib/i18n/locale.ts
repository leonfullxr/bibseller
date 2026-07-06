// Locale model + the framework-free routing/detection helpers (D17). English
// lives at the root, Spanish under `/es`; the URL is authoritative for what
// renders, while the detection chain below only decides the first-visit
// redirect target. Kept import-free so hooks and unit tests can use it directly.

export type Locale = 'en' | 'es';
export const defaultLocale: Locale = 'en';

// Presence of this cookie means "the visitor has a locale already" - it is what
// makes the first-visit redirect fire exactly once (see hooks.server.ts).
export const LOCALE_COOKIE = 'locale';
export const LOCALE_COOKIE_MAX_AGE = 60 * 60 * 24 * 365;

function isLocale(value: string | null | undefined): value is Locale {
	return value === 'en' || value === 'es';
}

/** The locale a path renders in: `/es` (or `/es/...`) is Spanish, everything else English. */
export function localeFromPath(pathname: string): Locale {
	return pathname === '/es' || pathname.startsWith('/es/') ? 'es' : 'en';
}

/** Drops the `/es` prefix, yielding the locale-free route path (`/es/races` -> `/races`). */
export function stripLocale(pathname: string): string {
	return localeFromPath(pathname) === 'es' ? pathname.slice(3) || '/' : pathname;
}

/** Re-applies a locale to a locale-free path: ('es','/races') -> '/es/races', ('en',…) unchanged. */
export function pathForLocale(locale: Locale, basePath: string): string {
	if (locale === 'en') return basePath;
	return basePath === '/' ? '/es' : `/es${basePath}`;
}

/**
 * Picks the best supported locale from an Accept-Language header, honouring
 * quality weights (`es-ES,es;q=0.9,en;q=0.8`). Unsupported tags are ignored;
 * an empty/missing header falls back to English.
 */
export function detectFromAcceptLanguage(header: string | null): Locale {
	if (!header) return defaultLocale;
	const ranked = header
		.split(',')
		.map((part) => {
			const [tag, ...params] = part.trim().split(';');
			const qParam = params.map((p) => p.trim()).find((p) => p.startsWith('q='));
			// A missing q defaults to 1; a malformed q (q=abc -> NaN) is coerced to 0
			// so it sorts last and the comparator never sees NaN (unstable order).
			const q = qParam ? Number.parseFloat(qParam.slice(2)) : 1;
			return { lang: tag.toLowerCase().split('-')[0], q: Number.isFinite(q) ? q : 0 };
		})
		.filter((t) => isLocale(t.lang))
		.sort((a, b) => b.q - a.q);
	return (ranked[0]?.lang as Locale) ?? defaultLocale;
}

/**
 * Whether to suggest Spanish to a visitor with no settled locale: their geo
 * country (Cloudflare's `cf-ipcountry`) is ES, or - when no geo signal is
 * available (e.g. local dev) - their browser prefers Spanish. A soft signal:
 * it raises the dismissible banner, never a redirect (D17). `country` of `''`
 * or `null` means "unknown" and falls through to the browser.
 */
export function suggestsSpanish(country: string | null, acceptLanguage: string | null): boolean {
	if (country) return country.toUpperCase() === 'ES';
	return detectFromAcceptLanguage(acceptLanguage) === 'es';
}

// ponytail: substring heuristic, not a maintained bot list. Its only job is to
// keep crawlers on the canonical URL they requested (no detection redirect), so
// hreflang stays clean. Swap for a real list if SEO ever misbehaves.
const BOT_UA = /bot|crawl|spider|slurp|bingpreview|facebookexternalhit|embedly|whatsapp|telegram/i;
export function isBot(userAgent: string | null): boolean {
	return !!userAgent && BOT_UA.test(userAgent);
}
