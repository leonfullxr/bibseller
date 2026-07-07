<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { getI18n } from '$lib/i18n';

	let { policy }: { policy: TransferPolicy } = $props();
	const { t } = getI18n();
</script>

<span class="pill {policy}">{t(`policy.label.${policy}`)}</span>

<style>
	/* Policy tones are semantic (green = resale allowed) - not brand. The 1px
	   tone border keeps the pill crisp on the warm paper surfaces. */
	.platform_sale {
		background: var(--emerald-100);
		border: 1px solid var(--emerald-300);
		color: var(--emerald-800);
	}

	.official_only {
		background: var(--sky-100);
		border: 1px solid var(--sky-200);
		color: var(--sky-800);
	}

	.connect_only {
		background: var(--amber-100);
		border: 1px solid var(--amber-300);
		color: var(--amber-800);
	}

	.unknown {
		background: var(--slate-200);
		border: 1px solid var(--slate-300);
		color: var(--slate-700);
	}
</style>
