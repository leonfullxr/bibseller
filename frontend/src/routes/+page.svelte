<script lang="ts">
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, link } = getI18n();

	const steps = $derived([
		{ icon: 'list', title: t('home.step1Title'), desc: t('home.step1Desc') },
		{ icon: 'chat', title: t('home.step2Title'), desc: t('home.step2Desc') },
		{ icon: 'transfer', title: t('home.step3Title'), desc: t('home.step3Desc') },
		{ icon: 'check', title: t('home.step4Title'), desc: t('home.step4Desc') }
	]);

	const journey = $derived([
		{ n: 1, who: 'seller', icon: 'list', label: t('home.j1Title') },
		{ n: 2, who: 'buyer', icon: 'search', label: t('home.j2Title') },
		{ n: 3, who: 'seller', icon: 'chat', label: t('home.j3Title') },
		{ n: 4, who: 'buyer', icon: 'transfer', label: t('home.j4Title') },
		{ n: 5, who: 'seller', icon: 'handover', label: t('home.j5Title') },
		{ n: 6, who: 'buyer', icon: 'medal', label: t('home.j6Title') }
	]);
	const sellerSteps = $derived(journey.filter((s) => s.who === 'seller'));
	const buyerSteps = $derived(journey.filter((s) => s.who === 'buyer'));

	const modes = $derived([
		{ name: t('home.modePlatformSaleName'), desc: t('home.modePlatformSaleDesc') },
		{ name: t('home.modeOfficialName'), desc: t('home.modeOfficialDesc') },
		{ name: t('home.modeConnectName'), desc: t('home.modeConnectDesc') }
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

<section class="how" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="how-steps">
		{#each steps as step (step.title)}
			<li class="how-step">
				<span class="how-icon"><Icon name={step.icon} /></span>
				<h3>{step.title}</h3>
				<p>{step.desc}</p>
			</li>
		{/each}
	</ol>
	<p class="how-note">{t('home.howNote')}</p>
</section>

<section class="journey" aria-labelledby="journey-title">
	<h2 id="journey-title">{t('home.journeyTitle')}</h2>
	<p class="journey-lead">{t('home.journeyLead')}</p>
	<div class="lanes">
		<div class="lane lane-seller">
			<div class="lane-head">
				<span class="lane-avatar"><Icon name="person" /></span>
				<span class="lane-name">{t('home.journeySeller')}</span>
			</div>
			<ol class="lane-steps">
				{#each sellerSteps as step (step.n)}
					<li class="lane-step">
						<span class="lane-num" aria-hidden="true">{step.n}</span>
						<span class="lane-icon"><Icon name={step.icon} /></span>
						<span class="lane-label">{step.label}</span>
					</li>
				{/each}
			</ol>
		</div>
		<div class="lane lane-buyer">
			<div class="lane-head">
				<span class="lane-avatar"><Icon name="person" /></span>
				<span class="lane-name">{t('home.journeyBuyer')}</span>
			</div>
			<ol class="lane-steps">
				{#each buyerSteps as step (step.n)}
					<li class="lane-step">
						<span class="lane-num" aria-hidden="true">{step.n}</span>
						<span class="lane-icon"><Icon name={step.icon} /></span>
						<span class="lane-label">{step.label}</span>
					</li>
				{/each}
			</ol>
		</div>
	</div>
</section>

<section class="modes">
	{#each modes as mode (mode.name)}
		<div class="mode">
			<h2>{mode.name}</h2>
			<p>{mode.desc}</p>
		</div>
	{/each}
</section>

<section class="contact" aria-labelledby="contact-title">
	<h2 id="contact-title">{t('home.contactTitle')}</h2>
	<p class="contact-lead">{t('home.contactLead')}</p>
	<a class="contact-cta" href={link(resolve('/contact'))}>{t('home.contactCta')}</a>
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

	/* How it works: four friendly icon cards, linked by an arrow. */
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

	.how-steps {
		list-style: none;
		margin: 1.75rem 0 0;
		padding: 0;
		display: grid;
		gap: 1rem;
	}

	@media (min-width: 640px) {
		.how-steps {
			grid-template-columns: repeat(4, minmax(0, 1fr));
		}
	}

	.how-step {
		position: relative;
		border-radius: 0.85rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem 1rem;
	}

	/* arrow between cards: points down on mobile, right between columns on desktop */
	.how-step:not(:last-child)::after {
		content: '';
		position: absolute;
		left: 50%;
		bottom: -0.85rem;
		transform: translateX(-50%);
		width: 0;
		height: 0;
		border-left: 6px solid transparent;
		border-right: 6px solid transparent;
		border-top: 7px solid var(--slate-300);
	}

	@media (min-width: 640px) {
		.how-step:not(:last-child)::after {
			left: auto;
			bottom: auto;
			right: -0.8rem;
			top: 50%;
			transform: translateY(-50%);
			border-top: 6px solid transparent;
			border-bottom: 6px solid transparent;
			border-left: 7px solid var(--slate-300);
			border-right: 0;
		}
	}

	.how-icon {
		display: grid;
		place-items: center;
		width: 3rem;
		height: 3rem;
		margin: 0 auto;
		border-radius: 9999px;
		background: var(--emerald-50);
		color: var(--emerald-600);
		font-size: 1.5rem;
	}

	.how-step h3 {
		margin-top: 0.85rem;
		font-size: 1rem;
		font-weight: 600;
	}

	.how-step p {
		margin-top: 0.35rem;
		font-size: 0.85rem;
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

	/* Buyer and seller journey: two people, each heading a column of their own
	   steps; the numbers 1-6 run between the columns so it reads as a back-and-forth. */
	.journey {
		padding-block: 2rem;
		text-align: center;
	}

	.journey h2 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
		letter-spacing: -0.015em;
	}

	.journey-lead {
		margin: 0.5rem auto 0;
		max-width: 34rem;
		font-size: 0.95rem;
		line-height: 1.5rem;
		color: var(--slate-600);
	}

	.lanes {
		margin: 2rem auto 0;
		max-width: 22rem;
		display: grid;
		gap: 1.75rem;
		text-align: left;
	}

	.lane-head {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.6rem;
		margin-bottom: 1rem;
	}

	.lane-avatar {
		display: grid;
		place-items: center;
		width: 2.75rem;
		height: 2.75rem;
		border-radius: 9999px;
		font-size: 1.6rem;
	}

	.lane-seller .lane-avatar {
		background: var(--emerald-100);
		color: var(--emerald-700);
	}

	.lane-buyer .lane-avatar {
		background: var(--sky-100);
		color: var(--sky-700);
	}

	.lane-name {
		font-size: 1.05rem;
		font-weight: 700;
	}

	.lane-seller .lane-name {
		color: var(--emerald-800);
	}

	.lane-buyer .lane-name {
		color: var(--sky-800);
	}

	.lane-steps {
		position: relative;
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		gap: 0.85rem;
	}

	/* connector running down each person's column to their steps */
	.lane-steps::before {
		content: '';
		position: absolute;
		left: 1.1rem;
		top: 0.5rem;
		bottom: 0.5rem;
		width: 2px;
		background: var(--slate-200);
	}

	.lane-step {
		position: relative;
		display: flex;
		align-items: center;
		gap: 0.7rem;
		min-height: 2.25rem;
		padding-left: 2.9rem;
	}

	.lane-num {
		position: absolute;
		left: 0.35rem;
		top: 50%;
		transform: translateY(-50%);
		width: 1.5rem;
		height: 1.5rem;
		border-radius: 9999px;
		display: grid;
		place-items: center;
		font-size: 0.8rem;
		font-weight: 700;
		color: white;
		background: var(--emerald-600);
	}

	.lane-buyer .lane-num {
		background: var(--sky-600);
	}

	.lane-icon {
		flex: none;
		display: grid;
		place-items: center;
		font-size: 1.4rem;
		color: var(--emerald-700);
	}

	.lane-buyer .lane-icon {
		color: var(--sky-700);
	}

	.lane-label {
		font-size: 0.95rem;
		font-weight: 600;
		color: var(--slate-900);
	}

	@media (min-width: 720px) {
		.lanes {
			max-width: 44rem;
			grid-template-columns: 1fr 1fr;
			gap: 2.5rem;
			align-items: start;
		}

		/* offset the buyer's steps half a row so the numbers zig-zag 1-2-3-4-5-6 */
		.lane-buyer .lane-steps {
			margin-top: 1.7rem;
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

	.contact {
		margin-top: 1rem;
		border-radius: 1rem;
		border: 1px solid var(--emerald-100);
		background: var(--emerald-50);
		padding: 2.5rem 1.5rem;
		text-align: center;
	}

	.contact h2 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
		letter-spacing: -0.015em;
	}

	.contact-lead {
		margin: 0.5rem auto 0;
		max-width: 32rem;
		font-size: 0.95rem;
		line-height: 1.5rem;
		color: var(--slate-600);
	}

	.contact-cta {
		margin-top: 1.25rem;
		display: inline-block;
		border-radius: 0.5rem;
		background: var(--emerald-600);
		padding: 0.6rem 1.4rem;
		font-size: 0.9rem;
		font-weight: 600;
		color: white;
	}

	.contact-cta:hover {
		background: var(--emerald-700);
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
</style>
