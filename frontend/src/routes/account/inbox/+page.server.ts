import { redirect } from '@sveltejs/kit';
import { apiGet } from '$lib/api/server';
import type { ChatThreadSummary } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies, fetch }) => {
	if (!locals.user) redirect(303, '/login');
	// Chat requires a verified email; the API would 403 the inbox otherwise.
	if (!locals.user.email_verified) redirect(303, '/settings');

	const data = await apiGet<{ items: ChatThreadSummary[] }>('/api/v1/threads', fetch, {
		headers: sessionHeader(cookies)
	});
	return { threads: data.items };
};
