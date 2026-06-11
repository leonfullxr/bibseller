/**
 * Joins an API base URL and a path without doubling or dropping slashes.
 * Pure function so it stays unit-testable outside the SvelteKit runtime.
 */
export function apiUrl(base: string, path: string): string {
	return `${base.replace(/\/+$/, '')}/${path.replace(/^\/+/, '')}`;
}
