<script lang="ts">
	import type { SubmitFunction } from '@sveltejs/kit';
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { formatDate, formatPrice } from '$lib/format';
	import { getI18n, listingStatusLabel } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, locale, link } = getI18n();

	// Active rows first, newest first within each group.
	const listings = $derived(
		[...data.listings].sort(
			(a, b) =>
				Number(b.status === 'active') - Number(a.status === 'active') ||
				b.created_at.localeCompare(a.created_at)
		)
	);

	// The publish flow redirects here with ?created=1 (sell/[slug] action).
	const created = $derived(page.url.searchParams.get('created') === '1');

	const confirmCancel: SubmitFunction = ({ cancel }) => {
		if (!window.confirm(t('myListings.cancelConfirm'))) cancel();
	};
</script>

<svelte:head>
	<title>{t('myListings.title')}</title>
</svelte:head>

<div class="head">
	<h1>{t('myListings.heading')}</h1>
	<a href={link(resolve('/sell'))} class="new">{t('sell.heading')}</a>
</div>

{#if created}
	<p class="feedback success" role="status">{t('myListings.created')}</p>
{/if}

{#if form?.error}
	<p class="feedback error" role="alert">{form.error}</p>
{/if}

{#if form?.cancelled}
	<p class="feedback success" role="status">{t('myListings.cancelled')}</p>
{/if}

{#if data.listings.length === 0}
	<p class="empty">
		{t('myListings.emptyPre')}
		<a href={link(resolve('/sell'))}>{t('myListings.listABib')}</a>.
	</p>
{:else}
	<ul class="listings">
		{#each listings as l (l.id)}
			<li class:done={l.status !== 'active'}>
				<div class="info">
					<a href={link(resolve('/races/[slug]', { slug: l.race_slug }))} class="race"
						>{l.race_name}</a
					>
					<p class="meta">
						{formatDate(l.event_date, locale)} - {formatPrice(l.price_cents, l.currency, locale) ??
							t('listingCard.priceOnRequest')}
					</p>
				</div>
				<div class="right">
					<span class="status {l.status}">{listingStatusLabel(t, l.status)}</span>
					<a href={link(resolve('/listings/[id]', { id: l.id }))} class="view"
						>{t('myListings.view')}</a
					>
					{#if l.status === 'active'}
						<a href={link(resolve('/account/listings/[id]/edit', { id: l.id }))} class="edit"
							>{t('myListings.edit')}</a
						>
						<form method="POST" action="?/cancel" use:enhance={confirmCancel}>
							<input type="hidden" name="id" value={l.id} />
							<button type="submit" class="cancel">{t('myListings.cancel')}</button>
						</form>
					{/if}
				</div>
			</li>
		{/each}
	</ul>
{/if}

<style>
	.head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
	}

	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.new {
		border-radius: 0.375rem;
		background: var(--emerald-600);
		padding: 0.375rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	.new:hover {
		background: var(--emerald-700);
	}

	.empty {
		margin-top: 2rem;
		color: var(--slate-600);
	}

	.empty a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.listings {
		margin-top: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		list-style: none;
		padding: 0;
	}

	.listings li {
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

	.listings li.done {
		background: var(--slate-50);
		opacity: 0.7;
	}

	.info {
		min-width: 0;
	}

	.race {
		font-weight: 600;
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.meta {
		margin-top: 0.125rem;
		font-size: 0.75rem;
		color: var(--slate-500);
	}

	.right {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.25rem 0.75rem;
		margin-left: auto;
		white-space: nowrap;
	}

	.status {
		border-radius: 9999px;
		padding: 0.125rem 0.5rem;
		font-size: 0.6875rem;
		font-weight: 600;
		background: var(--slate-200);
		color: var(--slate-700);
	}

	.status.active {
		background: var(--emerald-100);
		color: var(--emerald-800);
	}

	.status.sold {
		background: var(--sky-100);
		color: var(--sky-800);
	}

	.status.reserved {
		background: var(--amber-100);
		color: var(--amber-800);
	}

	.view,
	.edit {
		font-size: 0.875rem;
		color: var(--slate-700);
		text-decoration: underline;
	}

	.cancel {
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.375rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--slate-700);
	}

	.cancel:hover {
		background: var(--slate-100);
	}

	.feedback {
		margin-top: 1rem;
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

	.success {
		border: 1px solid var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}
</style>
