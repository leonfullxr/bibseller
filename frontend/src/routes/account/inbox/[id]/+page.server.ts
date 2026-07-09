import { error, redirect } from '@sveltejs/kit';
import { apiFetch, apiGet } from '$lib/api/server';
import type { ChatMessage, ChatThreadSummary } from '$lib/api/types';
import { sessionHeader } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, cookies, fetch }) => {
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

	// The default page is now the thread's tail - the newest messages, ascending
	// (#154). prev_cursor, when set, is the "load earlier" cursor the page uses
	// to walk backward on demand; the old forward-backfill loop (which opened
	// long threads on their oldest messages) is gone.
	const msgs = await apiGet<{ items: ChatMessage[]; prev_cursor: string | null }>(
		`/api/v1/threads/${params.id}/messages`,
		fetch,
		{ headers: sessionHeader(cookies) }
	);

	return {
		thread,
		messages: msgs.items,
		earlierCursor: msgs.prev_cursor,
		meId: locals.user.id
	};
};
