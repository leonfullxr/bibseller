<script lang="ts">
	// Terms of Service (M7.1, #53). DRAFT copy pending counsel review (D10).
	// The per-policy responsibility split iterates transferPolicies so all four
	// modes are always covered and stay consistent with the policy.label.* keys.
	import { getI18n } from '$lib/i18n';
	import { transferPolicies } from '$lib/policy';

	const { t } = getI18n();
</script>

<svelte:head>
	<title>{t('legal.terms.title')}</title>
	<meta name="description" content={t('legal.terms.metaDescription')} />
</svelte:head>

<h1>{t('legal.terms.heading')}</h1>
<p>{t('legal.terms.intro')}</p>

<h2>{t('legal.terms.roleHeading')}</h2>
<p>{t('legal.terms.roleBody')}</p>

<h2>{t('legal.terms.policyHeading')}</h2>
<p>{t('legal.terms.policyIntro')}</p>
<dl>
	{#each transferPolicies as policy (policy)}
		<dt>{t(`policy.label.${policy}`)}</dt>
		<dd>{t(`legal.terms.policy.${policy}`)}</dd>
	{/each}
</dl>

<h2>{t('legal.terms.userHeading')}</h2>
<p>{t('legal.terms.userBody')}</p>

<h2>{t('legal.terms.conductHeading')}</h2>
<p>{t('legal.terms.conductBody')}</p>

<h2>{t('legal.terms.liabilityHeading')}</h2>
<p>{t('legal.terms.liabilityBody')}</p>

<h2>{t('legal.terms.changesHeading')}</h2>
<p>{t('legal.terms.changesBody')}</p>
