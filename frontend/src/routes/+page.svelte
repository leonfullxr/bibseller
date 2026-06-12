<script lang="ts">
	import { resolve } from '$app/paths';
	import RaceCard from '$lib/components/RaceCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modes = [
		{
			name: 'Platform sale',
			desc: 'The race allows resale: list, chat, and pay securely through the platform.'
		},
		{
			name: 'Official process',
			desc: 'The race runs its own name change: we connect you and link the official procedure.'
		},
		{
			name: 'Connect only',
			desc: 'Restricted or unverified races: we provide the chat, the rest stays between you two.'
		}
	];
</script>

<svelte:head>
	<title>Bibseller — race bibs find new runners</title>
	<meta
		name="description"
		content="Non-profit, EU-wide marketplace connecting runners who can't start with runners who missed registration."
	/>
</svelte:head>

<section class="hero">
	<h1>Race bibs find <span>new runners</span></h1>
	<p class="tagline">
		Injured? Plans changed? Missed registration? A zero-commission marketplace that connects sellers
		and buyers of race bibs — always within each race's own rules.
	</p>

	<form method="GET" action={resolve('/races')} class="search">
		<input type="search" name="q" placeholder="Search a race or city…" />
		<button type="submit">Search</button>
	</form>
	<a href={resolve('/races')} class="browse-all">or browse all races</a>

	{#if data.apiStatus !== 'ok'}
		<div class="api-status">
			<span class="dot"></span>
			<span class="api-msg">
				API unreachable — run <code>make dev</code>
			</span>
		</div>
	{/if}
</section>

{#if data.upcoming.length > 0}
	<section class="upcoming">
		<div class="upcoming-head">
			<h2>Upcoming races</h2>
			<a href={resolve('/races')}>See all</a>
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
	Under construction — follow the
	<a href="https://github.com/leonfullxr/bibseller/issues/13" rel="external">roadmap</a>.
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
</style>
