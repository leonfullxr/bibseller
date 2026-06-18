<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { policyView } from '$lib/policy';

	let { policy, officialUrl = null }: { policy: TransferPolicy; officialUrl?: string | null } =
		$props();

	// The buy path exists ONLY for platform_sale, and even there it ships with M6
	// (honest disabled stub). Messaging is live (M5) via the contact composer on
	// the listing page, so there is no chat stub here.
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
</div>

<style>
	.cta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
	}

	.buy,
	.official {
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
</style>
