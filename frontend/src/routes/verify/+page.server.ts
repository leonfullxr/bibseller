import { apiFetch } from '$lib/api/server';
import type { PageServerLoad } from './$types';

/**
 * Landing page for the emailed verification link (GET /verify?token=...).
 * Consumes the token against the API server-to-server - the raw token never
 * touches client JS.
 */
export const load: PageServerLoad = async ({ url }) => {
	const token = url.searchParams.get('token');
	if (!token) return { status: 'missing' as const };
	try {
		const res = await apiFetch('/api/v1/auth/verify', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ token })
		});
		if (res.ok) return { status: 'ok' as const };
		if (res.status === 400) return { status: 'invalid' as const };
		return { status: 'error' as const };
	} catch {
		return { status: 'error' as const };
	}
};
