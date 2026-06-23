<script lang="ts">
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, link } = getI18n();

	const modes = $derived([
		{ name: t('home.modePlatformSaleName'), desc: t('home.modePlatformSaleDesc') },
		{ name: t('home.modeOfficialName'), desc: t('home.modeOfficialDesc') },
		{ name: t('home.modeConnectName'), desc: t('home.modeConnectDesc') }
	]);

	const steps = $derived([
		{ n: 1, title: t('home.step1Title'), desc: t('home.step1Desc') },
		{ n: 2, title: t('home.step2Title'), desc: t('home.step2Desc') },
		{ n: 3, title: t('home.step3Title'), desc: t('home.step3Desc') },
		{ n: 4, title: t('home.step4Title'), desc: t('home.step4Desc') }
	]);
</script>

<svelte:head>
	<title>{t('home.title')}</title>
	<meta name="description" content={t('home.metaDescription')} />
</svelte:head>

<section class="hero">
	<h1>{t('home.heroTitle')} <span>{t('home.heroTitleHighlight')}</span></h1>
	<p class="tagline">{t('home.tagline')}</p>

	<form method="GET" action={link(resolve('/races'))} class="search">
		<input type="search" name="q" placeholder={t('home.searchPlaceholder')} />
		<button type="submit">{t('home.search')}</button>
	</form>
	<a href={link(resolve('/races'))} class="browse-all">{t('home.browseAll')}</a>

	{#if data.apiStatus !== 'ok'}
		<div class="api-status">
			<span class="dot"></span>
			<span class="api-msg">
				{t('home.apiUnreachable')} <code>make dev</code>
			</span>
		</div>
	{/if}
</section>

<section class="how" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="cycle">
		{#each steps as step (step.n)}
			<li>
				<span class="num" aria-hidden="true">{step.n}</span>
				<h3>{step.title}</h3>
				<p>{step.desc}</p>
			</li>
		{/each}
	</ol>
	<p class="how-note">{t('home.howNote')}</p>
</section>

{#if data.upcoming.length > 0}
	<section class="upcoming">
		<div class="upcoming-head">
			<h2>{t('home.upcoming')}</h2>
			<a href={link(resolve('/races'))}>{t('home.seeAll')}</a>
		</div>
		<div class="grid">
			{#each data.upcoming as race (race.id)}
				<RaceCard {race} />
			{/each}
		</div>
	</section>
{/if}

<section class="modes">
	{#each modes as mode (mode.name)}
		<div class="mode">
			<h2>{mode.name}</h2>
			<p>{mode.desc}</p>
		</div>
	{/each}
</section>

<p class="construction">
	{t('home.underConstruction')}
	<a href="https://github.com/leonfullxr/bibseller/issues/13" rel="external">{t('home.roadmap')}</a
	>.
</p>

<style>
	.hero {
		padding-block: 2rem;
		text-align: center;
	}

	.hero h1 {
		font-size: 2.25rem;
		line-height: 2.5rem;
		font-weight: 800;
		letter-spacing: -0.025em;
	}

	@media (min-width: 640px) {
		.hero h1 {
			font-size: 3rem;
			line-height: 1;
		}
	}

	.hero h1 span {
		color: var(--emerald-600);
	}

	.tagline {
		margin: 1rem auto 0;
		max-width: 42rem;
		font-size: 1.125rem;
		line-height: 1.75rem;
		color: var(--slate-600);
	}

	.search {
		margin: 2rem auto 0;
		display: flex;
		max-width: 28rem;
		align-items: center;
		gap: 0.5rem;
	}

	.search input {
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.search button {
		border-radius: 0.375rem;
		background: var(--emerald-600);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		white-space: nowrap;
		color: white;
	}

	.search button:hover {
		background: var(--emerald-700);
	}

	.browse-all {
		margin-top: 0.75rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.api-status {
		margin-top: 1.5rem;
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

	.upcoming {
		padding-block: 2rem;
	}

	.upcoming-head {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
	}

	.upcoming-head h2 {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.upcoming-head a {
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--emerald-700);
		text-decoration: underline;
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

	.modes {
		display: grid;
		gap: 1rem;
		padding-block: 2rem;
	}

	@media (min-width: 640px) {
		.modes {
			grid-template-columns: repeat(3, minmax(0, 1fr));
		}
	}

	.mode {
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.25rem;
	}

	.mode h2 {
		font-weight: 600;
	}

	.mode p {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.construction {
		padding-block: 1rem;
		text-align: center;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-400);
	}

	.construction a {
		text-decoration: underline;
	}

	.construction a:hover {
		color: var(--slate-600);
	}

	.how {
		padding-block: 2rem;
		text-align: center;
	}

	.how h2 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
		letter-spacing: -0.015em;
	}

	/* Mobile-first: a vertical numbered timeline. */
	.cycle {
		list-style: none;
		margin: 1.5rem auto 0;
		padding: 0;
		display: grid;
		gap: 1.25rem;
		max-width: 22rem;
		text-align: left;
	}

	.cycle li {
		position: relative;
		padding-left: 2.75rem;
	}

	.cycle .num {
		position: absolute;
		left: 0;
		top: 0;
		display: grid;
		place-items: center;
		width: 1.75rem;
		height: 1.75rem;
		border-radius: 9999px;
		background: var(--emerald-600);
		font-size: 0.875rem;
		font-weight: 700;
		color: white;
	}

	/* connecting line down to the next badge */
	.cycle li:not(:last-child)::after {
		content: '';
		position: absolute;
		left: 0.8125rem;
		top: 1.75rem;
		bottom: -1.25rem;
		width: 2px;
		background: var(--slate-200);
	}

	.cycle h3 {
		font-size: 1rem;
		font-weight: 600;
	}

	.cycle p {
		margin-top: 0.25rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.how-note {
		margin: 1.5rem auto 0;
		max-width: 32rem;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-400);
	}

	/* Desktop: arrange the four steps clockwise around a ring (the "cycle"). */
	/* ponytail: 4 fixed positions, not trig - revisit only if the step count changes. */
	@media (min-width: 640px) {
		.cycle {
			position: relative;
			display: block;
			width: min(78vw, 30rem);
			height: min(78vw, 30rem);
			max-width: none;
			margin: 2.5rem auto 0;
			text-align: center;
		}

		.cycle::before {
			content: '';
			position: absolute;
			inset: 19%;
			border: 2px dashed var(--slate-300);
			border-radius: 9999px;
		}

		.cycle::after {
			content: '\21BB';
			position: absolute;
			inset: 0;
			display: grid;
			place-items: center;
			font-size: 2.25rem;
			color: var(--emerald-600);
		}

		.cycle li {
			position: absolute;
			width: 10rem;
			padding-left: 0;
		}

		.cycle li:not(:last-child)::after {
			content: none;
		}

		.cycle .num {
			position: static;
			margin: 0 auto 0.5rem;
		}

		.cycle li:nth-child(1) {
			top: 0;
			left: 50%;
			transform: translate(-50%, 0);
		}

		.cycle li:nth-child(2) {
			top: 50%;
			right: 0;
			transform: translate(0, -50%);
		}

		.cycle li:nth-child(3) {
			bottom: 0;
			left: 50%;
			transform: translate(-50%, 0);
		}

		.cycle li:nth-child(4) {
			top: 50%;
			left: 0;
			transform: translate(0, -50%);
		}
	}
</style>
