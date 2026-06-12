/**
 * Server-side API client: SvelteKit `load` functions call the Go API
 * directly (not through the Vite/Caddy proxy) using this helper.
 * Browser-side code uses relative `/api/...` URLs instead.
 */
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { apiUrl } from './url';

const base = () => env.API_URL || 'http://localhost:8080';

export function apiFetch(path: string, init?: RequestInit): Promise<Response> {
	return fetch(apiUrl(base(), path), init);
}

/** GET a JSON payload, translating API failures into SvelteKit errors. */
export async function apiGet<T>(path: string, fetchFn: typeof fetch = fetch): Promise<T> {
	let res: Response;
	try {
		res = await fetchFn(apiUrl(base(), path));
	} catch {
		error(502, 'The API is unreachable.');
	}
	if (res.status === 404) error(404, 'Not found');
	if (res.status === 400) error(400, 'Invalid request');
	if (!res.ok) error(502, 'The API returned an unexpected error.');
	return res.json() as Promise<T>;
}
