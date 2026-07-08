<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ListingSummary } from '$lib/api/types';
	import { formatPrice } from '$lib/format';
	import { getI18n } from '$lib/i18n';

	let { listing }: { listing: ListingSummary } = $props();
	const { t, locale, link } = getI18n();

	const price = $derived(formatPrice(listing.price_cents, listing.currency, locale));
	const original = $derived(formatPrice(listing.original_price_cents, listing.currency, locale));
	const belowFace = $derived(
		listing.price_cents != null &&
			listing.original_price_cents != null &&
			listing.price_cents < listing.original_price_cents
	);
</script>

<a href={link(resolve('/listings/[id]', { id: listing.id }))} class="card">
	<div class="price-row">
		<span class="price">{price ?? t('listingCard.priceOnRequest')}</span>
		{#if belowFace && original}
			<span class="original">{original}</span>
		{/if}
	</div>
	{#if belowFace}
		<span class="deal">{t('listingCard.belowFace')}</span>
	{/if}
	{#if listing.description}
		<p class="desc">{listing.description}</p>
	{/if}
	<p class="seller">{t('listingCard.listedBy', { name: listing.seller_name })}</p>
</a>

<style>
	/* Same bib-tag family as RaceCard: hard ink border, poster shadow that
	   flattens a step on hover. */
	.card {
		display: block;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 1rem;
		box-shadow: var(--shadow-hard-sm);
		transition:
			translate 0.1s,
			box-shadow 0.1s;
	}

	.card:hover {
		translate: 1px 1px;
		box-shadow: 2px 2px 0 var(--ink);
	}

	.price-row {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.5rem;
	}

	/* The price is the bib number: condensed, heavy, ink. */
	.price {
		font-family: var(--font-display);
		font-size: 1.5rem;
		line-height: 1.75rem;
		font-weight: 600;
		color: var(--ink);
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
