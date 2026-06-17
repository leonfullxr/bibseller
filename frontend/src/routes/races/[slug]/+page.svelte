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
	<title>{race.name} - bibs for sale - Bibseller</title>
	<meta
		name="description"
		content="Bibs for {race.name} ({formatDate(race.event_date)}, {race.city})."
	/>
</svelte:head>

<nav>
	<a href={resolve('/races')}>Back to all races</a>
</nav>

<header>
	<div class="title-row">
		<h1>{race.name}</h1>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="meta">
		{formatDate(race.event_date)} - {race.city}, {race.country}
		{#if race.distance}
			- {race.distance}{/if}
		- <span class="sport">{race.sport}</span>
	</p>
	{#if race.website_url}
		<a href={race.website_url} rel="external nofollow noopener" target="_blank" class="website">
			Race website
		</a>
	{/if}
</header>

<div class="callout-wrap">
	<PolicyCallout
		policy={race.transfer_policy}
		officialUrl={race.official_transfer_url}
		notes={race.policy_notes}
	/>
</div>

<section>
	<div class="section-head">
		<h2>
			{race.active_listings}
			{race.active_listings === 1 ? 'bib' : 'bibs'} for sale
		</h2>
		<a href={resolve('/sell/[slug]', { slug: race.slug })} class="sell-cta">Sell your bib</a>
	</div>
	{#if data.listings.length === 0}
		<div class="empty">
			<p>No bibs listed for this race yet.</p>
			<p class="hint">Selling yours? Listing opens soon.</p>
		</div>
	{:else}
		<div class="grid">
			{#each data.listings as listing (listing.id)}
				<ListingCard {listing} />
			{/each}
		</div>
	{/if}
</section>

<style>
	.section-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
	}

	.sell-cta {
		border-radius: 0.375rem;
		background: var(--emerald-600);
		padding: 0.375rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
		white-space: nowrap;
	}

	.sell-cta:hover {
		background: var(--emerald-700);
	}

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

	header {
		margin-top: 1rem;
	}

	.title-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.75rem;
	}

	h1 {
		font-size: 1.875rem;
		line-height: 2.25rem;
		font-weight: 700;
		letter-spacing: -0.025em;
	}

	.meta {
		margin-top: 0.5rem;
		color: var(--slate-600);
	}

	.sport {
		text-transform: capitalize;
	}

	.website {
		margin-top: 0.25rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.callout-wrap {
		margin-top: 1.5rem;
	}

	section {
		margin-top: 2rem;
	}

	h2 {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.empty {
		margin-top: 1rem;
		border-radius: 0.5rem;
		border: 1px dashed var(--slate-300);
		padding: 2.5rem;
		text-align: center;
	}

	.empty p {
		font-weight: 500;
		color: var(--slate-600);
	}

	.empty .hint {
		margin-top: 0.25rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 400;
		color: var(--slate-500);
	}

	.grid {
		margin-top: 1rem;
		display: grid;
		gap: 1rem;
	}

	@media (min-width: 640px) {
		.grid {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	@media (min-width: 1024px) {
		.grid {
			grid-template-columns: repeat(3, minmax(0, 1fr));
		}
	}
</style>
