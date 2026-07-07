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

<!-- The card is a bib tag: punched holes on the top band, the event date as
     the "number", hard poster shadow that flattens on hover. -->
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
		border: 2px solid var(--ink);
		border-radius: 0.5rem;
		background: white;
		box-shadow: var(--shadow-hard);
		overflow: hidden;
		transition:
			translate 0.1s,
			box-shadow 0.1s;
	}

	.card:hover {
		translate: 2px 2px;
		box-shadow: 2px 2px 0 var(--ink);
	}

	/* Top band with the zip-tie holes. */
	.strip {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		background: var(--paper-2);
		border-bottom: 2px solid var(--ink);
		padding: 0.375rem 1.75rem;
	}

	.strip::before,
	.strip::after {
		content: '';
		position: absolute;
		top: 50%;
		translate: 0 -50%;
		width: 0.625rem;
		height: 0.625rem;
		border-radius: 9999px;
		background: var(--paper);
		border: 2px solid var(--ink);
	}

	.strip::before {
		left: 0.5rem;
	}

	.strip::after {
		right: 0.5rem;
	}

	.date {
		font-family: var(--font-display);
		font-size: 1rem;
		line-height: 1.5rem;
		font-weight: 800;
		letter-spacing: 0.02em;
		text-transform: uppercase;
	}

	.body {
		padding: 0.875rem 1rem 1rem;
	}

	h3 {
		font-size: 1.375rem;
		line-height: 1.625rem;
		font-weight: 800;
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
		border-radius: 0.25rem;
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
