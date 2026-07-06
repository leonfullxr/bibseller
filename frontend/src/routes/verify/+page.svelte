<script lang="ts">
	import { onMount } from 'svelte';
	import { invalidate } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, link } = getI18n();

	// On success, refresh the layout's user so the "verify your email" banner
	// clears immediately. Targeted ('app:user') so this page's own load does not
	// re-run and re-POST the now-consumed token (which would flip it to "invalid").
	onMount(() => {
		if (data.status === 'ok') invalidate('app:user');
	});
</script>

<svelte:head><title>{t('verify.title')}</title></svelte:head>

<section class="verify">
	{#if data.status === 'ok'}
		<h1>{t('verify.okHeading')}</h1>
		<p>{t('verify.okBody')}</p>
		<a class="btn btn-primary" href={link(resolve('/'))}>{t('verify.continue')}</a>
	{:else if data.status === 'invalid'}
		<h1>{t('verify.invalidHeading')}</h1>
		<p>{t('verify.invalidBody')}</p>
		<a class="btn btn-primary" href={link(resolve('/(auth)/login'))}>{t('verify.signIn')}</a>
	{:else if data.status === 'missing'}
		<h1>{t('verify.missingHeading')}</h1>
		<p>{t('verify.missingBody')}</p>
		<a class="btn btn-primary" href={link(resolve('/'))}>{t('verify.home')}</a>
	{:else}
		<h1>{t('verify.errorHeading')}</h1>
		<p>{t('verify.errorBody')}</p>
		<a class="btn btn-primary" href={link(resolve('/'))}>{t('verify.home')}</a>
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
		margin-top: 1.5rem;
	}
</style>
