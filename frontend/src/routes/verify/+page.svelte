<script lang="ts">
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, link } = getI18n();
</script>

<svelte:head><title>{t('verify.title')}</title></svelte:head>

<section class="verify">
	{#if data.status === 'ok'}
		<h1>{t('verify.okHeading')}</h1>
		<p>{t('verify.okBody')}</p>
		<a href={link(resolve('/'))}>{t('verify.continue')}</a>
	{:else if data.status === 'invalid'}
		<h1>{t('verify.invalidHeading')}</h1>
		<p>{t('verify.invalidBody')}</p>
		<a href={link(resolve('/login'))}>{t('verify.signIn')}</a>
	{:else if data.status === 'missing'}
		<h1>{t('verify.missingHeading')}</h1>
		<p>{t('verify.missingBody')}</p>
		<a href={link(resolve('/'))}>{t('verify.home')}</a>
	{:else}
		<h1>{t('verify.errorHeading')}</h1>
		<p>{t('verify.errorBody')}</p>
		<a href={link(resolve('/'))}>{t('verify.home')}</a>
	{/if}
</section>

<style>
	.verify {
		max-width: 28rem;
		margin-inline: auto;
		text-align: center;
	}

	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	p {
		margin-top: 0.5rem;
		color: var(--slate-600);
	}

	a {
		display: inline-block;
		margin-top: 1.5rem;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	a:hover {
		background: var(--slate-700);
	}
</style>
