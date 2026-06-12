<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';

	let { policy }: { policy: TransferPolicy } = $props();

	const labels: Record<TransferPolicy, string> = {
		platform_sale: 'Resale allowed',
		official_only: 'Official transfer',
		connect_only: 'Chat only',
		unknown: 'Policy unverified'
	};
</script>

<span class="badge {policy}">{labels[policy]}</span>

<style>
	.badge {
		display: inline-block;
		border-radius: 9999px;
		padding: 0.125rem 0.625rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 600;
		white-space: nowrap;
	}

	.platform_sale {
		background: var(--emerald-100);
		color: var(--emerald-800);
	}

	.official_only {
		background: var(--sky-100);
		color: var(--sky-800);
	}

	.connect_only {
		background: var(--amber-100);
		color: var(--amber-800);
	}

	.unknown {
		background: var(--slate-200);
		color: var(--slate-700);
	}
</style>
