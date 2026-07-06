<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { form }: PageProps = $props();
	const { t, link } = getI18n();
	const { busy, submit } = pendingForm();
</script>

<svelte:head>
	<title>{t('login.title')}</title>
</svelte:head>

<h1>{t('nav.login')}</h1>

<form method="POST" use:enhance={submit}>
	<label for="email">{t('auth.email')}</label>
	<input
		id="email"
		name="email"
		type="email"
		required
		autocomplete="email"
		class="field"
		value={form?.email ?? ''}
	/>

	<label for="password">{t('auth.password')}</label>
	<input
		id="password"
		name="password"
		type="password"
		required
		autocomplete="current-password"
		class="field"
	/>

	{#if form?.error}
		<p class="alert" role="alert">{form.error}</p>
	{/if}

	<button type="submit" class="btn btn-primary" disabled={busy.value}>{t('nav.login')}</button>
</form>

<p class="alt"><a href={link(resolve('/(auth)/forgot'))}>{t('login.forgot')}</a></p>
<p class="alt">
	{t('login.newHere')} <a href={link(resolve('/(auth)/register'))}>{t('login.createAccount')}</a>
</p>
