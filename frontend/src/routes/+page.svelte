<script lang="ts">
	import { afterNavigate } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import Icon from '$lib/components/Icon.svelte';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import type { Page, RaceSummary } from '$lib/api/types';
	import { todayISO } from '$lib/format';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, plural, locale, link } = getI18n();

	// Instant search (progressive enhancement): typing fetches results directly
	// and overrides the grid; null means "show the server-loaded results". The
	// no-JS path is the plain GET form. Navigations (pill clicks, form submits)
	// hand authority back to the server data.
	let live: RaceSummary[] | null = $state(null);
	const races = $derived(live ?? data.races);
	afterNavigate(() => (live = null));

	let timer: ReturnType<typeof setTimeout> | undefined;
	let ctrl: AbortController | undefined;
	function instantSearch(e: Event) {
		const q = (e.currentTarget as HTMLInputElement).value.trim();
		clearTimeout(timer);
		timer = setTimeout(async () => {
			ctrl?.abort();
			ctrl = new AbortController();
			const params = new SvelteURLSearchParams();
			if (q) params.set('q', q);
			if (data.filters.country) params.set('country', data.filters.country);
			if (data.filters.sport) params.set('sport', data.filters.sport);
			params.set('date_from', todayISO());
			params.set('limit', '12');
			try {
				const res = await fetch(`/api/v1/races?${params}`, { signal: ctrl.signal });
				if (res.ok) live = ((await res.json()) as Page<RaceSummary>).items;
			} catch {
				// Aborted or offline: keep whatever is on screen.
			}
		}, 300);
	}

	// Country pills come from the live per-country counts (empty when that
	// fetch degraded); keep an active filter visible so it can be toggled off.
	const countries = $derived.by(() => {
		const list = Object.keys(data.countryCounts).sort();
		if (data.filters.country && !list.includes(data.filters.country)) {
			list.push(data.filters.country);
			list.sort();
		}
		return list;
	});
	const countryNames = $derived(new Intl.DisplayNames([locale], { type: 'region' }));
	const sports = ['running', 'trail', 'triathlon', 'cycling', 'obstacle', 'other'] as const;

	// Query string for the current URL with one pill's param toggled; the caller
	// prefixes link(resolve('/')) inline so the lint rule can see the resolve().
	function toggled(key: 'country' | 'sport', value: string): string {
		const params = new SvelteURLSearchParams(page.url.searchParams);
		if (params.get(key) === value) params.delete(key);
		else params.set(key, value);
		const qs = params.toString();
		return qs ? `?${qs}` : '';
	}

	// Carry the quick filters into the full catalog.
	const catalogQuery = $derived.by(() => {
		const qs = page.url.searchParams.toString();
		return qs ? `?${qs}` : '';
	});

	const steps = $derived([
		{ icon: 'list', title: t('home.step1Title'), desc: t('home.step1Desc') },
		{ icon: 'chat', title: t('home.step2Title'), desc: t('home.step2Desc') },
		{ icon: 'transfer', title: t('home.step3Title'), desc: t('home.step3Desc') },
		{ icon: 'check', title: t('home.step4Title'), desc: t('home.step4Desc') }
	]);
</script>

<svelte:head>
	<title>{t('home.title')}</title>
	<meta name="description" content={t('home.metaDescription')} />
</svelte:head>

<section class="intro">
	<h1>{t('home.heroTitle')} <span>{t('home.heroTitleHighlight')}</span></h1>
	<p class="tagline">{t('home.tagline')}</p>

	<form method="GET" action={link(resolve('/'))} class="search" role="search">
		{#if data.filters.country}
			<input type="hidden" name="country" value={data.filters.country} />
		{/if}
		{#if data.filters.sport}
			<input type="hidden" name="sport" value={data.filters.sport} />
		{/if}
		<input
			type="search"
			name="q"
			class="search-input"
			value={data.filters.q}
			placeholder={t('home.searchPlaceholder')}
			aria-label={t('races.filter.search')}
			oninput={instantSearch}
		/>
		<button type="submit" class="btn btn-primary">{t('home.search')}</button>
	</form>

	{#if data.apiStatus !== 'ok'}
		<div class="api-status">
			<span class="dot"></span>
			<span class="api-msg">
				{t('home.apiUnreachable')} <code>make dev</code>
			</span>
		</div>
	{/if}
</section>

<nav class="pills" aria-label={t('races.filtersSummary')}>
	{#each countries as c (c)}
		<a
			class="qpill"
			class:active={data.filters.country === c}
			href="{link(resolve('/'))}{toggled('country', c)}"
		>
			{countryNames.of(c) ?? c}
		</a>
	{/each}
	{#each sports as s (s)}
		<a
			class="qpill"
			class:active={data.filters.sport === s}
			href="{link(resolve('/'))}{toggled('sport', s)}"
		>
			{t(`sport.${s}`)}
		</a>
	{/each}
</nav>

<section class="results">
	<div class="meta">
		<p class="count" role="status">{plural('races.resultCount', races.length)}</p>
		<a class="more-filters" href="{link(resolve('/races'))}{catalogQuery}"
			>{t('home.moreFilters')}</a
		>
	</div>

	{#if races.length === 0}
		<div class="empty">
			<p>{t('races.empty')}</p>
			<a href={link(resolve('/'))}>{t('races.clearFilters')}</a>
		</div>
	{:else}
		<div class="grid">
			{#each races as race (race.id)}
				<RaceCard {race} />
			{/each}
		</div>
	{/if}
</section>

<section class="how panel" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="steps">
		{#each steps as step (step.icon)}
			<li>
				<span class="step-icon" aria-hidden="true"><Icon name={step.icon} /></span>
				<div>
					<h3>{step.title}</h3>
					<p>{step.desc}</p>
				</div>
			</li>
		{/each}
	</ol>
	<p class="contact-line">
		{t('home.contactLead')}
		<a href={link(resolve('/contact'))}>{t('home.contactCta')}</a>
	</p>
</section>

<p class="construction">
	{t('home.underConstruction')}
	<a href="https://github.com/leonfullxr/bibseller/issues/13" rel="external">{t('home.roadmap')}</a
	>.
</p>

<style>
	/* Compact intro: the page is an app surface, the search is the hero. */
	.intro {
		padding: 2rem 0 0.5rem;
		text-align: center;
	}

	.intro h1 {
		margin-inline: auto;
		max-width: 40rem;
		font-size: 1.875rem;
		line-height: 1.15;
	}

	@media (min-width: 640px) {
		.intro h1 {
			font-size: 2.25rem;
		}
	}

	.intro h1 span {
		color: var(--brand-600);
	}

	.tagline {
		margin: 0.625rem auto 0;
		max-width: 36rem;
		font-size: 0.9375rem;
		line-height: 1.5rem;
		color: var(--slate-600);
	}

	/* The large search field. */
	.search {
		margin: 1.5rem auto 0;
		display: flex;
		max-width: 36rem;
		align-items: stretch;
		gap: 0.5rem;
	}

	.search-input {
		width: 100%;
		border-radius: 0.75rem;
		border: 1px solid var(--slate-300);
		background: white;
		color: var(--ink);
		padding: 0.8125rem 1.125rem;
		font-size: 1.0625rem;
		line-height: 1.5rem;
		box-shadow: var(--shadow-hard-sm);
	}

	/* Border swap for any focus; the global :focus-visible ring still applies
	   for keyboard users. */
	.search-input:focus {
		border-color: var(--brand-600);
	}

	.search .btn {
		border-radius: 0.75rem;
		padding-inline: 1.5rem;
		white-space: nowrap;
	}

	.api-status {
		margin-top: 1rem;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.dot {
		display: inline-block;
		height: 0.5rem;
		width: 0.5rem;
		border-radius: 9999px;
		background: var(--amber-500);
	}

	.api-msg {
		color: var(--slate-600);
	}

	.api-msg code {
		border-radius: 0.25rem;
		background: var(--slate-200);
		padding: 0.125rem 0.25rem;
	}

	/* Quick filter pills: one horizontally scrollable row. */
	.pills {
		margin-top: 1.25rem;
		display: flex;
		gap: 0.5rem;
		overflow-x: auto;
		scrollbar-width: none; /* pills scroll by touch/wheel; the bar is noise */
		padding-bottom: 0.375rem; /* room so the h-scrollbar never overlaps the pills */
	}

	.qpill {
		flex-shrink: 0;
		border-radius: 9999px;
		border: 1px solid var(--slate-300);
		background: white;
		color: var(--slate-700);
		padding: 0.3125rem 0.875rem;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		font-weight: 600;
		white-space: nowrap;
		transition:
			background-color 0.15s,
			border-color 0.15s,
			color 0.15s;
	}

	.qpill:hover {
		border-color: var(--slate-400);
		background: var(--slate-50);
	}

	.qpill.active {
		border-color: var(--brand-600);
		background: var(--brand-600);
		color: white;
	}

	.qpill.active:hover {
		border-color: var(--brand-700);
		background: var(--brand-700);
	}

	/* Live results. */
	.results {
		margin-top: 0.75rem;
	}

	.meta {
		display: flex;
		flex-wrap: wrap;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.5rem 0.75rem;
	}

	.count {
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.more-filters {
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: var(--brand-700);
	}

	.more-filters:hover {
		color: var(--brand-800);
		text-decoration: underline;
	}

	.empty {
		margin-top: 1rem;
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

	/* Marketing strip: the four steps condensed to one compact row. */
	.how {
		margin-top: 2.5rem;
	}

	.how h2 {
		font-size: 1.125rem;
		line-height: 1.5rem;
	}

	.steps {
		list-style: none;
		margin: 1rem 0 0;
		padding: 0;
		display: grid;
		gap: 1rem;
	}

	@media (min-width: 720px) {
		.steps {
			grid-template-columns: repeat(4, minmax(0, 1fr));
		}
	}

	.steps li {
		display: flex;
		gap: 0.625rem;
	}

	.step-icon {
		display: grid;
		place-items: center;
		flex-shrink: 0;
		width: 2rem;
		height: 2rem;
		border-radius: 9999px;
		background: var(--brand-50);
		color: var(--brand-700);
		font-size: 1rem;
	}

	.steps h3 {
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
	}

	.steps p {
		margin-top: 0.125rem;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.contact-line {
		margin-top: 1.25rem;
		border-top: 1px solid var(--slate-200);
		padding-top: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.contact-line a {
		font-weight: 600;
		color: var(--brand-700);
		text-decoration: underline;
	}

	.contact-line a:hover {
		color: var(--brand-800);
	}

	.construction {
		padding-block: 1.5rem 1rem;
		text-align: center;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	.construction a {
		text-decoration: underline;
	}

	.construction a:hover {
		color: var(--slate-600);
	}
</style>
