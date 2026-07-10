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

	// Ordered 1-6, alternating seller/buyer; the layout reads the `who` to place
	// each step on its side of the centre flow line.
	const journey = $derived([
		{ n: 1, who: 'seller', icon: 'list', label: t('home.j1Title') },
		{ n: 2, who: 'buyer', icon: 'search', label: t('home.j2Title') },
		{ n: 3, who: 'seller', icon: 'chat', label: t('home.j3Title') },
		{ n: 4, who: 'buyer', icon: 'transfer', label: t('home.j4Title') },
		{ n: 5, who: 'seller', icon: 'handover', label: t('home.j5Title') },
		{ n: 6, who: 'buyer', icon: 'medal', label: t('home.j6Title') }
	]);

	// Modes carry their policy tone (semantic colors, not brand).
	const modes = $derived([
		{ tone: 'allowed', name: t('home.modePlatformSaleName'), desc: t('home.modePlatformSaleDesc') },
		{ tone: 'official', name: t('home.modeOfficialName'), desc: t('home.modeOfficialDesc') },
		{ tone: 'connect', name: t('home.modeConnectName'), desc: t('home.modeConnectDesc') }
	]);
</script>

<svelte:head>
	<title>{t('home.title')}</title>
	<meta name="description" content={t('home.metaDescription')} />
</svelte:head>

<!-- The opening: a dawn-gradient field (races start at first light) with the
     headline and search left-aligned inside it. -->
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
		<h2 class="upcoming-title">{t('home.upcoming')}</h2>
		<div class="rail">
			{#each data.upcoming as race (race.id)}
				<div class="rail-item"><RaceCard {race} /></div>
			{/each}
		</div>
		<a class="browse-btn btn btn-outline" href={link(resolve('/races'))}
			>{t('home.browseAllRaces')}</a
		>
	</section>
{/if}

<!-- The four steps: a real sequence, numbered cards. -->
<section class="how" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="steps">
		{#each steps as step, i (step.icon)}
			<li class="step">
				<span class="step-n" aria-hidden="true">{i + 1}</span>
				<h3>{step.title}</h3>
				<p>{step.desc}</p>
			</li>
		{/each}
	</ol>
	<p class="how-note">{t('home.howNote')}</p>
</section>

<!-- The six-step handover as paired columns: the seller's moves on the
     left, the buyer's on the right, read in numbered zigzag order. -->
<section class="journey" aria-labelledby="journey-title">
	<h2 id="journey-title">{t('home.journeyTitle')}</h2>
	<p class="journey-lead">{t('home.journeyLead')}</p>
	<div class="duet">
		<span class="col-head">{t('home.journeySeller')}</span>
		<span class="col-head">{t('home.journeyBuyer')}</span>
		{#each journey as step (step.n)}
			<div class="move {step.who}">
				<span class="move-n" aria-hidden="true">{step.n}</span>
				<span class="move-icon" aria-hidden="true"><Icon name={step.icon} /></span>
				<span class="move-label">{step.label}</span>
				<span class="move-who"
					>{step.who === 'seller' ? t('home.journeySeller') : t('home.journeyBuyer')}</span
				>
			</div>
		{/each}
	</div>
</section>

<section class="modes">
	{#each modes as mode (mode.name)}
		<div class="mode {mode.tone}">
			<h2>{mode.name}</h2>
			<p>{mode.desc}</p>
		</div>
	{/each}
</section>

<section class="contact" aria-labelledby="contact-title">
	<h2 id="contact-title">{t('home.contactTitle')}</h2>
	<p class="contact-lead">{t('home.contactLead')}</p>
	<a class="contact-cta btn btn-primary" href={link(resolve('/contact'))}>{t('home.contactCta')}</a>
</section>

<p class="construction">
	{t('home.underConstruction')}
	<a href="https://github.com/leonfullxr/bibseller/issues/13" rel="external">{t('home.roadmap')}</a
	>.
</p>

<style>
	/* Dawn field: the one gradient moment on the site. Violet first light over
	   porcelain, everything inside left-aligned. */
	.hero {
		border: 1px solid var(--brand-100);
		border-radius: 1.25rem;
		background:
			radial-gradient(52rem 30rem at 92% -20%, var(--brand-300) 0%, transparent 55%),
			linear-gradient(160deg, var(--brand-100) 0%, var(--brand-50) 40%, white 100%);
		padding: 3.5rem 1.5rem 3rem;
	}

	@media (min-width: 640px) {
		.hero {
			padding: 4.5rem 3.5rem 4rem;
		}
	}

	.hero h1 {
		max-width: 46rem;
		font-size: 2.75rem;
		line-height: 1.02;
		font-weight: 800;
	}

	@media (min-width: 640px) {
		.hero h1 {
			font-size: 4rem;
		}
	}

	.hero h1 span {
		color: var(--brand-600);
	}

	.tagline {
		margin-top: 1.25rem;
		max-width: 38rem;
		font-size: 1.0625rem;
		line-height: 1.75rem;
		color: var(--ink-2);
	}

	.search {
		margin-top: 2rem;
		display: flex;
		max-width: 32rem;
		align-items: stretch;
		gap: 0.5rem;
	}

	.search input {
		width: 100%;
		border-radius: 9999px;
		border: 1px solid var(--slate-200);
		background: white;
		color: var(--ink);
		padding: 0.6875rem 1.25rem;
		font-size: 1rem;
		line-height: 1.5rem;
		box-shadow: var(--shadow-hard-sm);
	}

	/* Border swap for any focus; the global :focus-visible ring still applies
	   for keyboard users. */
	.search input:focus {
		border-color: var(--brand-600);
	}

	.search button {
		border-radius: 9999px;
		background: var(--brand-600);
		padding: 0.6875rem 1.625rem;
		font-size: 0.9375rem;
		line-height: 1.5rem;
		font-weight: 700;
		white-space: nowrap;
		color: white;
		transition: background-color 0.15s;
	}

	.search button:hover {
		background: var(--brand-700);
	}

	.browse-all {
		margin-top: 1.125rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: var(--brand-700);
		text-decoration: underline;
	}

	.browse-all:hover {
		color: var(--brand-800);
	}

	.api-status {
		margin-top: 1.5rem;
		display: flex;
		align-items: center;
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
		padding-block: 3rem 1.5rem;
	}

	/* Section headings: chunky display grotesk, left-aligned. */
	.upcoming-title,
	.how h2,
	.journey h2,
	.contact h2 {
		font-size: 1.75rem;
		line-height: 2.25rem;
		font-weight: 700;
	}

	/* A hand-scrolled rail of upcoming races; cards snap into place. */
	.rail {
		margin-top: 1.25rem;
		display: flex;
		gap: 1rem;
		overflow-x: auto;
		scroll-snap-type: x proximity;
		padding: 0.25rem 0.25rem 1rem;
	}

	.rail-item {
		flex: 0 0 17rem;
		max-width: 80vw;
		scroll-snap-align: start;
	}

	.browse-btn {
		margin-top: 0.75rem;
	}

	.how {
		padding-block: 2rem;
	}

	/* The four steps as a numbered card row. */
	.steps {
		list-style: none;
		margin: 1.5rem 0 0;
		padding: 0;
		display: grid;
		gap: 1rem;
	}

	@media (min-width: 720px) {
		.steps {
			grid-template-columns: repeat(4, minmax(0, 1fr));
		}
	}

	.step {
		border: 1px solid var(--slate-200);
		border-radius: var(--radius-lg);
		background: white;
		box-shadow: var(--shadow-hard-sm);
		padding: 1.375rem 1.25rem;
	}

	.step-n {
		font-family: var(--font-display);
		font-size: 2rem;
		line-height: 2.25rem;
		font-weight: 800;
		color: var(--brand-600);
	}

	.step h3 {
		margin-top: 0.5rem;
		font-size: 1.125rem;
		line-height: 1.5rem;
		font-weight: 700;
	}

	.step p {
		margin-top: 0.375rem;
		font-size: 0.875rem;
		line-height: 1.375rem;
		color: var(--slate-600);
	}

	.how-note {
		margin-top: 1.25rem;
		max-width: 36rem;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	.journey {
		padding-block: 2rem;
	}

	.journey-lead {
		margin-top: 0.5rem;
		max-width: 36rem;
		font-size: 0.9375rem;
		line-height: 1.5rem;
		color: var(--slate-600);
	}

	/* The duet: seller column left, buyer column right; auto-placement puts
	   odd steps (seller) left and even (buyer) right, so the numbered zigzag
	   reads spatially. Mobile stacks 1-6 with the number carrying order. */
	.duet {
		margin-top: 1.75rem;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.875rem 1.25rem;
	}

	.col-head {
		font-size: 0.75rem;
		line-height: 1.25rem;
		font-weight: 800;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		padding-bottom: 0.25rem;
	}

	.col-head:first-child {
		color: var(--brand-700);
	}

	.col-head:nth-child(2) {
		color: var(--sky-700);
	}

	.move {
		display: flex;
		align-items: center;
		gap: 0.875rem;
		border: 1px solid var(--slate-200);
		border-left-width: 3px;
		border-radius: var(--radius);
		background: white;
		box-shadow: var(--shadow-hard-sm);
		padding: 1rem 1.125rem;
	}

	.move.seller {
		border-left-color: var(--brand-600);
	}

	.move.buyer {
		border-left-color: var(--sky-600);
	}

	.move-n {
		font-family: var(--font-display);
		font-size: 1.5rem;
		line-height: 1.75rem;
		font-weight: 800;
		color: var(--brand-600);
		min-width: 1.25ch;
	}

	.move.buyer .move-n {
		color: var(--sky-700);
	}

	.move-icon {
		display: grid;
		place-items: center;
		width: 2.25rem;
		height: 2.25rem;
		flex-shrink: 0;
		border-radius: 9999px;
		background: var(--paper-2);
		font-size: 1.125rem;
		color: var(--slate-600);
	}

	.move-label {
		min-width: 0;
		font-size: 1rem;
		line-height: 1.375rem;
		font-weight: 600;
	}

	/* The columns already name the lane on wide screens. */
	.move-who {
		display: none;
		margin-left: auto;
		font-size: 0.6875rem;
		line-height: 1.25rem;
		font-weight: 800;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--slate-500);
	}

	@media (max-width: 39.9375rem) {
		.duet {
			grid-template-columns: 1fr;
		}

		.col-head {
			display: none;
		}

		/* Single column: name the lane on each card instead. */
		.move-who {
			display: inline;
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

	/* Each mode wears its policy tone (semantic, never brand). */
	.mode {
		border: 1px solid var(--slate-200);
		border-top-width: 3px;
		border-radius: var(--radius-lg);
		background: white;
		padding: 1.25rem;
		box-shadow: var(--shadow-hard-sm);
	}

	.mode.allowed {
		border-top-color: var(--emerald-600);
	}

	.mode.official {
		border-top-color: var(--sky-600);
	}

	.mode.connect {
		border-top-color: var(--amber-500);
	}

	.mode h2 {
		font-size: 1.1875rem;
		font-weight: 700;
	}

	.mode p {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.375rem;
		color: var(--slate-600);
	}

	/* The closing invitation: a quiet violet-tinted panel. */
	.contact {
		margin-top: 1rem;
		border-radius: var(--radius-lg);
		border: 1px solid var(--brand-100);
		background: var(--brand-50);
		padding: 2.5rem 1.5rem;
		text-align: center;
	}

	.contact-lead {
		margin: 0.5rem auto 0;
		max-width: 32rem;
		font-size: 0.9375rem;
		line-height: 1.5rem;
		color: var(--ink-2);
	}

	.contact-cta {
		margin-top: 1.25rem;
	}

	.construction {
		padding-block: 1.5rem 0.5rem;
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
