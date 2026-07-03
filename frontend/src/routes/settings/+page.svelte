<script lang="ts">
	import { enhance } from '$app/forms';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t } = getI18n();
	const profile = pendingForm();
	const password = pendingForm();
	const sessions = pendingForm();

	// Mirrors the API allowlist (backend/internal/user); the catalog's country set.
	const countries = ['AT', 'BE', 'DE', 'ES', 'FR', 'IT', 'NL', 'PL', 'PT'];
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<h1>{t('settings.heading')}</h1>

<section class="panel">
	<h2>{t('settings.profile')}</h2>

	<form method="POST" action="?/profile" use:enhance={profile.submit}>
		<label for="display_name">{t('register.displayName')}</label>
		<input
			id="display_name"
			name="display_name"
			type="text"
			required
			minlength="2"
			maxlength="50"
			value={form?.value ?? data.user.display_name}
		/>

		<label for="locale">{t('lang.switch')}</label>
		<select id="locale" name="locale" value={data.user.locale}>
			<option value="en">{t('lang.en')}</option>
			<option value="es">{t('lang.es')}</option>
		</select>

		<label for="country">{t('settings.country')}</label>
		<select id="country" name="country" value={data.user.country ?? ''}>
			<option value="">{t('settings.countryNotSet')}</option>
			{#each countries as c (c)}<option value={c}>{c}</option>{/each}
		</select>

		{#if form?.error}
			<p class="feedback error" role="alert">{form.error}</p>
		{:else if form?.success}
			<p class="feedback success" role="status">{t('settings.profileUpdated')}</p>
		{/if}

		<button type="submit" disabled={profile.busy.value}>{t('settings.save')}</button>
	</form>
</section>

<section class="panel">
	<h2>{t('settings.password')}</h2>

	<form method="POST" action="?/changePassword" use:enhance={password.submit}>
		<label for="current_password">{t('settings.currentPassword')}</label>
		<input
			id="current_password"
			name="current_password"
			type="password"
			required
			autocomplete="current-password"
		/>

		<label for="new_password">{t('reset.newPassword')}</label>
		<input
			id="new_password"
			name="new_password"
			type="password"
			required
			minlength="8"
			autocomplete="new-password"
		/>

		<label for="confirm_password">{t('settings.confirmNewPassword')}</label>
		<input
			id="confirm_password"
			name="confirm_password"
			type="password"
			required
			minlength="8"
			autocomplete="new-password"
		/>

		{#if form?.pwError}
			<p class="feedback error" role="alert">{form.pwError}</p>
		{:else if form?.pwSuccess}
			<p class="feedback success" role="status">{t('settings.passwordChanged')}</p>
		{/if}

		<button type="submit" disabled={password.busy.value}>{t('settings.changePassword')}</button>
	</form>
</section>

<section class="panel">
	<h2>{t('settings.sessions')}</h2>

	<p class="note">{t('settings.sessionsNote')}</p>

	<form method="POST" action="?/logoutAll" use:enhance={sessions.submit}>
		<button type="submit" disabled={sessions.busy.value}>{t('settings.logoutAll')}</button>
	</form>
</section>

<section class="panel">
	<h2>{t('settings.deleteAccount')}</h2>

	<p class="note">{t('settings.deleteNote')}</p>

	<button type="button" disabled title={t('settings.deleteTitle')}>
		{t('settings.deleteSoon')}
	</button>
</section>

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.panel {
		margin-top: 1.5rem;
		max-width: 28rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem;
	}

	h2 {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.note {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	form {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.5rem;
	}

	label {
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		color: var(--slate-600);
	}

	input,
	select {
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	.feedback {
		border-radius: 0.375rem;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 500;
	}

	.error {
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		color: var(--amber-900);
	}

	.success {
		border: 1px solid var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	button {
		margin-top: 0.25rem;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: white;
	}

	button:hover:not(:disabled) {
		background: var(--slate-700);
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
