<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { getI18n } from '$lib/i18n';
	import { policyView } from '$lib/policy';

	let { policy, officialUrl = null }: { policy: TransferPolicy; officialUrl?: string | null } =
		$props();

	const { t } = getI18n();

	// The buy path exists ONLY for platform_sale, and even there it ships with M6
	// (honest disabled stub). Messaging is live (M5) via the contact composer on
	// the listing page, so there is no chat stub here.
	const action = $derived(policyView[policy].primaryAction);
</script>

<div class="cta">
	{#if action === 'buy'}
		<button type="button" disabled class="btn btn-primary buy" title={t('listingCta.buyTitle')}>
			{t('listingCta.buy')}
		</button>
	{:else if action === 'official' && officialUrl}
		<a href={officialUrl} rel="external nofollow noopener" target="_blank" class="btn official">
			{t('policy.officialLink')}
		</a>
	{/if}
</div>

<style>
	.cta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
	}

	.buy:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* The official-transfer action keeps the sky "official" tone (PolicyBadge,
	   PolicyCallout) rather than the brand primary. */
	.official {
		background: var(--sky-600);
	}

	.official:hover {
		background: var(--sky-700);
	}
</style>
