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

<section class="panel">
	<h1>{t('forgot.heading')}</h1>

	{#if form?.sent}
		<p class="feedback ok" role="status">{t('forgot.sent')}</p>
		<p class="alt"><a href={link(resolve('/login'))}>{t('forgot.backToLogin')}</a></p>
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
				value={form?.email ?? ''}
			/>

			{#if form?.error}
				<p class="feedback" role="alert">{form.error}</p>
			{/if}

			<button type="submit" disabled={busy.value}>{t('forgot.submit')}</button>
		</form>

		<p class="alt"><a href={link(resolve('/login'))}>{t('forgot.backToLogin')}</a></p>
	{/if}
</section>

<style>
	.panel {
		margin-inline: auto;
		max-width: 24rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem;
	}

	h1 {
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 700;
	}

	.lede {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
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

	input {
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.feedback {
		margin-top: 0.75rem;
		border-radius: 0.375rem;
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 500;
		color: var(--amber-900);
	}

	.feedback.ok {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	button {
		margin-top: 1rem;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: white;
	}

	button:hover {
		background: var(--slate-700);
	}

	button:disabled {
		opacity: 0.6;
		cursor: default;
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
