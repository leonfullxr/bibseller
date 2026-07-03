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

<section class="how" aria-labelledby="how-title">
	<h2 id="how-title">{t('home.howTitle')}</h2>
	<ol class="how-steps">
		{#each steps as step, i (step.icon)}
			<li class="how-step">
				<span class="how-icon"><Icon name={step.icon} /></span>
				<h3>{step.title}</h3>
				<p>{step.desc}</p>
			</li>
			{#if i < steps.length - 1}
				<li class="how-arrow" role="presentation" aria-hidden="true"><Icon name="arrow" /></li>
			{/if}
		{/each}
	</ol>
	<p class="how-note">{t('home.howNote')}</p>
</section>

<section class="journey" aria-labelledby="journey-title">
	<h2 id="journey-title">{t('home.journeyTitle')}</h2>
	<p class="journey-lead">{t('home.journeyLead')}</p>
	<div class="flow">
		<div class="flow-heads">
			<div class="flow-head lane-seller">
				<span class="lane-avatar"><Icon name="person" /></span>
				<span class="lane-name">{t('home.journeySeller')}</span>
			</div>
			<div class="flow-head lane-buyer">
				<span class="lane-avatar"><Icon name="person" /></span>
				<span class="lane-name">{t('home.journeyBuyer')}</span>
			</div>
		</div>
		<ol class="flow-steps">
			{#each journey as step (step.n)}
				<li class="flow-step flow-{step.who}">
					<div class="flow-card">
						<span class="flow-icon"><Icon name={step.icon} /></span>
						<span class="flow-label">{step.label}</span>
					</div>
				</li>
			{/each}
		</ol>
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
		text-align: center;
	}

	.upcoming-title {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
		letter-spacing: -0.015em;
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
		margin-right: 1rem;
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
		margin-top: 1.5rem;
		display: inline-block;
		border-radius: 0.5rem;
		background: var(--emerald-600);
		padding: 0.6rem 1.4rem;
		font-size: 0.9rem;
		font-weight: 600;
		color: white;
	}

	.browse-btn:hover {
		background: var(--emerald-700);
	}

	/* How it works: icon cards joined by arrows that scale with the viewport. */
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
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}

	.how-step {
		width: 100%;
		border-radius: 0.85rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem 1rem;
	}

	.how-arrow {
		flex: none;
		display: grid;
		place-items: center;
		color: var(--emerald-600);
		font-size: clamp(1.25rem, 6vw, 1.75rem);
		transform: rotate(90deg);
	}

	@media (min-width: 640px) {
		.how-steps {
			flex-direction: row;
			align-items: stretch;
			gap: clamp(0.5rem, 2.5vw, 1.75rem);
		}

		.how-step {
			flex: 1 1 0;
			width: auto;
		}

		.how-arrow {
			transform: none;
			font-size: clamp(1rem, 2.5vw, 1.75rem);
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
		color: var(--slate-500);
	}

	/* Buyer and seller journey: two people either side of a centre flow line;
	   each step connects to that line with a coloured node + connector. */
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

	.flow {
		margin: 2rem auto 0;
		max-width: 30rem;
	}

	.flow-heads {
		display: flex;
		justify-content: center;
		gap: 1.5rem;
		margin-bottom: 1.25rem;
	}

	.flow-head {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.lane-avatar {
		display: grid;
		place-items: center;
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 9999px;
		font-size: 1.5rem;
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
		font-weight: 700;
	}

	.lane-seller .lane-name {
		color: var(--emerald-800);
	}

	.lane-buyer .lane-name {
		color: var(--sky-800);
	}

	.flow-steps {
		position: relative;
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		gap: 1rem;
	}

	/* the centre flow line */
	.flow-steps::before {
		content: '';
		position: absolute;
		left: 1rem;
		top: 0.5rem;
		bottom: 0.5rem;
		width: 2px;
		background: var(--slate-200);
	}

	.flow-step {
		position: relative;
		display: flex;
		min-height: 2.5rem;
		padding-left: 2.85rem;
	}

	/* node sitting on the flow line */
	.flow-step::before {
		content: '';
		position: absolute;
		left: 1rem;
		top: 50%;
		transform: translate(-50%, -50%);
		width: 0.8rem;
		height: 0.8rem;
		border-radius: 9999px;
		border: 2px solid white;
		z-index: 1;
	}

	/* connector from the node to the card */
	.flow-step::after {
		content: '';
		position: absolute;
		left: 1rem;
		top: 50%;
		width: 1.85rem;
		height: 2px;
	}

	.flow-seller::before {
		background: var(--emerald-600);
	}

	.flow-seller::after {
		background: var(--emerald-600);
	}

	.flow-buyer::before {
		background: var(--sky-600);
	}

	.flow-buyer::after {
		background: var(--sky-600);
	}

	.flow-card {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		width: 100%;
		border-radius: 0.85rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 0.7rem 0.95rem;
	}

	.flow-seller .flow-card {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
	}

	.flow-buyer .flow-card {
		border-color: var(--sky-200);
		background: var(--sky-50);
	}

	.flow-icon {
		flex: none;
		display: grid;
		place-items: center;
		font-size: 1.4rem;
		color: var(--emerald-700);
	}

	.flow-buyer .flow-icon {
		color: var(--sky-700);
	}

	.flow-label {
		font-size: 0.95rem;
		font-weight: 600;
		color: var(--slate-900);
	}

	@media (min-width: 720px) {
		.flow {
			max-width: 46rem;
		}

		.flow-heads {
			justify-content: space-between;
			padding-inline: 3rem;
		}

		.flow-steps::before {
			left: 50%;
			transform: translateX(-50%);
		}

		.flow-step {
			padding-left: 0;
		}

		.flow-card {
			width: calc(50% - 2.85rem);
		}

		.flow-seller {
			justify-content: flex-start;
		}

		.flow-buyer {
			justify-content: flex-end;
		}

		.flow-seller .flow-card {
			flex-direction: row-reverse;
			text-align: right;
		}

		.flow-step::before {
			left: 50%;
		}

		.flow-seller::after {
			left: auto;
			right: 50%;
			width: 2.85rem;
		}

		.flow-buyer::after {
			left: 50%;
			width: 2.85rem;
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
		color: var(--slate-500);
	}

	.construction a {
		text-decoration: underline;
	}

	.construction a:hover {
		color: var(--slate-600);
	}
</style>
