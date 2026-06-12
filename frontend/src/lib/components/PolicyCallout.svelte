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
	<div class="callout sale">
		<p class="title">This race allows bib resale.</p>
		<p>
			Agree with the seller in chat, then pay securely through the platform — funds are held until
			the transfer is confirmed. Zero commission.
		</p>
		{#if notes}<p class="notes">“{notes}”</p>{/if}
	</div>
{:else if policy === 'official_only'}
	<div class="callout official">
		<p class="title">This race runs its own official name-change process.</p>
		<p>
			Find each other and agree on the details here — the transfer itself (and any official fee)
			goes through the race organizer. The platform never handles money for this race.
		</p>
		{#if officialUrl}
			<a href={officialUrl} rel="external nofollow noopener" target="_blank" class="official-link">
				Official transfer process ↗
			</a>
		{/if}
		{#if notes}<p class="notes">“{notes}”</p>{/if}
	</div>
{:else}
	<div class="callout restricted">
		<p class="title">
			{policy === 'unknown'
				? 'Transfer policy not verified yet — treat this race as chat-only.'
				: 'This race restricts bib transfers.'}
		</p>
		<p>
			The platform only connects you: it handles no money here and takes no responsibility for any
			arrangement between you and the other party. The race's own rules apply — check them before
			agreeing to anything.
		</p>
		{#if notes}<p class="notes">“{notes}”</p>{/if}
	</div>
{/if}

<style>
	.callout {
		border-radius: 0.5rem;
		border: 1px solid;
		padding: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.title {
		font-weight: 600;
	}

	.title + p {
		margin-top: 0.25rem;
	}

	.notes {
		margin-top: 0.5rem;
		font-style: italic;
	}

	.sale {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	.sale .notes {
		color: color-mix(in srgb, var(--emerald-800) 80%, transparent);
	}

	.official {
		border-color: var(--sky-200);
		background: var(--sky-50);
		color: var(--sky-900);
	}

	.official .notes {
		color: color-mix(in srgb, var(--sky-800) 80%, transparent);
	}

	.official-link {
		margin-top: 0.75rem;
		display: inline-block;
		border-radius: 0.375rem;
		background: var(--sky-600);
		padding: 0.375rem 0.75rem;
		font-weight: 600;
		color: white;
	}

	.official-link:hover {
		background: var(--sky-700);
	}

	.restricted {
		border-color: var(--amber-300);
		background: var(--amber-50);
		color: var(--amber-900);
	}

	.restricted .notes {
		color: color-mix(in srgb, var(--amber-800) 80%, transparent);
	}
</style>
