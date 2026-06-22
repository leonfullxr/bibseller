<script lang="ts">
	// Contact / imprint (M7.2, #54), doubling as the DSA single point of contact.
	// DRAFT copy pending counsel review (D10).
	import { getI18n } from '$lib/i18n';

	const { t } = getI18n();
</script>

<svelte:head>
	<title>{t('legal.contact.title')}</title>
	<meta name="description" content={t('legal.contact.metaDescription')} />
</svelte:head>

<h1>{t('legal.contact.heading')}</h1>
<p>{t('legal.contact.intro')}</p>

<h2>{t('legal.contact.operatorHeading')}</h2>
<p>{t('legal.contact.operatorBody')}</p>

<h2>{t('legal.contact.dsaHeading')}</h2>
<p>{t('legal.contact.dsaBody')}</p>

<h2>{t('legal.contact.reportHeading')}</h2>
<p>{t('legal.contact.reportBody')}</p>
