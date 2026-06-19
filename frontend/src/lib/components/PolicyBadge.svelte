<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { getI18n } from '$lib/i18n';

	let { policy }: { policy: TransferPolicy } = $props();
	const { t } = getI18n();
</script>

<span class="badge {policy}">{t(`policy.label.${policy}`)}</span>

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
