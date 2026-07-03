<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, link } = getI18n();
	const { busy, submit } = pendingForm();
</script>

<svelte:head>
	<title>{t('reset.title')}</title>
</svelte:head>

<section class="panel">
	<h1>{t('reset.heading')}</h1>

	{#if form?.done}
		<p class="alert ok" role="status">{t('reset.done')}</p>
		<p class="alt"><a href={link(resolve('/login'))}>{t('nav.login')}</a></p>
	{:else if !data.token}
		<p class="alert" role="alert">{t('reset.missingToken')}</p>
		<p class="alt"><a href={link(resolve('/forgot'))}>{t('reset.requestLink')}</a></p>
	{:else}
		<form method="POST" use:enhance={submit}>
			<input type="hidden" name="token" value={data.token} />

			<label for="password">{t('reset.newPassword')}</label>
			<input
				id="password"
				name="password"
				type="password"
				required
				minlength="8"
				autocomplete="new-password"
				class="field"
			/>

			<label for="confirm">{t('reset.confirmPassword')}</label>
			<input
				id="confirm"
				name="confirm"
				type="password"
				required
				minlength="8"
				autocomplete="new-password"
				class="field"
			/>

			{#if form?.error}
				<p class="alert" role="alert">{form.error}</p>
			{/if}

			<button type="submit" class="btn btn-primary" disabled={busy.value}
				>{t('reset.submit')}</button
			>
		</form>
	{/if}
</section>

<style>
	.panel {
		margin-inline: auto;
		max-width: 24rem;
	}

	h1 {
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 700;
	}

	form {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	label {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		color: var(--slate-600);
	}

	.alert {
		margin-top: 0.75rem;
	}

	.alert.ok {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	button {
		margin-top: 1rem;
	}

	.alt {
		margin-top: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.alt a {
		color: var(--emerald-700);
		text-decoration: underline;
	}
</style>
