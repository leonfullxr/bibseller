<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ListingSummary } from '$lib/api/types';
	import { formatPrice } from '$lib/format';

	let { listing }: { listing: ListingSummary } = $props();

	const price = $derived(formatPrice(listing.price_cents, listing.currency));
	const original = $derived(formatPrice(listing.original_price_cents, listing.currency));
	const belowFace = $derived(
		listing.price_cents != null &&
			listing.original_price_cents != null &&
			listing.price_cents < listing.original_price_cents
	);
</script>

<a href={resolve('/listings/[id]', { id: listing.id })} class="card">
	<div class="price-row">
		<span class="price">{price ?? 'Price on request'}</span>
		{#if belowFace && original}
			<span class="original">{original}</span>
		{/if}
	</div>
	{#if belowFace}
		<span class="deal">below face value</span>
	{/if}
	{#if listing.description}
		<p class="desc">{listing.description}</p>
	{/if}
	<p class="seller">Listed by {listing.seller_name}</p>
</a>

<style>
	.card {
		display: block;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 1rem;
		transition:
			border-color 150ms,
			box-shadow 150ms;
	}

	.card:hover {
		border-color: var(--slate-300);
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
	}

	.price-row {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.price {
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 700;
		color: var(--slate-900);
	}

	.original {
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-400);
		text-decoration: line-through;
	}

	.deal {
		margin-top: 0.25rem;
		display: inline-block;
		border-radius: 9999px;
		background: var(--emerald-100);
		padding: 0.125rem 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 600;
		color: var(--emerald-800);
	}

	.desc {
		margin-top: 0.5rem;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		overflow: hidden;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.seller {
		margin-top: 0.75rem;
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}
</style>
