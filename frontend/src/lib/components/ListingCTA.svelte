<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { policyView } from '$lib/policy';

	let { policy, officialUrl = null }: { policy: TransferPolicy; officialUrl?: string | null } =
		$props();

	// The buy path exists ONLY for platform_sale - and even there it ships
	// with M6. Chat ships with M5. Until then: honest disabled stubs.
	const action = $derived(policyView[policy].primaryAction);
</script>

<div class="cta">
	{#if action === 'buy'}
		<button type="button" disabled class="buy" title="Secure checkout arrives with payments (M6)">
			Buy securely - coming soon
		</button>
	{:else if action === 'official' && officialUrl}
		<a href={officialUrl} rel="external nofollow noopener" target="_blank" class="official">
			Official transfer process
		</a>
	{/if}
	<button type="button" disabled class="chat" title="Chat arrives with M5">
		Message seller - coming soon
	</button>
</div>

<style>
	.cta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
	}

	.buy,
	.official,
	.chat {
		border-radius: 0.375rem;
		padding: 0.5rem 1rem;
		font-weight: 600;
	}

	.buy {
		background: var(--emerald-600);
		color: white;
		opacity: 0.5;
		cursor: not-allowed;
	}

	.official {
		display: inline-block;
		background: var(--sky-600);
		color: white;
	}

	.official:hover {
		background: var(--sky-700);
	}

	.chat {
		border: 1px solid var(--slate-300);
		color: var(--slate-500);
		opacity: 0.6;
		cursor: not-allowed;
	}
</style>
