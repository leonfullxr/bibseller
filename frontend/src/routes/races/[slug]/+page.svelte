<script lang="ts">
	import { resolve } from '$app/paths';
	import ListingCard from '$lib/components/ListingCard.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate } from '$lib/format';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const race = $derived(data.race);
</script>

<svelte:head>
	<title>{race.name} — bibs for sale — Bibseller</title>
	<meta
		name="description"
		content="Bibs for {race.name} ({formatDate(race.event_date)}, {race.city})."
	/>
</svelte:head>

<nav class="text-sm">
	<a href={resolve('/races')} class="text-slate-500 hover:text-slate-900">← All races</a>
</nav>

<header class="mt-4">
	<div class="flex flex-wrap items-center gap-3">
		<h1 class="text-3xl font-bold tracking-tight">{race.name}</h1>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="mt-2 text-slate-600">
		{formatDate(race.event_date)} · {race.city}, {race.country}
		{#if race.distance}· {race.distance}{/if}
		· <span class="capitalize">{race.sport}</span>
	</p>
	{#if race.website_url}
		<a
			href={race.website_url}
			rel="external nofollow noopener"
			target="_blank"
			class="mt-1 inline-block text-sm text-emerald-700 underline"
		>
			Race website ↗
		</a>
	{/if}
</header>

<div class="mt-6">
	<PolicyCallout
		policy={race.transfer_policy}
		officialUrl={race.official_transfer_url}
		notes={race.policy_notes}
	/>
</div>

<section class="mt-8">
	<h2 class="text-lg font-semibold">
		{race.active_listings}
		{race.active_listings === 1 ? 'bib' : 'bibs'} for sale
	</h2>
	{#if data.listings.length === 0}
		<div class="mt-4 rounded-lg border border-dashed border-slate-300 p-10 text-center">
			<p class="font-medium text-slate-600">No bibs listed for this race yet.</p>
			<p class="mt-1 text-sm text-slate-500">Selling yours? Listing opens soon.</p>
		</div>
	{:else}
		<div class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each data.listings as listing (listing.id)}
				<ListingCard {listing} />
			{/each}
		</div>
	{/if}
</section>
