<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { getI18n } from '$lib/i18n';

	let { policy }: { policy: TransferPolicy } = $props();
	const { t } = getI18n();
</script>

<span class="pill {policy}">{t(`policy.label.${policy}`)}</span>

<style>
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
