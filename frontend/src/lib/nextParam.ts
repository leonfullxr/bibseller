/**
 * Clamp a ?next= redirect target to a same-site path. Rejects absolute URLs,
 * scheme-relative //host, and /\host (browsers treat \ as / in URLs).
 */
export function safeNext(next: string | null): string {
	return next !== null && next.startsWith('/') && !next.startsWith('//') && !next.startsWith('/\\')
		? next
		: '/';
}
