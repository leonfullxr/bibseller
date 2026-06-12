<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';

	let { policy, officialUrl = null }: { policy: TransferPolicy; officialUrl?: string | null } =
		$props();

	// The buy path exists ONLY for platform_sale — and even there it ships
	// with M6. Chat ships with M5. Until then: honest disabled stubs.
</script>

<div class="flex flex-wrap gap-3">
	{#if policy === 'platform_sale'}
		<button
			type="button"
			disabled
			class="cursor-not-allowed rounded-md bg-emerald-600 px-4 py-2 font-semibold text-white opacity-50"
			title="Secure checkout arrives with payments (M6)"
		>
			Buy securely — coming soon
		</button>
	{:else if policy === 'official_only' && officialUrl}
		<a
			href={officialUrl}
			rel="external nofollow noopener"
			target="_blank"
			class="rounded-md bg-sky-600 px-4 py-2 font-semibold text-white hover:bg-sky-700"
		>
			Official transfer process ↗
		</a>
	{/if}
	<button
		type="button"
		disabled
		class="cursor-not-allowed rounded-md border border-slate-300 px-4 py-2 font-semibold text-slate-500 opacity-60"
		title="Chat arrives with M5"
	>
		Message seller — coming soon
	</button>
</div>
