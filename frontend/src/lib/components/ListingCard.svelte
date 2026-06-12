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

<a
	href={resolve('/listings/[id]', { id: listing.id })}
	class="block rounded-lg border border-slate-200 bg-white p-4 transition hover:border-slate-300 hover:shadow-sm"
>
	<div class="flex items-baseline justify-between gap-2">
		<span class="text-xl font-bold text-slate-900">{price ?? 'Price on request'}</span>
		{#if belowFace && original}
			<span class="text-sm text-slate-400 line-through">{original}</span>
		{/if}
	</div>
	{#if belowFace}
		<span
			class="mt-1 inline-block rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-semibold text-emerald-800"
		>
			below face value
		</span>
	{/if}
	{#if listing.description}
		<p class="mt-2 line-clamp-2 text-sm text-slate-600">{listing.description}</p>
	{/if}
	<p class="mt-3 text-xs text-slate-500">Listed by {listing.seller_name}</p>
</a>
