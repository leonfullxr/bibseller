<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';

	let {
		policy,
		officialUrl = null,
		notes = null
	}: {
		policy: TransferPolicy;
		officialUrl?: string | null;
		notes?: string | null;
	} = $props();
</script>

{#if policy === 'platform_sale'}
	<div class="rounded-lg border border-emerald-200 bg-emerald-50 p-4 text-sm text-emerald-900">
		<p class="font-semibold">This race allows bib resale.</p>
		<p class="mt-1">
			Agree with the seller in chat, then pay securely through the platform — funds are held until
			the transfer is confirmed. Zero commission.
		</p>
		{#if notes}<p class="mt-2 text-emerald-800/80 italic">“{notes}”</p>{/if}
	</div>
{:else if policy === 'official_only'}
	<div class="rounded-lg border border-sky-200 bg-sky-50 p-4 text-sm text-sky-900">
		<p class="font-semibold">This race runs its own official name-change process.</p>
		<p class="mt-1">
			Find each other and agree on the details here — the transfer itself (and any official fee)
			goes through the race organizer. The platform never handles money for this race.
		</p>
		{#if officialUrl}
			<a
				href={officialUrl}
				rel="external nofollow noopener"
				target="_blank"
				class="mt-3 inline-block rounded-md bg-sky-600 px-3 py-1.5 font-semibold text-white hover:bg-sky-700"
			>
				Official transfer process ↗
			</a>
		{/if}
		{#if notes}<p class="mt-2 text-sky-800/80 italic">“{notes}”</p>{/if}
	</div>
{:else}
	<div class="rounded-lg border border-amber-300 bg-amber-50 p-4 text-sm text-amber-900">
		<p class="font-semibold">
			{policy === 'unknown'
				? 'Transfer policy not verified yet — treat this race as chat-only.'
				: 'This race restricts bib transfers.'}
		</p>
		<p class="mt-1">
			The platform only connects you: it handles no money here and takes no responsibility for any
			arrangement between you and the other party. The race's own rules apply — check them before
			agreeing to anything.
		</p>
		{#if notes}<p class="mt-2 text-amber-800/80 italic">“{notes}”</p>{/if}
	</div>
{/if}
