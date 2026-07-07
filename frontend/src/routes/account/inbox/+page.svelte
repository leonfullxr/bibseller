<script lang="ts">
	import { resolve } from '$app/paths';
	import { formatWhen } from '$lib/format';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t: tr, plural, locale, link } = getI18n();
</script>

<svelte:head>
	<title>{tr('inbox.title')}</title>
</svelte:head>

<h1>{tr('inbox.heading')}</h1>

{#if data.threads.length === 0}
	<p class="empty">
		{tr('inbox.emptyPre')}
		<a href={link(resolve('/races'))}>{tr('inbox.emptyRacesLink')}</a>
		{tr('inbox.emptyPost')}
	</p>
{:else}
	<ul class="threads">
		{#each data.threads as t (t.id)}
			<li>
				<a
					href={link(resolve('/account/inbox/[id]', { id: t.id }))}
					class="row"
					class:unread={t.unread_count > 0}
				>
					<div class="who">
						<span class="name">{t.other_party}</span>
						<span class="tag">{t.role === 'buyer' ? tr('role.seller') : tr('role.buyer')}</span>
					</div>
					<span class="race">{t.race_name}</span>
					<div class="right">
						{#if t.last_message_at}
							<span class="when">{formatWhen(t.last_message_at, locale)}</span>
						{/if}
						{#if t.unread_count > 0}
							<span class="badge" aria-label={plural('nav.inboxUnread', t.unread_count)}
								>{t.unread_count}</span
							>
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
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 0.25rem 1rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 0.75rem 1rem;
	}

	.row:hover {
		border-color: var(--slate-300);
	}

	.row.unread {
		border-color: var(--brand-300);
		background: var(--brand-50);
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
		min-width: 8rem;
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
		margin-left: auto;
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
		background: var(--brand-700);
		padding: 0.125rem 0.375rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 700;
		color: white;
	}
</style>
