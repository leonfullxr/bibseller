import { error, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { ChatMessage, ChatThreadSummary } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, cookies }) => {
	if (!locals.user) redirect(303, '/login');
	if (!locals.user.email_verified) redirect(303, '/settings');

	// The inbox carries this thread's header context (other party, race). Reuse
	// it rather than add a single-thread endpoint - a user's inbox is small.
	let threadsRes: Response;
	try {
		threadsRes = await apiFetch('/api/v1/threads', { headers: sessionHeader(cookies) });
	} catch {
		error(502, 'The API is unreachable.');
	}
	if (!threadsRes.ok) error(502, 'Could not load the conversation.');
	const inbox = (await threadsRes.json()) as { items: ChatThreadSummary[] };
	const thread = inbox.items.find((t) => t.id === params.id);
	if (!thread) error(404, 'Conversation not found');

	let msgsRes: Response;
	try {
		msgsRes = await apiFetch(`/api/v1/threads/${params.id}/messages`, {
			headers: sessionHeader(cookies)
		});
	} catch {
		error(502, 'The API is unreachable.');
	}
	if (!msgsRes.ok) error(502, 'Could not load messages.');
	const msgs = (await msgsRes.json()) as { items: ChatMessage[] };

	return { thread, messages: msgs.items, meId: locals.user.id };
};
