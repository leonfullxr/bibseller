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

<h1 class="text-2xl font-bold">Browse races</h1>

<form method="GET" action={resolve('/races')} class="mt-4 flex flex-wrap items-end gap-3">
	<label class="flex flex-col gap-1 text-xs font-medium text-slate-600">
		Search
		<input
			type="search"
			name="q"
			value={data.filters.q}
			placeholder="Race or city…"
			class="w-44 rounded-md border border-slate-300 px-2.5 py-1.5 text-sm"
		/>
	</label>
	<label class="flex flex-col gap-1 text-xs font-medium text-slate-600">
		Country
		<select
			name="country"
			value={data.filters.country}
			class="rounded-md border border-slate-300 px-2 py-1.5 text-sm"
		>
			<option value="">All</option>
			{#each countries as c (c)}<option value={c}>{c}</option>{/each}
		</select>
	</label>
	<label class="flex flex-col gap-1 text-xs font-medium text-slate-600">
		Sport
		<select
			name="sport"
			value={data.filters.sport}
			class="rounded-md border border-slate-300 px-2 py-1.5 text-sm capitalize"
		>
			<option value="">All</option>
			{#each sports as s (s)}<option value={s}>{s}</option>{/each}
		</select>
	</label>
	<label class="flex flex-col gap-1 text-xs font-medium text-slate-600">
		Transfer policy
		<select
			name="policy"
			value={data.filters.policy}
			class="rounded-md border border-slate-300 px-2 py-1.5 text-sm"
		>
			<option value="">All</option>
			{#each policies as p (p.value)}<option value={p.value}>{p.label}</option>{/each}
		</select>
	</label>
	<button
		type="submit"
		class="rounded-md bg-slate-900 px-4 py-1.5 text-sm font-semibold text-white hover:bg-slate-700"
	>
		Filter
	</button>
</form>

{#if data.races.length === 0}
	<div class="mt-12 rounded-lg border border-dashed border-slate-300 p-10 text-center">
		<p class="font-medium text-slate-600">No races match those filters.</p>
		<a href={resolve('/races')} class="mt-2 inline-block text-sm text-emerald-700 underline">
			Clear filters
		</a>
	</div>
{:else}
	<div class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
		{#each data.races as race (race.id)}
			<RaceCard {race} />
		{/each}
	</div>
	{#if nextQuery}
		<div class="mt-8 text-center">
			<a
				href="{resolve('/races')}?{nextQuery}"
				class="inline-block rounded-md border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-white"
			>
				Next page →
			</a>
		</div>
	{/if}
{/if}
