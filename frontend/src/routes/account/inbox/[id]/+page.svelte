<script lang="ts">
	import { onMount, tick, untrack } from 'svelte';
	import { resolve } from '$app/paths';
	import { formatDate, formatTime } from '$lib/format';
	import { getI18n } from '$lib/i18n';
	import { requiresAck } from '$lib/policy';
	import type { MessageKey } from '$lib/i18n';
	import type { ChatMessage } from '$lib/api/types';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, locale, link } = getI18n();

	// $derived so they track navigation between threads - SvelteKit reuses this
	// component across the [id] param rather than remounting.
	const meId = $derived(data.meId);
	const threadId = $derived(data.thread.id);

	// Seeded from the server load for SSR/first paint, then owned locally as
	// polling and sending append to it. untrack captures only the initial value.
	let messages = $state<ChatMessage[]>(untrack(() => data.messages));
	let draft = $state('');
	let sending = $state(false);
	let error = $state('');
	let list = $state<HTMLDivElement>();
	let files = $state<FileList | null>(null);
	let fileInput = $state<HTMLInputElement>();
	let polling = false; // in-flight guard so slow polls cannot overlap
	let notice = $state(''); // status line for block / report actions
	let mounted = $state(false); // gates local-timezone times (see formatTime)
	let stale = $state(false); // true after 3 consecutive poll failures
	let pollFails = 0;
	let preview = $state(''); // data-URL thumbnail of the attached image

	// Re-seed and jump to the latest when navigating to a different thread.
	$effect(() => {
		if (threadId) {
			messages = data.messages;
			tick().then(() => list?.scrollTo({ top: list.scrollHeight }));
		}
	});

	// The poll cursor is the newest id we hold (UUIDv7 -> time-ordered).
	function cursor(): string {
		return messages.length ? messages[messages.length - 1].id : '';
	}

	async function merge(incoming: ChatMessage[]) {
		const seen = new Set(messages.map((m) => m.id));
		const fresh = incoming.filter((m) => !seen.has(m.id));
		if (!fresh.length) return;
		// Only auto-scroll when the reader is pinned near the bottom - a poll must
		// not yank someone away from re-reading older messages.
		const pinned = !list || list.scrollHeight - list.scrollTop - list.clientHeight < 48;
		messages = [...messages, ...fresh];
		await tick();
		if (pinned) list?.scrollTo({ top: list.scrollHeight });
	}

	async function poll() {
		if (polling) return; // a previous poll is still in flight; skip this tick
		polling = true;
		try {
			const res = await fetch(`/api/v1/threads/${threadId}/messages?since=${cursor()}`, {
				credentials: 'same-origin'
			});
			if (res.ok) {
				await merge(((await res.json()) as { items: ChatMessage[] }).items);
				pollFails = 0;
				stale = false;
			} else {
				stale = ++pollFails >= 3;
			}
		} catch {
			/* transient; the next tick retries */
			stale = ++pollFails >= 3;
		} finally {
			polling = false;
		}
	}

	onMount(() => {
		mounted = true; // hydration is done; times may switch to the local timezone
		let id: ReturnType<typeof setInterval> | undefined;
		function schedule() {
			clearInterval(id);
			// D13: 3-5s active polling. A hidden/backgrounded tab backs off ~10x
			// (#96) - there's nothing new to show a tab nobody is looking at.
			id = setInterval(poll, document.visibilityState === 'hidden' ? 30_000 : 4_000);
		}
		schedule();
		document.addEventListener('visibilitychange', schedule);
		return () => {
			clearInterval(id);
			document.removeEventListener('visibilitychange', schedule);
		};
	});

	// The API's stable send-error codes -> translated messages (backend
	// internal/chat/http.go). Unlisted codes fall back to chat.sendFailed.
	const sendErrorKeys: Record<string, MessageKey> = {
		invalid_image: 'chat.invalidImage',
		unsupported_image: 'chat.unsupportedImage',
		blocked: 'chat.blockedSend',
		image_too_large: 'chat.imageTooLarge'
	};

	// Enter sends; Shift+Enter inserts a newline. Routes through the form's
	// onsubmit, so the same empty/in-flight guard applies. IME composition
	// (isComposing) must not send.
	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
			e.preventDefault();
			(e.currentTarget as HTMLTextAreaElement).form?.requestSubmit();
		}
	}

	function clearFile() {
		files = null;
		if (fileInput) fileInput.value = '';
	}

	// Thumbnail preview of the attached image. Skips files over the API's 5 MiB
	// cap - they can never send anyway, and reading them wastes memory.
	$effect(() => {
		const f = files?.[0];
		preview = '';
		if (!f || f.size > 5 << 20) return;
		let cancelled = false; // local to this effect run; not the poll-death `stale` above
		const reader = new FileReader();
		reader.onload = () => {
			if (!cancelled) preview = reader.result as string;
		};
		reader.readAsDataURL(f);
		return () => {
			cancelled = true;
			reader.abort();
		};
	});

	// UTC date-part of a timestamp, for the day separators. Deterministic on
	// server and client, so it needs no `mounted` gate.
	function dayOf(iso: string): string {
		return new Date(iso).toISOString().slice(0, 10);
	}

	// Consecutive same-sender, same-day messages render as one visual group:
	// tight gaps, tail corner on the group's last bubble.
	function sameGroup(a: ChatMessage, b: ChatMessage): boolean {
		return a.sender_id === b.sender_id && dayOf(a.created_at) === dayOf(b.created_at);
	}

	// Header avatar letter; spread handles astral-plane initials.
	const initial = $derived([...data.thread.other_party][0]?.toUpperCase() ?? '?');

	// Chip label next to the thumbnail (user data, no i18n).
	const fileName = $derived(files?.[0]?.name ?? '');

	// Bubble timestamps: UTC on SSR/first paint (hydration-safe), the viewer's
	// own timezone once mounted (formatTime's contract in $lib/format).
	function bubbleTime(iso: string): string {
		if (!mounted)
			return new Intl.DateTimeFormat(locale, { timeStyle: 'short', timeZone: 'UTC' }).format(
				new Date(iso)
			);
		return formatTime(iso, locale);
	}

	async function send(e: SubmitEvent) {
		e.preventDefault();
		const text = draft.trim();
		const file = files?.[0] ?? null;
		if ((!text && !file) || sending) return;
		sending = true;
		error = '';
		try {
			let res: Response;
			if (file) {
				// Image upload: multipart, with the text as an optional caption. The
				// browser sets the multipart Content-Type (boundary) itself.
				const fd = new FormData();
				fd.append('image', file);
				if (text) fd.append('body', text);
				res = await fetch(`/api/v1/threads/${threadId}/messages`, {
					method: 'POST',
					credentials: 'same-origin',
					body: fd
				});
			} else {
				res = await fetch(`/api/v1/threads/${threadId}/messages`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					credentials: 'same-origin',
					body: JSON.stringify({ body: text })
				});
			}
			if (res.ok) {
				await merge([(await res.json()) as ChatMessage]);
				// Sending always snaps to the bottom, even if the reader had
				// scrolled up (merge only scrolls while pinned).
				list?.scrollTo({ top: list.scrollHeight });
				draft = '';
				clearFile();
			} else if (res.status === 413) {
				error = t('chat.imageTooLarge');
			} else if (res.status === 429) {
				error = t('chat.tooFast');
			} else {
				// Map the API's stable error code to a translated message.
				const detail = (await res.json().catch(() => null)) as {
					error?: { code?: string };
				} | null;
				error = t(sendErrorKeys[detail?.error?.code ?? ''] ?? 'chat.sendFailed');
			}
		} catch {
			error = t('chat.networkError');
		} finally {
			sending = false;
		}
	}

	async function blockUser() {
		if (!window.confirm(t('chat.blockConfirm', { name: data.thread.other_party }))) return;
		try {
			const res = await fetch('/api/v1/blocks', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'same-origin',
				body: JSON.stringify({ blocked_id: data.thread.other_party_id })
			});
			notice = res.ok ? t('chat.blocked') : t('chat.blockFailed');
		} catch {
			notice = t('chat.networkRetry');
		}
	}

	async function unblockUser() {
		try {
			const res = await fetch(`/api/v1/blocks/${data.thread.other_party_id}`, {
				method: 'DELETE',
				credentials: 'same-origin'
			});
			notice = res.ok ? t('chat.unblocked') : t('chat.unblockFailed');
		} catch {
			notice = t('chat.networkRetry');
		}
	}

	async function reportMessage(id: string) {
		if (!window.confirm(t('chat.reportConfirm'))) return;
		try {
			const res = await fetch('/api/v1/reports', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'same-origin',
				body: JSON.stringify({ subject_type: 'message', subject_id: id, reason: 'other' })
			});
			notice = res.ok ? t('chat.messageReported') : t('chat.messageReportFailed');
		} catch {
			notice = t('chat.messageReportFailed');
		}
	}
</script>

<svelte:head>
	<title>{t('chat.title', { name: data.thread.other_party })}</title>
</svelte:head>

<nav><a href={link(resolve('/account/inbox'))}>{t('chat.back')}</a></nav>

<section class="chat">
	<header class="head">
		<span class="avatar" aria-hidden="true">{initial}</span>
		<div class="who">
			<h1>{data.thread.other_party}</h1>
			<p class="about">
				{t('chat.about')}
				<a href={link(resolve('/listings/[id]', { id: data.thread.listing_id }))}
					>{data.thread.race_name}</a
				>
			</p>
		</div>
		<details class="safety">
			<summary>{t('chat.safetySummary')}</summary>
			<div class="safety-actions">
				<button type="button" onclick={blockUser}>{t('chat.block')}</button>
				<button type="button" onclick={unblockUser}>{t('chat.unblock')}</button>
				{#if notice}<span class="notice" role="status">{notice}</span>{/if}
			</div>
		</details>
	</header>

	{#if requiresAck(data.thread.transfer_policy)}
		<p class="alert policy-note" role="note">{t('chat.policyReminder')}</p>
	{/if}

	<!-- The tab stop lets keyboard users scroll the log; svelte-ignore because
	     role="log" is not interactive, but a scrollable region needs focus. -->
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<div class="messages" bind:this={list} role="log" tabindex="0" aria-label={t('chat.logAria')}>
		{#each messages as m, i (m.id)}
			{#if i === 0 || dayOf(messages[i - 1].created_at) !== dayOf(m.created_at)}
				<div class="day"><span>{formatDate(dayOf(m.created_at), locale)}</span></div>
			{/if}
			<div
				class="msg"
				class:mine={m.sender_id === meId}
				class:grouped={i > 0 && sameGroup(messages[i - 1], m)}
				class:tail={i === messages.length - 1 || !sameGroup(m, messages[i + 1])}
			>
				{#if m.has_image}
					<img
						class="image"
						src={`/api/v1/threads/${threadId}/messages/${m.id}/image`}
						alt={m.body ?? t('chat.sharedImage')}
						loading="lazy"
					/>
				{/if}
				{#if m.body}<p class="body">{m.body}</p>{/if}
				<div class="msg-foot">
					{#if m.sender_id !== meId}
						<button type="button" class="report-msg" onclick={() => reportMessage(m.id)}
							>{t('chat.reportMsg')}</button
						>
					{/if}
					<span class="time">{bubbleTime(m.created_at)}</span>
				</div>
			</div>
		{/each}
	</div>

	{#if stale}
		<p class="stale" role="status">{t('chat.connectionLost')}</p>
	{/if}

	<form class="composer" onsubmit={send}>
		{#if error}
			<p class="alert" role="alert">{error}</p>
		{/if}
		{#if preview}
			<div class="chip">
				<img class="chip-img" src={preview} alt={t('chat.previewAlt')} />
				<span class="chip-name">{fileName}</span>
				<button
					type="button"
					class="chip-clear"
					onclick={clearFile}
					aria-label={t('chat.clearImage')}
				>
					<svg
						viewBox="0 0 24 24"
						width="14"
						height="14"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						aria-hidden="true"
					>
						<path d="M18 6 6 18M6 6l12 12" />
					</svg>
				</button>
			</div>
		{/if}
		<div class="dock">
			<label class="attach">
				<input
					type="file"
					accept="image/jpeg,image/png"
					aria-label={t('chat.attachAria')}
					bind:files
					bind:this={fileInput}
				/>
				<svg
					viewBox="0 0 24 24"
					width="20"
					height="20"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					<path
						d="M21.44 11.05 12.25 20.24a6 6 0 0 1-8.49-8.49l8.57-8.57a4 4 0 1 1 5.66 5.66l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48"
					/>
				</svg>
			</label>
			<textarea
				class="field"
				bind:value={draft}
				rows="2"
				maxlength="4000"
				aria-label={t('chat.messageAria')}
				placeholder={t('chat.messagePlaceholder')}
				onkeydown={onKeydown}></textarea>
			<button type="submit" class="btn btn-primary" disabled={sending}
				>{sending ? t('chat.sending') : t('chat.send')}</button
			>
		</div>
	</form>
</section>

<style>
	nav {
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	nav a {
		color: var(--slate-500);
	}

	nav a:hover {
		color: var(--slate-900);
	}

	/* One seamless chat card: header bar / policy strip / message canvas /
	   composer dock, separated by hairlines instead of stacked fragments. */
	.chat {
		margin-top: 1rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.75rem;
		background: white;
		overflow: hidden;
	}

	.head {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid var(--slate-200);
	}

	.avatar {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		flex-shrink: 0;
		border-radius: 9999px;
		background: var(--emerald-100);
		color: var(--emerald-800);
		font-weight: 700;
	}

	.who {
		flex: 1;
		min-width: 0;
	}

	h1 {
		font-size: 1rem;
		line-height: 1.5rem;
		font-weight: 600;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.about {
		font-size: 0.8125rem;
		color: var(--slate-500);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.about a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	/* Safety cluster: a quiet chip that opens as a small popover, so blocking
	   controls stay one click away without crowding the header. */
	.safety {
		position: relative;
		font-size: 0.8125rem;
	}

	.safety summary {
		list-style: none;
		cursor: pointer;
		padding: 0.375rem 0.625rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-300);
		color: var(--slate-600);
		font-weight: 500;
		white-space: nowrap;
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.safety summary::-webkit-details-marker {
		display: none;
	}

	.safety summary:hover {
		background: var(--slate-100);
	}

	.safety-actions {
		position: absolute;
		right: 0;
		top: calc(100% + 0.375rem);
		z-index: 10;
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: 0.5rem;
		max-width: min(20rem, 80vw);
		width: max-content;
		padding: 0.5rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		box-shadow: 0 4px 12px rgb(15 23 42 / 0.08);
	}

	.safety button {
		border-radius: 0.375rem;
		background: none;
		color: var(--slate-600);
		border: 1px solid var(--slate-300);
		padding: 0.25rem 0.625rem;
		font-size: 0.8125rem;
		font-weight: 500;
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.safety button:hover {
		background: var(--slate-100);
	}

	.notice {
		color: var(--emerald-700);
	}

	/* Policy reminder as a full-width strip inside the card. */
	.policy-note {
		border-radius: 0;
		border-left: none;
		border-right: none;
		border-top: none;
	}

	.messages {
		display: flex;
		flex-direction: column;
		max-height: 60dvh;
		min-height: 12rem;
		overflow-y: auto;
		scrollbar-gutter: stable;
		background: var(--slate-50);
		padding: 1rem;
	}

	.messages:focus-visible {
		outline: 2px solid var(--emerald-600);
		outline-offset: -2px;
	}

	.day {
		align-self: center;
		margin: 0.875rem 0 0.375rem;
	}

	.day:first-child {
		margin-top: 0.125rem;
	}

	.day span {
		display: inline-block;
		border-radius: 9999px;
		padding: 0.125rem 0.625rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		background: var(--slate-100);
		border: 1px solid var(--slate-200);
		color: var(--slate-600);
	}

	/* Bubbles: theirs = white on the left, mine = solid emerald on the right
	   (white text on emerald-700 is 5.5:1). Same-sender runs group tightly;
	   the last bubble of a group gets the tail corner. */
	.msg {
		max-width: 75%;
		align-self: flex-start;
		margin-top: 0.75rem;
		border-radius: 1rem;
		border: 1px solid var(--slate-200);
		background: white;
		color: var(--slate-900);
		padding: 0.5rem 0.75rem;
		box-shadow: 0 1px 2px rgb(15 23 42 / 0.04);
	}

	.msg.grouped {
		margin-top: 0.125rem;
	}

	.msg:first-child,
	.day + .msg {
		margin-top: 0;
	}

	.msg.tail {
		border-bottom-left-radius: 0.25rem;
	}

	.msg.mine {
		align-self: flex-end;
		border-color: var(--emerald-700);
		background: var(--emerald-700);
		color: white;
		border-bottom-left-radius: 1rem;
	}

	.msg.mine.tail {
		border-bottom-right-radius: 0.25rem;
	}

	.body {
		white-space: pre-wrap;
		overflow-wrap: anywhere;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: inherit;
	}

	.image {
		display: block;
		max-width: 100%;
		max-height: 20rem;
		border-radius: 0.625rem;
		margin-bottom: 0.375rem;
	}

	.msg-foot {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.75rem;
		margin-top: 0.125rem;
	}

	.time {
		font-size: 0.6875rem;
		line-height: 1rem;
		color: var(--slate-500);
	}

	.mine .time {
		color: var(--emerald-100);
	}

	.msg-foot button {
		margin-right: auto;
		background: none;
		color: var(--slate-500);
		padding: 0;
		font-size: 0.6875rem;
		font-weight: 400;
	}

	.msg-foot button:hover {
		background: none;
		color: var(--amber-700);
	}

	/* Declutter: on hover-capable devices the report link appears only while the
	   bubble is hovered or holds focus; touch keeps it always visible. */
	@media (hover: hover) {
		.msg .report-msg {
			opacity: 0;
			pointer-events: none;
		}

		.msg:hover .report-msg,
		.msg:focus-within .report-msg {
			opacity: 1;
			pointer-events: auto;
		}
	}

	.stale {
		text-align: center;
		font-size: 0.8125rem;
		padding: 0.375rem 0.75rem;
		background: var(--amber-50);
		border-top: 1px solid var(--amber-100);
		color: var(--amber-800);
	}

	/* Composer dock: attach / auto-growing field / send on one baseline. */
	.composer {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 0.75rem;
		border-top: 1px solid var(--slate-200);
		background: white;
	}

	.dock {
		display: flex;
		align-items: flex-end;
		gap: 0.5rem;
	}

	/* The real file input stays focusable inside the label (so Tab + Enter
	   still open the picker) but is visually hidden behind the icon button. */
	.attach {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		flex-shrink: 0;
		position: relative;
		border-radius: 0.5rem;
		color: var(--slate-500);
		cursor: pointer;
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.attach input {
		position: absolute;
		width: 1px;
		height: 1px;
		opacity: 0;
		overflow: hidden;
	}

	.attach:hover {
		background: var(--slate-100);
		color: var(--slate-700);
	}

	.attach:focus-within {
		outline: 2px solid var(--emerald-600);
		outline-offset: 2px;
	}

	textarea {
		flex: 1;
		min-width: 0;
		font: inherit;
		font-size: 0.875rem;
		line-height: 1.25rem;
		resize: none;
		min-height: 2.5rem;
		max-height: 8rem;
		field-sizing: content; /* auto-grow; browsers without it keep rows=2 */
	}

	.dock .btn {
		flex-shrink: 0;
	}

	.chip {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.chip-img {
		height: 3rem;
		width: 3rem;
		object-fit: cover;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
	}

	.chip-name {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-size: 0.8125rem;
		color: var(--slate-600);
	}

	.composer .chip-clear {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.75rem;
		height: 1.75rem;
		flex-shrink: 0;
		border-radius: 9999px;
		background: none;
		color: var(--slate-600);
		border: 1px solid var(--slate-300);
		font-size: 1rem;
		line-height: 1;
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.composer .chip-clear:hover {
		background: var(--slate-100);
		color: var(--slate-900);
	}
</style>
