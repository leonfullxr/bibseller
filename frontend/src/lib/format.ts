/** Locale-aware display helpers. Locale stays 'en' until i18n lands (M8). */

export function formatPrice(
	cents: number | null | undefined,
	currency = 'EUR',
	locale = 'en'
): string | null {
	if (cents == null) return null;
	const whole = cents % 100 === 0;
	return new Intl.NumberFormat(locale, {
		style: 'currency',
		currency,
		minimumFractionDigits: whole ? 0 : 2,
		maximumFractionDigits: whole ? 0 : 2
	}).format(cents / 100);
}

/** Formats a YYYY-MM-DD date without timezone drift. */
export function formatDate(isoDate: string, locale = 'en'): string {
	return new Intl.DateTimeFormat(locale, { dateStyle: 'long', timeZone: 'UTC' }).format(
		new Date(`${isoDate}T00:00:00Z`)
	);
}

/**
 * Formats an RFC 3339 timestamp as a date + time. Pinned to UTC and an explicit
 * locale so server and client render identically (no hydration mismatch);
 * consistent with formatDate showing UTC. Localized times wait for i18n (M8).
 */
export function formatDateTime(iso: string, locale = 'en'): string {
	return new Intl.DateTimeFormat(locale, {
		dateStyle: 'medium',
		timeStyle: 'short',
		timeZone: 'UTC'
	}).format(new Date(iso));
}

/** Today as YYYY-MM-DD (UTC) - the default lower bound for race browsing. */
export function todayISO(): string {
	return new Date().toISOString().slice(0, 10);
}
