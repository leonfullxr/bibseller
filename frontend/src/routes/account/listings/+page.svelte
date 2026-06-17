<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { formatDate, formatPrice } from '$lib/format';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
</script>

<svelte:head>
	<title>My listings - Bibseller</title>
</svelte:head>

<div class="head">
	<h1>My listings</h1>
	<a href={resolve('/sell')} class="new">Sell a bib</a>
</div>

{#if form?.error}
	<p class="feedback error" role="alert">{form.error}</p>
{/if}

{#if data.listings.length === 0}
	<p class="empty">You have no listings yet. <a href={resolve('/sell')}>List a bib</a>.</p>
{:else}
	<ul class="listings">
		{#each data.listings as l (l.id)}
			<li>
				<div class="info">
					<a href={resolve('/races/[slug]', { slug: l.race_slug })} class="race">{l.race_name}</a>
					<p class="meta">
						{formatDate(l.event_date)} - {formatPrice(l.price_cents, l.currency) ??
							'Price on request'}
					</p>
				</div>
				<div class="right">
					<span class="status {l.status}">{l.status}</span>
					{#if l.status === 'active'}
						<a href={resolve('/account/listings/[id]/edit', { id: l.id })} class="edit">Edit</a>
						<form method="POST" action="?/cancel" use:enhance>
							<input type="hidden" name="id" value={l.id} />
							<button type="submit" class="cancel">Cancel</button>
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
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 0.75rem 1rem;
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
		align-items: center;
		gap: 0.75rem;
		white-space: nowrap;
	}

	.status {
		border-radius: 9999px;
		padding: 0.125rem 0.5rem;
		font-size: 0.6875rem;
		font-weight: 600;
		text-transform: capitalize;
		background: var(--slate-200);
		color: var(--slate-700);
	}

	.status.active {
		background: var(--emerald-100);
		color: var(--emerald-800);
	}

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
</style>
