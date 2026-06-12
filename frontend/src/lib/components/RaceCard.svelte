<script lang="ts">
	import { resolve } from '$app/paths';
	import type { RaceSummary } from '$lib/api/types';
	import { formatDate } from '$lib/format';
	import PolicyBadge from './PolicyBadge.svelte';

	let { race }: { race: RaceSummary } = $props();
</script>

<a
	href={resolve('/races/[slug]', { slug: race.slug })}
	class="block rounded-lg border border-slate-200 bg-white p-4 transition hover:border-slate-300 hover:shadow-sm"
>
	<div class="flex items-start justify-between gap-2">
		<h3 class="font-semibold text-slate-900">{race.name}</h3>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="mt-1 text-sm text-slate-600">
		{formatDate(race.event_date)} · {race.city}, {race.country}
	</p>
	<div class="mt-3 flex items-center gap-2 text-xs text-slate-500">
		{#if race.distance}
			<span class="rounded bg-slate-100 px-1.5 py-0.5 font-medium">{race.distance}</span>
		{/if}
		<span class="rounded bg-slate-100 px-1.5 py-0.5 font-medium capitalize">{race.sport}</span>
		<span class="ml-auto font-medium {race.active_listings > 0 ? 'text-emerald-700' : ''}">
			{race.active_listings}
			{race.active_listings === 1 ? 'bib' : 'bibs'} listed
		</span>
	</div>
</a>
