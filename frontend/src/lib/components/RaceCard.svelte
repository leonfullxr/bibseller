<script lang="ts">
	import { resolve } from '$app/paths';
	import type { RaceSummary } from '$lib/api/types';
	import { formatDate } from '$lib/format';
	import { getI18n, sportLabel } from '$lib/i18n';
	import PolicyBadge from './PolicyBadge.svelte';

	let { race }: { race: RaceSummary } = $props();
	const { t, locale, link, plural } = getI18n();

	const bibs = $derived(plural('raceCard.bibs', race.active_listings));
</script>

<!-- Catalog plate: a hairline-framed card, the date as a small-caps
     bordeaux eyebrow, the race name in the journal serif. -->
<a href={link(resolve('/races/[slug]', { slug: race.slug }))} class="card">
	<div class="strip">
		<span class="date">{formatDate(race.event_date, locale)}</span>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<div class="body">
		<h3>{race.name}</h3>
		<p class="meta">{race.city}, {race.country}</p>
		<div class="tags">
			{#if race.distance}
				<span class="tag">{race.distance}</span>
			{/if}
			<span class="tag sport">{sportLabel(t, race.sport)}</span>
			<span class="count" class:active={race.active_listings > 0}>
				{bibs}
			</span>
		</div>
	</div>
</a>

<style>
	.card {
		display: block;
		border: 1px solid var(--slate-200);
		border-radius: 0.25rem;
		background: white;
		box-shadow: var(--shadow-hard-sm);
		overflow: hidden;
		transition:
			box-shadow 0.15s,
			border-color 0.15s;
	}

	.card:hover {
		border-color: var(--slate-400);
		box-shadow: var(--shadow-hard);
	}

	.strip {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		border-bottom: 1px solid var(--slate-200);
		padding: 0.5rem 1.125rem;
	}

	.date {
		font-size: 0.6875rem;
		line-height: 1.25rem;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: var(--brand-700);
	}

	.body {
		padding: 0.875rem 1.125rem 1.125rem;
	}

	h3 {
		font-family: var(--font-display);
		font-size: 1.375rem;
		line-height: 1.75rem;
		font-weight: 600;
		color: var(--ink);
	}

	.meta {
		margin-top: 0.25rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.tags {
		margin-top: 0.75rem;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}

	.tag {
		border-radius: 0.125rem;
		background: var(--slate-100);
		border: 1px solid var(--slate-200);
		padding: 0.125rem 0.375rem;
		font-weight: 500;
	}

	.sport {
		text-transform: capitalize;
	}

	.count {
		margin-left: auto;
		font-weight: 600;
	}

	.count.active {
		color: var(--emerald-700);
	}
</style>
