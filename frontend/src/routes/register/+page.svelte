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
	<title>{t('register.title')}</title>
</svelte:head>

<section class="panel">
	<h1>{t('register.heading')}</h1>

	<form method="POST" use:enhance={submit}>
		<label for="display_name">{t('register.displayName')}</label>
		<input
			id="display_name"
			name="display_name"
			type="text"
			required
			minlength="2"
			maxlength="50"
			class="field"
			value={form?.display_name ?? ''}
		/>

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
		<!-- minlength mirrors the server rule; autocomplete="new-password"
		     tells password managers to offer a generated one. -->
		<input
			id="password"
			name="password"
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
			>{t('register.heading')}</button
		>
	</form>

	<p class="alt">
		{t('register.haveAccount')} <a href={link(resolve('/login'))}>{t('nav.login')}</a>
	</p>
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
