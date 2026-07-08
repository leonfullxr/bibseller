<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import RaceMap from '$lib/components/RaceMap.svelte';
	import { formatDate } from '$lib/format';
	import { getI18n } from '$lib/i18n';
	import { sportLabel } from '$lib/i18n/messages';
	import { transferPolicies } from '$lib/policy';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, plural, locale, link } = getI18n();

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
	const sports = ['running', 'trail', 'triathlon', 'cycling', 'obstacle', 'other'] as const;
	// ponytail: hardcoded distinct distances of the curated catalog (#6 intake);
	// derive from a server-side aggregate when the catalog outgrows this list.
	const knownDistances = ['20k', '42k', '70.3', 'half', 'marathon'];
	const distances = $derived.by(() => {
		const list = [...knownDistances];
		if (data.filters.distance && !list.includes(data.filters.distance)) {
			list.push(data.filters.distance);
			list.sort();
		}
		return list;
	});
	const policies = $derived(
		transferPolicies.map((value) => ({ value, label: t(`policy.label.${value}`) }))
	);

	const hasMap = $derived(Object.keys(data.countryCounts).length > 0);

	// Active filters as removable chips; each chip links to the same URL minus
	// its own param (and minus the cursor - removing a filter restarts paging).
	const chips = $derived.by(() => {
		const f = data.filters;
		const out: { key: string; label: string }[] = [];
		if (f.q) out.push({ key: 'q', label: `“${f.q}”` });
		if (f.country) out.push({ key: 'country', label: countryNames.of(f.country) ?? f.country });
		if (f.sport) out.push({ key: 'sport', label: sportLabel(t, f.sport) });
		if (f.distance) out.push({ key: 'distance', label: f.distance });
		if (f.policy)
			out.push({
				key: 'policy',
				label: policies.find((p) => p.value === f.policy)?.label ?? f.policy
			});
		if (f.date_from)
			out.push({
				key: 'date_from',
				label: `${t('races.filter.dateFrom')} ${formatDate(f.date_from, locale)}`
			});
		if (f.date_to)
			out.push({
				key: 'date_to',
				label: `${t('races.filter.dateTo')} ${formatDate(f.date_to, locale)}`
			});
		return out;
	});
	const anyFilter = $derived(chips.length > 0);

	// Query string for the current URL minus one filter; the caller prefixes
	// link(resolve('/races')) inline so the lint rule can see the resolve().
	function without(key: string): string {
		const params = new SvelteURLSearchParams(page.url.searchParams);
		params.delete(key);
		params.delete('cursor');
		const qs = params.toString();
		return qs ? `?${qs}` : '';
	}

	// Selects and date inputs apply on change when JS is around; the Apply
	// button stays as the no-JS and keyboard path. Text search keeps its
	// Enter-to-submit and must not fire on blur.
	function autoSubmit(e: Event) {
		const el = e.target;
		if (el instanceof HTMLSelectElement || (el instanceof HTMLInputElement && el.type === 'date')) {
			el.form?.requestSubmit();
		}
	}

	// SSR renders the map open so no-JS readers always see it (the summary
	// toggle is hidden from 48rem up, where the map pane is permanent). On
	// mount, small screens collapse it to the summary toggle.
	let mapOpen = $state(true);
	onMount(() => {
		if (!matchMedia('(min-width: 48rem)').matches) mapOpen = false;
	});

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

<div class="layout" class:has-map={hasMap}>
	{#if hasMap}
		<details class="map-pane" bind:open={mapOpen}>
			<summary>{t('races.mapSummary')}</summary>
			<div class="panel map-panel">
				<RaceMap
					counts={data.countryCounts}
					cities={data.cities}
					country={data.filters.country}
					filters={data.filters}
				/>
			</div>
		</details>
	{/if}

	<div class="results">
		<form
			method="GET"
			action={link(resolve('/races'))}
			class="panel filterbar"
			onchange={autoSubmit}
		>
			<label class="search">
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
				{t('races.filter.distance')}
				<select name="distance" class="field" value={data.filters.distance}>
					<option value="">{t('races.filter.all')}</option>
					{#each distances as d (d)}<option value={d}>{d}</option>{/each}
				</select>
			</label>
			<label>
				{t('races.filter.policy')}
				<select name="policy" class="field" value={data.filters.policy}>
					<option value="">{t('races.filter.all')}</option>
					{#each policies as p (p.value)}<option value={p.value}>{p.label}</option>{/each}
				</select>
			</label>
			<label>
				{t('races.filter.dateFrom')}
				<input type="date" name="date_from" class="field" value={data.filters.date_from} />
			</label>
			<label>
				{t('races.filter.dateTo')}
				<input type="date" name="date_to" class="field" value={data.filters.date_to} />
			</label>
			<button type="submit" class="btn btn-primary">{t('races.filter.submit')}</button>
			{#if anyFilter}
				<a class="clear" href={link(resolve('/races'))}>{t('races.clearFilters')}</a>
			{/if}
		</form>

		<div class="meta">
			<p class="count" role="status">
				{data.nextCursor
					? t('races.resultCountMore', { n: String(data.races.length) })
					: plural('races.resultCount', data.races.length)}
			</p>
			{#if anyFilter}
				<ul class="chips">
					{#each chips as chip (chip.key)}
						<li>
							<a
								class="chip"
								href="{link(resolve('/races'))}{without(chip.key)}"
								aria-label={t('races.removeFilter', { name: chip.label })}
							>
								{chip.label}
								<svg
									viewBox="0 0 24 24"
									width="12"
									height="12"
									fill="none"
									stroke="currentColor"
									stroke-width="2.5"
									stroke-linecap="round"
									aria-hidden="true"
								>
									<path d="M18 6 6 18M6 6l12 12" />
								</svg>
							</a>
						</li>
					{/each}
				</ul>
			{/if}
		</div>

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
	</div>
</div>

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	/* Stacked by default (mobile); the desktop split grid is below. */
	.layout {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	/* Map pane: a collapsible <details> below 48rem, a permanent panel above
	   (summary hidden), the sticky left pane of the split view from 64rem up. */
	.map-pane summary {
		cursor: pointer;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: var(--slate-700);
	}

	.map-pane[open] summary {
		margin-bottom: 0.75rem;
	}

	.map-panel {
		padding: 1rem;
	}

	.map-panel :global(.race-map) {
		margin-top: 0;
	}

	.results {
		min-width: 0;
	}

	/* Horizontal filter bar: compact fields in a wrapping row; the Apply
	   button bottom-aligns with the fields (labels sit above each field). */
	.filterbar {
		display: flex;
		flex-wrap: wrap;
		align-items: end;
		gap: 0.75rem;
		padding: 1rem;
	}

	.filterbar label {
		display: flex;
		min-width: 0;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--slate-500);
	}

	.search {
		flex: 1 1 11rem;
	}

	.filterbar .field {
		max-width: 100%;
		padding: 0.375rem 0.625rem;
		font-weight: 400;
		text-transform: none;
		letter-spacing: normal;
	}

	.search .field {
		width: 100%;
	}

	.sport {
		text-transform: capitalize;
	}

	.filterbar button {
		padding-block: 0.4375rem;
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

	.meta {
		margin-top: 1rem;
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.5rem 0.75rem;
	}

	.count {
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.chips {
		list-style: none;
		padding: 0;
		display: flex;
		flex-wrap: wrap;
		gap: 0.375rem;
	}

	.chip {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		border-radius: 9999px;
		border: 1px solid var(--slate-200);
		background: var(--slate-100);
		color: var(--slate-700);
		padding: 0.125rem 0.625rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.chip:hover {
		background: var(--slate-200);
		color: var(--slate-900);
	}

	.empty {
		margin-top: 1.5rem;
	}

	.empty a {
		margin-top: 0.5rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.grid {
		margin-top: 1rem;
		display: grid;
		gap: 1rem;
	}

	.more {
		margin-top: 2rem;
		text-align: center;
	}

	/* 48rem-64rem: map permanent on top, shorter (~18rem tall); grid 2-3 cols. */
	@media (min-width: 48rem) {
		.map-pane summary {
			display: none;
		}

		.map-panel {
			/* RaceMap caps its width from this height budget. */
			--map-max-h: 18rem;
		}

		.grid {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	@media (min-width: 56rem) and (max-width: 63.9375rem) {
		.grid {
			grid-template-columns: repeat(3, minmax(0, 1fr));
		}
	}

	/* Desktop split view: sticky map pane left (~45%), scrolling results right. */
	@media (min-width: 64rem) {
		.layout.has-map {
			display: grid;
			grid-template-columns: minmax(0, 9fr) minmax(0, 11fr);
			gap: 1.5rem;
			align-items: start;
		}

		.map-pane {
			/* 4.5rem = sticky app bar (3.5rem) + breathing room. */
			position: sticky;
			top: 4.5rem;
		}

		.map-panel {
			display: flex;
			align-items: center;
			justify-content: center;
			height: calc(100dvh - 5.5rem);
			/* Panel height minus its padding and the map-hint line. */
			--map-max-h: max(14rem, calc(100dvh - 10.5rem));
		}

		.map-panel :global(.race-map) {
			width: 100%;
		}
	}
</style>
