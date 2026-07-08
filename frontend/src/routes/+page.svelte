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
		<h2 class="upcoming-title">{t('home.upcoming')}</h2>
		<div class="marquee">
			<div class="marquee-track" style="--marquee-duration: {data.upcoming.length * 6}s">
				{#each data.upcoming as race (race.id)}
					<div class="marquee-item"><RaceCard {race} /></div>
				{/each}
				{#each data.upcoming as race (race.id + '-dup')}
					<div class="marquee-item dup" inert aria-hidden="true"><RaceCard {race} /></div>
				{/each}
			</div>
		</div>
		<a class="browse-btn" href={link(resolve('/races'))}>{t('home.browseAllRaces')}</a>
	</section>
{/if}

<!-- The four steps as checkpoints along a race course: a track line with
     numbered markers, checkered flag at the finish. -->
<section class="how" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="course">
		{#each steps as step, i (step.icon)}
			<li class="checkpoint">
				<span class="marker" class:finish={i === steps.length - 1} aria-hidden="true"
					>{#if i < steps.length - 1}{i + 1}{/if}</span
				>
				<h3>{step.title}</h3>
				<p>{step.desc}</p>
			</li>
		{/each}
	</ol>
	<p class="how-note">{t('home.howNote')}</p>
</section>

<!-- The six-step handover as a race timing board: ink leaderboard rows,
     orange positions, a lane tag naming whose move it is. -->
<section class="journey" aria-labelledby="journey-title">
	<h2 id="journey-title">{t('home.journeyTitle')}</h2>
	<p class="journey-lead">{t('home.journeyLead')}</p>
	<ol class="board">
		{#each journey as step (step.n)}
			<li class="row {step.who}">
				<span class="row-n" aria-hidden="true">{String(step.n).padStart(2, '0')}</span>
				<span class="row-icon" aria-hidden="true"><Icon name={step.icon} /></span>
				<span class="row-label">{step.label}</span>
				<span class="row-who {step.who}"
					>{step.who === 'seller' ? t('home.journeySeller') : t('home.journeyBuyer')}</span
				>
			</li>
		{/each}
	</ol>
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
	/* Full-bleed ink hero with a finish-line diagonal at the bottom - the page
	   opens like a race poster, not a web form. */
	.hero {
		margin-top: -2.5rem; /* cancel main's padding so the ink meets the header */
		margin-inline: calc(50% - 50vw);
		padding: 4rem 1rem 6rem;
		background: var(--ink);
		color: var(--paper);
		text-align: center;
		clip-path: polygon(0 0, 100% 0, 100% calc(100% - 2.75rem), 0 100%);
	}

	.hero h1 {
		margin-inline: auto;
		max-width: 56rem;
		font-size: 3.25rem;
		line-height: 0.95;
		font-weight: 800;
		color: white;
	}

	@media (min-width: 640px) {
		.hero h1 {
			font-size: 5rem;
		}
	}

	.hero h1 span {
		color: var(--brand-500);
	}

	.tagline {
		margin: 1.25rem auto 0;
		max-width: 42rem;
		font-size: 1.125rem;
		line-height: 1.75rem;
		color: var(--slate-300);
	}

	.search {
		margin: 2.25rem auto 0;
		display: flex;
		max-width: 30rem;
		align-items: stretch;
		gap: 0.5rem;
	}

	.search input {
		width: 100%;
		border-radius: 0.25rem;
		border: 2px solid transparent;
		background: white;
		color: var(--ink);
		padding: 0.75rem 1rem;
		font-size: 1rem;
		line-height: 1.5rem;
	}

	/* Border swap for any focus; the global :focus-visible ring still applies
	   for keyboard users (visible against the ink hero). */
	.search input:focus {
		border-color: var(--brand-500);
	}

	.search button {
		border-radius: 0.25rem;
		background: var(--brand-700);
		padding: 0.75rem 1.5rem;
		font-family: var(--font-display);
		font-size: 1.125rem;
		line-height: 1.5rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		white-space: nowrap;
		color: white;
		transition: background-color 0.15s;
	}

	.search button:hover {
		background: var(--brand-800);
	}

	.browse-all {
		margin-top: 1rem;
		display: inline-block;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--brand-300);
		text-decoration: underline;
	}

	.browse-all:hover {
		color: var(--brand-100);
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
		padding-block: 2.5rem 2rem;
		text-align: center;
	}

	/* Section headings share one device: condensed 800 with a short orange
	   underline bar - start-line tape under the sign. */
	.upcoming-title,
	.how h2,
	.journey h2,
	.contact h2 {
		font-size: 1.75rem;
		line-height: 2.25rem;
		font-weight: 800;
	}

	.upcoming-title::after,
	.how h2::after,
	.journey h2::after {
		content: '';
		display: block;
		margin: 0.5rem auto 0;
		width: 3rem;
		height: 0.25rem;
		background: var(--brand-500);
	}

	/* Auto-scrolling "roulette" of a few races. The list is rendered twice (the
	   copy is inert + aria-hidden) so the translateX(-50%) loop is seamless. */
	.marquee {
		margin-top: 1.5rem;
		overflow: hidden;
		-webkit-mask-image: linear-gradient(
			to right,
			transparent,
			#000 3rem,
			#000 calc(100% - 3rem),
			transparent
		);
		mask-image: linear-gradient(
			to right,
			transparent,
			#000 3rem,
			#000 calc(100% - 3rem),
			transparent
		);
	}

	.marquee-track {
		display: flex;
		align-items: stretch;
		width: max-content;
		animation: marquee var(--marquee-duration, 36s) linear infinite;
	}

	.marquee:hover .marquee-track,
	.marquee:focus-within .marquee-track {
		animation-play-state: paused;
	}

	.marquee-item {
		flex: 0 0 17rem;
		max-width: 80vw;
		margin-right: 1.25rem;
		/* room for the bib cards' hard shadows inside the overflow clip */
		padding: 2px 6px 6px 2px;
		text-align: left;
	}

	@keyframes marquee {
		from {
			transform: translateX(0);
		}
		to {
			transform: translateX(-50%);
		}
	}

	/* No auto-motion when the user prefers reduced motion; allow manual scroll. */
	@media (prefers-reduced-motion: reduce) {
		.marquee {
			overflow-x: auto;
		}

		.marquee-track {
			animation: none;
		}

		/* The duplicate half exists only to loop the animation seamlessly. With the
		   animation stopped it would just be inert, unclickable copies that a manual
		   scroll could reach (reads as broken links), so drop it from layout (#92). */
		.marquee-item.dup {
			display: none;
		}
	}

	.browse-btn {
		margin-top: 1.75rem;
		display: inline-block;
		border-radius: 0.25rem;
		background: var(--brand-700);
		padding: 0.625rem 1.5rem;
		font-family: var(--font-display);
		font-size: 1.0625rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: white;
		box-shadow: var(--shadow-hard-sm);
		transition:
			translate 0.1s,
			box-shadow 0.1s,
			background-color 0.15s;
	}

	.browse-btn:hover {
		background: var(--brand-800);
		translate: 1px 1px;
		box-shadow: 2px 2px 0 var(--ink);
	}

	/* How it works: icon cards joined by arrows that scale with the viewport. */
	.how {
		padding-block: 2rem;
		text-align: center;
	}

	/* The course: a track line through numbered checkpoint markers. */
	.course {
		list-style: none;
		margin: 2.25rem 0 0;
		padding: 0;
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 1.75rem;
		text-align: left;
	}

	.course::before {
		content: '';
		position: absolute;
		left: 1.25rem;
		top: 1rem;
		bottom: 1rem;
		width: 3px;
		background: var(--ink);
	}

	.checkpoint {
		position: relative;
		padding-left: 3.75rem;
	}

	.marker {
		position: absolute;
		left: 0;
		top: 0;
		display: grid;
		place-items: center;
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 9999px;
		background: var(--brand-500);
		border: 3px solid var(--ink);
		color: var(--ink);
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 800;
	}

	/* The last checkpoint is the finish: checkered flag, no number. */
	.marker.finish {
		border-radius: 0.25rem;
		background: repeating-conic-gradient(var(--ink) 0% 25%, white 0% 50%) 50% / 0.625rem 0.625rem;
	}

	.checkpoint h3 {
		font-size: 1.375rem;
		line-height: 1.625rem;
		font-weight: 800;
	}

	.checkpoint p {
		margin-top: 0.25rem;
		max-width: 34rem;
		font-size: 0.9rem;
		line-height: 1.4rem;
		color: var(--slate-600);
	}

	@media (min-width: 720px) {
		.course {
			flex-direction: row;
			gap: 1.5rem;
			text-align: center;
		}

		.course::before {
			left: 3rem;
			right: 3rem;
			top: 1.25rem;
			bottom: auto;
			width: auto;
			height: 3px;
		}

		.checkpoint {
			flex: 1 1 0;
			padding-left: 0;
			padding-top: 3.5rem;
		}

		.marker {
			left: 50%;
			translate: -50%;
		}
	}

	.how-note {
		margin: 1.5rem auto 0;
		max-width: 32rem;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	.journey {
		padding-block: 2rem;
		text-align: center;
	}

	.journey-lead {
		margin: 0.5rem auto 0;
		max-width: 34rem;
		font-size: 0.95rem;
		line-height: 1.5rem;
		color: var(--slate-600);
	}

	/* The timing board: an ink leaderboard. Each row is lane-coded - an
	   accent edge, a tinted wash, and a lane-colored position number - so the
	   seller/buyer alternation reads at a glance. */
	.board {
		list-style: none;
		margin: 2rem auto 0;
		padding: 0;
		max-width: 46rem;
		background: var(--ink);
		border-radius: 0.75rem;
		box-shadow: var(--shadow-hard);
		overflow: hidden;
		text-align: left;
	}

	.row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.5rem 1rem;
		padding: 1.125rem 1.5rem;
		color: white;
	}

	.row + .row {
		border-top: 1px solid var(--ink-2);
	}

	.row.seller {
		box-shadow: inset 0.25rem 0 0 var(--brand-500);
		background: rgb(249 115 22 / 0.07);
	}

	.row.buyer {
		box-shadow: inset 0.25rem 0 0 var(--sky-200);
		background: rgb(186 230 253 / 0.05);
	}

	.row-n {
		min-width: 2.5ch;
		font-family: var(--font-display);
		font-size: 1.75rem;
		line-height: 1.75rem;
		font-weight: 800;
	}

	.row.seller .row-n {
		color: var(--brand-500);
	}

	.row.buyer .row-n {
		color: var(--sky-200);
	}

	.row-icon {
		display: grid;
		place-items: center;
		width: 2.75rem;
		height: 2.75rem;
		flex-shrink: 0;
		border-radius: 9999px;
		background: var(--ink-2);
		font-size: 1.375rem;
		color: var(--slate-300);
	}

	.row-label {
		flex: 1;
		min-width: 0;
		font-size: 1.125rem;
		line-height: 1.5rem;
		font-weight: 600;
	}

	.row-who {
		flex-shrink: 0;
		border-radius: 9999px;
		padding: 0.1875rem 0.75rem;
		font-family: var(--font-display);
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.row-who.seller {
		background: var(--brand-500);
		color: var(--ink);
	}

	.row-who.buyer {
		background: var(--sky-200);
		color: var(--ink);
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
		border-top: 4px solid var(--brand-500);
		background: white;
		padding: 1.25rem;
	}

	.mode h2 {
		font-size: 1.25rem;
		font-weight: 700;
	}

	.mode p {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	/* The closing block answers the hero: an orange slab with an ink button. */
	.contact {
		margin-top: 1rem;
		border-radius: 0.75rem;
		background: var(--brand-700);
		padding: 2.5rem 1.5rem;
		text-align: center;
		color: white;
	}

	.contact h2 {
		color: white;
	}

	.contact-lead {
		margin: 0.5rem auto 0;
		max-width: 32rem;
		font-size: 0.95rem;
		line-height: 1.5rem;
		color: var(--brand-100);
	}

	.contact-cta {
		margin-top: 1.25rem;
		display: inline-block;
		border-radius: 0.25rem;
		background: var(--ink);
		padding: 0.625rem 1.5rem;
		font-family: var(--font-display);
		font-size: 1.0625rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: white;
		transition: background-color 0.15s;
	}

	.contact-cta:hover {
		background: var(--ink-2);
	}

	.construction {
		padding-block: 1rem;
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
