<script lang="ts">
	import { resolve } from '$app/paths';
	import ListingCTA from '$lib/components/ListingCTA.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate, formatPrice } from '$lib/format';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const listing = $derived(data.listing);
	const race = $derived(data.listing.race);

	const price = $derived(formatPrice(listing.price_cents, listing.currency));
	const original = $derived(formatPrice(listing.original_price_cents, listing.currency));
	const belowFace = $derived(
		listing.price_cents != null &&
			listing.original_price_cents != null &&
			listing.price_cents < listing.original_price_cents
	);
	const available = $derived(listing.status === 'active');
</script>

<svelte:head>
	<title>Bib for {race.name} — Bibseller</title>
</svelte:head>

<nav class="text-sm">
	<a
		href={resolve('/races/[slug]', { slug: race.slug })}
		class="text-slate-500 hover:text-slate-900"
	>
		← {race.name}
	</a>
</nav>

<div class="mt-4 rounded-lg border border-slate-200 bg-white p-6 {available ? '' : 'opacity-75'}">
	<div class="flex flex-wrap items-start justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold">Bib for {race.name}</h1>
			<p class="mt-1 text-slate-600">
				{formatDate(race.event_date)} · {race.city}, {race.country}
				{#if race.distance}· {race.distance}{/if}
			</p>
		</div>
		<PolicyBadge policy={race.transfer_policy} />
	</div>

	{#if !available}
		<div class="mt-4 rounded-md bg-slate-100 p-3 text-sm font-semibold text-slate-700">
			This listing is no longer available ({listing.status}).
		</div>
	{/if}

	<div class="mt-6 flex items-baseline gap-3">
		<span class="text-4xl font-extrabold tracking-tight">{price ?? 'Price on request'}</span>
		{#if belowFace && original}
			<span class="text-lg text-slate-400 line-through">{original}</span>
			<span class="rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-semibold text-emerald-800">
				below face value
			</span>
		{/if}
	</div>

	{#if listing.description}
		<p class="mt-4 max-w-prose text-slate-700">{listing.description}</p>
	{/if}
	<p class="mt-4 text-sm text-slate-500">
		Listed by {listing.seller_name} on {formatDate(listing.created_at.slice(0, 10))}
	</p>

	{#if available}
		<div class="mt-6">
			<ListingCTA policy={race.transfer_policy} officialUrl={race.official_transfer_url} />
		</div>
	{/if}
</div>

<div class="mt-6">
	<PolicyCallout policy={race.transfer_policy} officialUrl={race.official_transfer_url} />
</div>
