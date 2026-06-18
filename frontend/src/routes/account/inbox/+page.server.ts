import { error, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { ChatThreadSummary } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.user) redirect(303, '/login');
	// Chat requires a verified email; the API would 403 the inbox otherwise.
	if (!locals.user.email_verified) redirect(303, '/settings');

	let res: Response;
	try {
		res = await apiFetch('/api/v1/threads', { headers: sessionHeader(cookies) });
	} catch {
		error(502, 'The API is unreachable.');
	}
	if (!res.ok) error(502, 'Could not load your inbox.');

	const data = (await res.json()) as { items: ChatThreadSummary[] };
	return { threads: data.items };
};
