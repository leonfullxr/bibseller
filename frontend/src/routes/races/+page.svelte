<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const countries = ['AT', 'BE', 'DE', 'ES', 'FR', 'IT', 'NL', 'PL', 'PT'];
	const sports = ['running', 'trail', 'triathlon', 'cycling', 'obstacle', 'other'];
	const policies = [
		{ value: 'platform_sale', label: 'Resale allowed' },
		{ value: 'official_only', label: 'Official transfer' },
		{ value: 'connect_only', label: 'Chat only' },
		{ value: 'unknown', label: 'Unverified' }
	];

	const nextQuery = $derived.by(() => {
		if (!data.nextCursor) return null;
		const params = new SvelteURLSearchParams(page.url.searchParams);
		params.set('cursor', data.nextCursor);
		return params.toString();
	});
</script>

<svelte:head>
	<title>Browse races — Bibseller</title>
	<meta name="description" content="Find race bibs for sale across EU running events." />
</svelte:head>

<h1>Browse races</h1>

<form method="GET" action={resolve('/races')} class="filters">
	<label>
		Search
		<input type="search" name="q" value={data.filters.q} placeholder="Race or city…" />
	</label>
	<label>
		Country
		<select name="country" value={data.filters.country}>
			<option value="">All</option>
			{#each countries as c (c)}<option value={c}>{c}</option>{/each}
		</select>
	</label>
	<label>
		Sport
		<select name="sport" value={data.filters.sport} class="sport">
			<option value="">All</option>
			{#each sports as s (s)}<option value={s}>{s}</option>{/each}
		</select>
	</label>
	<label>
		Transfer policy
		<select name="policy" value={data.filters.policy}>
			<option value="">All</option>
			{#each policies as p (p.value)}<option value={p.value}>{p.label}</option>{/each}
		</select>
	</label>
	<button type="submit">Filter</button>
</form>

{#if data.races.length === 0}
	<div class="empty">
		<p>No races match those filters.</p>
		<a href={resolve('/races')}>Clear filters</a>
	</div>
{:else}
	<div class="grid">
		{#each data.races as race (race.id)}
			<RaceCard {race} />
		{/each}
	</div>
	{#if nextQuery}
		<div class="more">
			<a href="{resolve('/races')}?{nextQuery}">Next page →</a>
		</div>
	{/if}
{/if}

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.filters {
		margin-top: 1rem;
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		gap: 0.75rem;
	}

	.filters label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		color: var(--slate-600);
	}

	.filters input,
	.filters select {
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.filters input {
		width: 11rem;
		padding: 0.375rem 0.625rem;
	}

	.filters select {
		padding: 0.375rem 0.5rem;
	}

	.sport {
		text-transform: capitalize;
	}

	.filters button {
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.375rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: white;
	}

	.filters button:hover {
		background: var(--slate-700);
	}

	.empty {
		margin-top: 3rem;
		border-radius: 0.5rem;
		border: 1px dashed var(--slate-300);
		padding: 2.5rem;
		text-align: center;
	}

	.empty p {
		font-weight: 500;
		color: var(--slate-600);
	}

	.empty a {
		margin-top: 0.5rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.grid {
		margin-top: 1.5rem;
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

	.more {
		margin-top: 2rem;
		text-align: center;
	}

	.more a {
		display: inline-block;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: var(--slate-700);
	}

	.more a:hover {
		background: white;
	}
</style>
