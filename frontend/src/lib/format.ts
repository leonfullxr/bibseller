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

/**
 * Time-aware short form for inbox rows: a timestamp from today (UTC) shows as
 * a time, anything older as a date. Same UTC pin as its siblings, so server
 * and client render identically.
 */
export function formatWhen(iso: string, locale = 'en'): string {
	const d = new Date(iso);
	const today = d.toISOString().slice(0, 10) === todayISO();
	return new Intl.DateTimeFormat(
		locale,
		today ? { timeStyle: 'short', timeZone: 'UTC' } : { dateStyle: 'medium', timeZone: 'UTC' }
	).format(d);
}

/** Today as YYYY-MM-DD (UTC) - the default lower bound for race browsing. */
export function todayISO(): string {
	return new Date().toISOString().slice(0, 10);
}
