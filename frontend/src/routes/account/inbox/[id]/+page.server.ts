import { error, redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/server';
import type { ChatMessage, ChatThreadSummary } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, cookies }) => {
	if (!locals.user) redirect(303, '/login');
	if (!locals.user.email_verified) redirect(303, '/settings');

	// Dedicated header endpoint (#97) - avoids fetching the whole inbox just to
	// render one thread's header.
	let threadRes: Response;
	try {
		threadRes = await apiFetch(`/api/v1/threads/${params.id}`, { headers: sessionHeader(cookies) });
	} catch {
		error(502, { message: 'The API is unreachable.', key: 'apiError.unreachable' });
	}
	// 403 (exists, not yours) folds into the same "not found" as a missing
	// thread - this visitor gets no more signal than "not found" either way.
	if (threadRes.status === 404 || threadRes.status === 403)
		error(404, { message: 'Conversation not found', key: 'apiError.not_found' });
	if (!threadRes.ok)
		error(502, { message: 'Could not load the conversation.', key: 'apiError.loadFailed' });
	const thread = (await threadRes.json()) as ChatThreadSummary;

	let msgsRes: Response;
	try {
		msgsRes = await apiFetch(`/api/v1/threads/${params.id}/messages`, {
			headers: sessionHeader(cookies)
		});
	} catch {
		error(502, { message: 'The API is unreachable.', key: 'apiError.unreachable' });
	}
	if (!msgsRes.ok) error(502, { message: 'Could not load messages.', key: 'apiError.loadFailed' });
	const msgs = (await msgsRes.json()) as { items: ChatMessage[]; next_cursor: string | null };

	// ponytail: ListMessages pages oldest-first (id ASC), so this walks forward
	// from the start of the thread, capped at 5 extra pages (~600 messages
	// total). On threads longer than that, the newest messages - not older
	// ones - are what stays unloaded, until a proper tail/reverse-cursor
	// fetch exists (#125).
	const items = msgs.items;
	let cursor = msgs.next_cursor;
	for (let i = 0; i < 5 && cursor; i++) {
		let res: Response;
		try {
			res = await apiFetch(`/api/v1/threads/${params.id}/messages?since=${cursor}`, {
				headers: sessionHeader(cookies)
			});
		} catch {
			break; // transient network error mid-backfill - keep what we already have
		}
		if (!res.ok) break;
		const page = (await res.json()) as { items: ChatMessage[]; next_cursor: string | null };
		items.push(...page.items);
		cursor = page.next_cursor;
	}

	return { thread, messages: items, meId: locals.user.id };
};
