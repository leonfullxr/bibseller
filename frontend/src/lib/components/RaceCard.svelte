<script lang="ts">
	import { resolve } from '$app/paths';
	import type { RaceSummary } from '$lib/api/types';
	import { formatDate } from '$lib/format';
	import PolicyBadge from './PolicyBadge.svelte';

	let { race }: { race: RaceSummary } = $props();
</script>

<a href={resolve('/races/[slug]', { slug: race.slug })} class="card">
	<div class="top">
		<h3>{race.name}</h3>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="meta">
		{formatDate(race.event_date)} - {race.city}, {race.country}
	</p>
	<div class="tags">
		{#if race.distance}
			<span class="tag">{race.distance}</span>
		{/if}
		<span class="tag sport">{race.sport}</span>
		<span class="count" class:active={race.active_listings > 0}>
			{race.active_listings}
			{race.active_listings === 1 ? 'bib' : 'bibs'} listed
		</span>
	</div>
</a>

<style>
	.card {
		display: block;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 1rem;
		transition:
			border-color 150ms,
			box-shadow 150ms;
	}

	.card:hover {
		border-color: var(--slate-300);
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
	}

	.top {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 0.5rem;
	}

	h3 {
		font-weight: 600;
		color: var(--slate-900);
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
		padding: 0.125rem 0.375rem;
		font-weight: 500;
	}

	.sport {
		text-transform: capitalize;
	}

	.count {
		margin-left: auto;
		font-weight: 500;
	}

	.count.active {
		color: var(--emerald-700);
	}
</style>
