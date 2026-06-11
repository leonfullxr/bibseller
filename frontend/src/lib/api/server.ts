/**
 * Server-side API client: SvelteKit `load` functions call the Go API
 * directly (not through the Vite/Caddy proxy) using this helper.
 * Browser-side code uses relative `/api/...` URLs instead.
 */
import { env } from '$env/dynamic/private';
import { apiUrl } from './url';

const base = () => env.API_URL || 'http://localhost:8080';

export function apiFetch(path: string, init?: RequestInit): Promise<Response> {
	return fetch(apiUrl(base(), path), init);
}
