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

<h1>{t('reset.heading')}</h1>

{#if form?.done}
	<p class="alert ok" role="status">{t('reset.done')}</p>
	<p class="alt"><a href={link(resolve('/(auth)/login'))}>{t('nav.login')}</a></p>
{:else if !data.token}
	<p class="alert" role="alert">{t('reset.missingToken')}</p>
	<p class="alt"><a href={link(resolve('/(auth)/forgot'))}>{t('reset.requestLink')}</a></p>
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

		<button type="submit" class="btn btn-primary" disabled={busy.value}>{t('reset.submit')}</button>
	</form>
{/if}
