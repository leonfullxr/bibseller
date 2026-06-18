<script lang="ts">
	import { resolve } from '$app/paths';
	import { formatDate } from '$lib/format';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Inbox - Bibseller</title>
</svelte:head>

<h1>Inbox</h1>

{#if data.threads.length === 0}
	<p class="empty">
		No conversations yet. Browse <a href={resolve('/races')}>races</a> and contact a seller to start one.
	</p>
{:else}
	<ul class="threads">
		{#each data.threads as t (t.id)}
			<li>
				<a
					href={resolve('/account/inbox/[id]', { id: t.id })}
					class="row"
					class:unread={t.unread_count > 0}
				>
					<div class="who">
						<span class="name">{t.other_party}</span>
						<span class="tag">{t.role === 'buyer' ? 'seller' : 'buyer'}</span>
					</div>
					<span class="race">{t.race_name}</span>
					<div class="right">
						{#if t.last_message_at}
							<span class="when">{formatDate(t.last_message_at.slice(0, 10))}</span>
						{/if}
						{#if t.unread_count > 0}
							<span class="badge">{t.unread_count}</span>
						{/if}
					</div>
				</a>
			</li>
		{/each}
	</ul>
{/if}

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.empty {
		margin-top: 1.5rem;
		color: var(--slate-600);
	}

	.empty a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.threads {
		margin-top: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		list-style: none;
		padding: 0;
	}

	.row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 0.75rem 1rem;
	}

	.row:hover {
		border-color: var(--slate-300);
	}

	.row.unread {
		border-color: var(--emerald-300);
		background: var(--emerald-50);
	}

	.who {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		min-width: 0;
	}

	.name {
		font-weight: 600;
		color: var(--slate-900);
	}

	.tag {
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
		text-transform: capitalize;
	}

	.race {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-size: 0.875rem;
		color: var(--slate-600);
	}

	.right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		white-space: nowrap;
	}

	.when {
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}

	.badge {
		display: inline-flex;
		min-width: 1.25rem;
		justify-content: center;
		border-radius: 9999px;
		background: var(--emerald-600);
		padding: 0.125rem 0.375rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 700;
		color: white;
	}
</style>
