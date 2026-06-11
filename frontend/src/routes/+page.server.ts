import { apiFetch } from '$lib/api/server';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		const res = await apiFetch('/api/healthz', { signal: AbortSignal.timeout(1500) });
		return { apiStatus: res.ok ? ('ok' as const) : ('down' as const) };
	} catch {
		return { apiStatus: 'down' as const };
	}
};
