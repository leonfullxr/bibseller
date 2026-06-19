<script lang="ts">
	import type { TransferPolicy } from '$lib/api/types';
	import { getI18n } from '$lib/i18n';
	import { policyView } from '$lib/policy';

	let {
		policy,
		officialUrl = null,
		notes = null
	}: {
		policy: TransferPolicy;
		officialUrl?: string | null;
		notes?: string | null;
	} = $props();

	const { t } = getI18n();
	const tone = $derived(policyView[policy].tone);
</script>

<div class="callout {tone}">
	<p class="title">{t(`policy.disclaimer.${policy}.title`)}</p>
	<p>{t(`policy.disclaimer.${policy}.body`)}</p>
	{#if tone === 'official' && officialUrl}
		<a href={officialUrl} rel="external nofollow noopener" target="_blank" class="official-link">
			{t('policy.officialLink')}
		</a>
	{/if}
	{#if notes}<p class="notes">“{notes}”</p>{/if}
</div>

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
