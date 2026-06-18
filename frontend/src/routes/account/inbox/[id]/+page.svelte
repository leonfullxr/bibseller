<script lang="ts">
	import { onMount, tick, untrack } from 'svelte';
	import { resolve } from '$app/paths';
	import { formatDateTime } from '$lib/format';
	import type { ChatMessage } from '$lib/api/types';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

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
		messages = [...messages, ...fresh];
		await tick();
		list?.scrollTo({ top: list.scrollHeight });
	}

	async function poll() {
		if (polling) return; // a previous poll is still in flight; skip this tick
		polling = true;
		try {
			const res = await fetch(`/api/v1/threads/${threadId}/messages?since=${cursor()}`, {
				credentials: 'same-origin'
			});
			if (res.ok) await merge(((await res.json()) as { items: ChatMessage[] }).items);
		} catch {
			/* transient; the next tick retries */
		} finally {
			polling = false;
		}
	}

	onMount(() => {
		const id = setInterval(poll, 4000); // D13: 3-5s polling
		return () => clearInterval(id);
	});

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
				draft = '';
				files = null;
				if (fileInput) fileInput.value = '';
			} else if (res.status === 413) {
				error = 'That image is too large (5 MB max).';
			} else if (res.status === 429) {
				error = 'You are sending messages too fast - wait a moment.';
			} else {
				// Surface the API's specific reason (bad image, message/caption too long, ...).
				const detail = (await res.json().catch(() => null)) as {
					error?: { message?: string };
				} | null;
				error = detail?.error?.message ?? 'Could not send your message. Try again.';
			}
		} catch {
			error = 'Network error - check your connection.';
		} finally {
			sending = false;
		}
	}

	async function blockUser() {
		if (
			!window.confirm(
				`Block ${data.thread.other_party}? Neither of you will be able to message the other.`
			)
		)
			return;
		try {
			const res = await fetch('/api/v1/blocks', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'same-origin',
				body: JSON.stringify({ blocked_id: data.thread.other_party_id })
			});
			notice = res.ok ? 'User blocked.' : 'Could not block the user.';
		} catch {
			notice = 'Network error - try again.';
		}
	}

	async function unblockUser() {
		try {
			const res = await fetch(`/api/v1/blocks/${data.thread.other_party_id}`, {
				method: 'DELETE',
				credentials: 'same-origin'
			});
			notice = res.ok ? 'User unblocked.' : 'Could not unblock the user.';
		} catch {
			notice = 'Network error - try again.';
		}
	}

	async function reportMessage(id: string) {
		if (!window.confirm('Report this message to the moderators?')) return;
		try {
			const res = await fetch('/api/v1/reports', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'same-origin',
				body: JSON.stringify({ subject_type: 'message', subject_id: id, reason: 'other' })
			});
			notice = res.ok ? 'Message reported.' : 'Could not report the message.';
		} catch {
			notice = 'Could not report the message.';
		}
	}
</script>

<svelte:head>
	<title>Chat with {data.thread.other_party} - Bibseller</title>
</svelte:head>

<nav><a href={resolve('/account/inbox')}>Back to inbox</a></nav>

<header class="head">
	<h1>{data.thread.other_party}</h1>
	<p class="about">
		about <a href={resolve('/listings/[id]', { id: data.thread.listing_id })}
			>{data.thread.race_name}</a
		>
	</p>
	<div class="safety">
		<button type="button" onclick={blockUser}>Block</button>
		<button type="button" onclick={unblockUser}>Unblock</button>
		{#if notice}<span class="notice" role="status">{notice}</span>{/if}
	</div>
</header>

<div class="messages" bind:this={list}>
	{#each messages as m (m.id)}
		<div class="msg" class:mine={m.sender_id === meId}>
			{#if m.has_image}
				<img
					class="image"
					src={`/api/v1/threads/${threadId}/messages/${m.id}/image`}
					alt={m.body ?? 'Shared image'}
					loading="lazy"
				/>
			{/if}
			{#if m.body}<p class="body">{m.body}</p>{/if}
			<div class="msg-foot">
				<span class="time">{formatDateTime(m.created_at)}</span>
				<button type="button" class="report-msg" onclick={() => reportMessage(m.id)}>report</button>
			</div>
		</div>
	{/each}
</div>

<form class="composer" onsubmit={send}>
	<textarea
		bind:value={draft}
		rows="3"
		maxlength="4000"
		aria-label="Your message"
		placeholder="Write a message, or attach an image..."
	></textarea>
	{#if error}
		<p class="feedback error" role="alert">{error}</p>
	{/if}
	<div class="actions">
		<input
			class="file"
			type="file"
			accept="image/jpeg,image/png"
			aria-label="Attach an image (JPEG or PNG)"
			bind:files
			bind:this={fileInput}
		/>
		<button type="submit" disabled={sending}>{sending ? 'Sending...' : 'Send'}</button>
	</div>
</form>

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

	.head {
		margin-top: 1rem;
	}

	h1 {
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 700;
	}

	.about {
		margin-top: 0.125rem;
		font-size: 0.875rem;
		color: var(--slate-600);
	}

	.about a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.messages {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 60vh;
		overflow-y: auto;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: var(--slate-50);
		padding: 1rem;
	}

	.msg {
		max-width: 80%;
		align-self: flex-start;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 0.5rem 0.75rem;
	}

	.msg.mine {
		align-self: flex-end;
		border-color: var(--emerald-200);
		background: var(--emerald-50);
	}

	.body {
		white-space: pre-wrap;
		overflow-wrap: anywhere;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-900);
	}

	.image {
		display: block;
		max-width: 100%;
		max-height: 20rem;
		border-radius: 0.375rem;
		margin-bottom: 0.375rem;
	}

	.time {
		margin-top: 0.25rem;
		display: block;
		font-size: 0.6875rem;
		line-height: 1rem;
		color: var(--slate-400);
	}

	.composer {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.actions {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
	}

	.file {
		font-size: 0.8125rem;
		color: var(--slate-600);
		min-width: 0;
	}

	textarea {
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font: inherit;
		font-size: 0.875rem;
		resize: vertical;
	}

	.feedback {
		border-radius: 0.375rem;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.error {
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		color: var(--amber-900);
	}

	button {
		align-self: flex-start;
		border-radius: 0.375rem;
		background: var(--emerald-600);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	button:hover:not(:disabled) {
		background: var(--emerald-700);
	}

	button:disabled {
		opacity: 0.6;
		cursor: default;
	}

	.safety {
		margin-top: 0.5rem;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.safety button {
		align-self: auto;
		background: none;
		color: var(--slate-600);
		border: 1px solid var(--slate-300);
		padding: 0.25rem 0.625rem;
		font-size: 0.8125rem;
		font-weight: 500;
	}

	.safety button:hover {
		background: var(--slate-100);
	}

	.notice {
		color: var(--emerald-700);
	}

	.msg-foot {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
	}

	.msg-foot button {
		align-self: auto;
		background: none;
		color: var(--slate-400);
		padding: 0;
		font-size: 0.6875rem;
		font-weight: 400;
	}

	.msg-foot button:hover {
		background: none;
		color: var(--amber-700);
	}
</style>
