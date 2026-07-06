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
	<title>{t('forgot.title')}</title>
</svelte:head>

<h1>{t('forgot.heading')}</h1>

{#if form?.sent}
	<p class="alert ok" role="status">{t('forgot.sent')}</p>
	<p class="alt"><a href={link(resolve('/(auth)/login'))}>{t('forgot.backToLogin')}</a></p>
{:else}
	<p class="lede">{t('forgot.lede')}</p>

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

		{#if form?.error}
			<p class="alert" role="alert">{form.error}</p>
		{/if}

		<button type="submit" class="btn btn-primary" disabled={busy.value}>{t('forgot.submit')}</button
		>
	</form>

	<p class="alt"><a href={link(resolve('/(auth)/login'))}>{t('forgot.backToLogin')}</a></p>
{/if}

<style>
	.lede {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.alert.ok {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}
</style>
