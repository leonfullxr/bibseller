<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import RaceMap from '$lib/components/RaceMap.svelte';
	import { getI18n } from '$lib/i18n';
	import { transferPolicies } from '$lib/policy';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, locale, link } = getI18n();

	// Country options come from the live per-country counts; the hardcoded
	// catalog set is the fallback when the map-counts fetch degraded to empty.
	const fallbackCountries = ['AT', 'BE', 'DE', 'ES', 'FR', 'IT', 'NL', 'PL', 'PT'];
	const countries = $derived.by(() => {
		const live = Object.keys(data.countryCounts).sort();
		const list = live.length ? live : [...fallbackCountries];
		// Keep the active filter selectable even if it has no upcoming races.
		if (data.filters.country && !list.includes(data.filters.country)) {
			list.push(data.filters.country);
			list.sort();
		}
		return list;
	});
	// Localized country names for the filter labels; fall back to the code.
	const countryNames = $derived(new Intl.DisplayNames([locale], { type: 'region' }));
	const anyFilter = $derived(
		Boolean(data.filters.q || data.filters.country || data.filters.sport || data.filters.policy)
	);
	const sports = ['running', 'trail', 'triathlon', 'cycling', 'obstacle', 'other'] as const;
	const policies = $derived(
		transferPolicies.map((value) => ({ value, label: t(`policy.label.${value}`) }))
	);

	const nextQuery = $derived.by(() => {
		if (!data.nextCursor) return null;
		const params = new SvelteURLSearchParams(page.url.searchParams);
		params.set('cursor', data.nextCursor);
		return params.toString();
	});
</script>

<svelte:head>
	<title>{t('races.title')}</title>
	<meta name="description" content={t('races.metaDescription')} />
</svelte:head>

<h1>{t('races.heading')}</h1>

{#if Object.keys(data.countryCounts).length > 0}
	<RaceMap
		counts={data.countryCounts}
		cities={data.cities}
		country={data.filters.country}
		filters={data.filters}
	/>
{/if}

<form method="GET" action={link(resolve('/races'))} class="filters">
	<label>
		{t('races.filter.search')}
		<input
			type="search"
			name="q"
			class="field"
			value={data.filters.q}
			placeholder={t('races.filter.searchPlaceholder')}
		/>
	</label>
	<label>
		{t('races.filter.country')}
		<select name="country" class="field" value={data.filters.country}>
			<option value="">{t('races.filter.all')}</option>
			{#each countries as c (c)}<option value={c}>{countryNames.of(c) ?? c}</option>{/each}
		</select>
	</label>
	<label>
		{t('races.filter.sport')}
		<select name="sport" value={data.filters.sport} class="field sport">
			<option value="">{t('races.filter.all')}</option>
			{#each sports as s (s)}<option value={s}>{t(`sport.${s}`)}</option>{/each}
		</select>
	</label>
	<label>
		{t('races.filter.policy')}
		<select name="policy" class="field" value={data.filters.policy}>
			<option value="">{t('races.filter.all')}</option>
			{#each policies as p (p.value)}<option value={p.value}>{p.label}</option>{/each}
		</select>
	</label>
	<button type="submit" class="btn btn-primary">{t('races.filter.submit')}</button>
	{#if anyFilter}
		<a class="clear" href={link(resolve('/races'))}>{t('races.clearFilters')}</a>
	{/if}
</form>

{#if data.races.length === 0}
	<div class="empty">
		<p>{t('races.empty')}</p>
		<a href={link(resolve('/races'))}>{t('races.clearFilters')}</a>
	</div>
{:else}
	<div class="grid">
		{#each data.races as race (race.id)}
			<RaceCard {race} />
		{/each}
	</div>
	{#if nextQuery}
		<div class="more">
			<a class="btn btn-outline" href="{link(resolve('/races'))}?{nextQuery}"
				>{t('races.nextPage')}</a
			>
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

	/* Compact sizing for the filter bar; the visual style comes from .field. */
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
		padding: 0.375rem 1rem;
	}

	.clear {
		align-self: center;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
		text-decoration: underline;
	}

	.clear:hover {
		color: var(--slate-900);
	}

	.empty {
		margin-top: 3rem;
	}

	.empty a {
		margin-top: 0.5rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
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
</style>
