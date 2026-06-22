<script lang="ts">
	// Privacy Policy (M7.2, #54). DRAFT copy pending counsel review (D10).
	import { getI18n } from '$lib/i18n';

	const { t } = getI18n();
</script>

<svelte:head>
	<title>{t('legal.privacy.title')}</title>
	<meta name="description" content={t('legal.privacy.metaDescription')} />
</svelte:head>

<h1>{t('legal.privacy.heading')}</h1>
<p>{t('legal.privacy.intro')}</p>

<h2>{t('legal.privacy.dataHeading')}</h2>
<p>{t('legal.privacy.dataBody')}</p>

<h2>{t('legal.privacy.basisHeading')}</h2>
<p>{t('legal.privacy.basisBody')}</p>

<h2>{t('legal.privacy.retentionHeading')}</h2>
<p>{t('legal.privacy.retentionBody')}</p>

<h2>{t('legal.privacy.rightsHeading')}</h2>
<p>{t('legal.privacy.rightsBody')}</p>

<h2>{t('legal.privacy.cookiesHeading')}</h2>
<p>{t('legal.privacy.cookiesBody')}</p>

<h2>{t('legal.privacy.contactHeading')}</h2>
<p>{t('legal.privacy.contactBody')}</p>
