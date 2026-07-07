<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { MessageKey } from '$lib/i18n';
	import { activeSection, sections, type Section } from './sections';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, link } = getI18n();
	const profile = pendingForm();
	const password = pendingForm();
	const sessions = pendingForm();

	// Mirrors the API allowlist (backend/internal/user); the catalog's country set.
	const countries = ['AT', 'BE', 'DE', 'ES', 'FR', 'IT', 'NL', 'PL', 'PT'];

	const active = $derived(activeSection(page.url.searchParams.get('section')));

	const navLabel: Record<Section, MessageKey> = {
		profile: 'settings.profile',
		security: 'settings.security',
		account: 'settings.account'
	};
	const hint: Record<Section, MessageKey> = {
		profile: 'settings.profileHint',
		security: 'settings.securityHint',
		account: 'settings.accountHint'
	};
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<h1>{t('settings.heading')}</h1>

<div class="wrap">
	<nav class="rail" aria-label={t('settings.navAria')}>
		{#each sections as s (s)}
			<a
				href="{link(resolve('/settings'))}?section={s}"
				aria-current={active === s ? 'page' : undefined}
				class:active={active === s}>{t(navLabel[s])}</a
			>
		{/each}
	</nav>

	<div class="pane">
		<h2>{t(navLabel[active])}</h2>
		<p class="hint">{t(hint[active])}</p>

		{#if active === 'profile'}
			<section class="panel">
				<form method="POST" action="?/profile" use:enhance={profile.submit}>
					<label for="display_name">{t('register.displayName')}</label>
					<input
						id="display_name"
						name="display_name"
						type="text"
						required
						minlength="2"
						maxlength="50"
						class="field"
						value={form?.value ?? data.user.display_name}
					/>

					<label for="locale">{t('lang.switch')}</label>
					<select id="locale" name="locale" class="field" value={data.user.locale}>
						<option value="en">{t('lang.en')}</option>
						<option value="es">{t('lang.es')}</option>
					</select>

					<label for="country">{t('settings.country')}</label>
					<select id="country" name="country" class="field" value={data.user.country ?? ''}>
						<option value="">{t('settings.countryNotSet')}</option>
						{#each countries as c (c)}<option value={c}>{c}</option>{/each}
					</select>

					{#if form?.error}
						<p class="alert" role="alert">{form.error}</p>
					{:else if form?.success}
						<p class="alert ok" role="status">{t('settings.profileUpdated')}</p>
					{/if}

					<button type="submit" class="btn btn-primary" disabled={profile.busy.value}
						>{t('settings.save')}</button
					>
				</form>
			</section>
		{:else if active === 'security'}
			<section class="panel">
				<h3>{t('settings.password')}</h3>

				<form
					method="POST"
					action="?section=security&/changePassword"
					use:enhance={password.submit}
				>
					<label for="current_password">{t('settings.currentPassword')}</label>
					<input
						id="current_password"
						name="current_password"
						type="password"
						required
						autocomplete="current-password"
						class="field"
					/>

					<label for="new_password">{t('reset.newPassword')}</label>
					<input
						id="new_password"
						name="new_password"
						type="password"
						required
						minlength="8"
						autocomplete="new-password"
						class="field"
					/>

					<label for="confirm_password">{t('settings.confirmNewPassword')}</label>
					<input
						id="confirm_password"
						name="confirm_password"
						type="password"
						required
						minlength="8"
						autocomplete="new-password"
						class="field"
					/>

					{#if form?.pwError}
						<p class="alert" role="alert">{form.pwError}</p>
					{:else if form?.pwSuccess}
						<p class="alert ok" role="status">{t('settings.passwordChanged')}</p>
					{/if}

					<button type="submit" class="btn btn-primary" disabled={password.busy.value}
						>{t('settings.changePassword')}</button
					>
				</form>
			</section>

			<section class="panel">
				<h3>{t('settings.sessions')}</h3>

				<p class="note">{t('settings.sessionsNote')}</p>

				<form method="POST" action="?section=security&/logoutAll" use:enhance={sessions.submit}>
					<button type="submit" class="btn btn-primary" disabled={sessions.busy.value}
						>{t('settings.logoutAll')}</button
					>
				</form>
			</section>
		{:else}
			<section class="panel">
				<h3>{t('settings.deleteAccount')}</h3>

				<p class="note">{t('settings.deleteNote')}</p>

				<button type="button" class="btn btn-outline" disabled title={t('settings.deleteTitle')}>
					{t('settings.deleteSoon')}
				</button>
			</section>
		{/if}
	</div>
</div>

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.wrap {
		margin-top: 1.5rem;
		display: grid;
		grid-template-columns: 14rem 1fr;
		gap: 2rem;
		align-items: start;
	}

	.rail {
		position: sticky;
		top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.rail a {
		padding: 0.5rem 0.75rem;
		border-radius: 0.375rem;
		border-left: 2px solid transparent;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 500;
		color: var(--slate-600);
		transition:
			background-color 0.15s,
			color 0.15s;
	}

	.rail a:hover {
		background: var(--slate-100);
	}

	.rail a.active {
		color: var(--brand-700);
		background: var(--brand-50);
		border-left-color: var(--brand-600);
	}

	.pane {
		min-width: 0;
		max-width: 34rem;
	}

	h2 {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.hint {
		margin-top: 0.125rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	h3 {
		font-size: 1rem;
		line-height: 1.5rem;
		font-weight: 600;
	}

	.panel {
		margin-top: 1rem;
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

	.field {
		width: 100%;
	}

	.alert.ok {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	form button,
	.panel > .btn {
		margin-top: 0.25rem;
	}

	/* Narrow screens: the rail becomes a horizontal strip above the pane. */
	@media (max-width: 47.9375rem) {
		.wrap {
			grid-template-columns: 1fr;
			gap: 1rem;
		}

		.rail {
			position: static;
			flex-direction: row;
			overflow-x: auto;
		}

		.rail a {
			white-space: nowrap;
			border-left: none;
			border-bottom: 2px solid transparent;
			border-radius: 0.375rem 0.375rem 0 0;
		}

		.rail a.active {
			border-bottom-color: var(--brand-600);
		}
	}
</style>
