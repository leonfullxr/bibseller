<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';

	let { policy }: { policy: TransferPolicy } = $props();

	const styles: Record<TransferPolicy, { label: string; classes: string }> = {
		platform_sale: { label: 'Resale allowed', classes: 'bg-emerald-100 text-emerald-800' },
		official_only: { label: 'Official transfer', classes: 'bg-sky-100 text-sky-800' },
		connect_only: { label: 'Chat only', classes: 'bg-amber-100 text-amber-800' },
		unknown: { label: 'Policy unverified', classes: 'bg-slate-200 text-slate-700' }
	};
</script>

<span
	class="inline-block rounded-full px-2.5 py-0.5 text-xs font-semibold whitespace-nowrap {styles[
		policy
	].classes}"
>
	{styles[policy].label}
</span>
