<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
</script>

<svelte:head>
	<title>Settings - Bibseller</title>
</svelte:head>

<h1>Settings</h1>

<section class="panel">
	<h2>Profile</h2>

	<form method="POST" action="?/profile" use:enhance>
		<label for="display_name">Display name</label>
		<input
			id="display_name"
			name="display_name"
			type="text"
			required
			minlength="2"
			maxlength="50"
			value={form?.value ?? data.user.display_name}
		/>

		{#if form?.error}
			<p class="feedback error" role="alert">{form.error}</p>
		{:else if form?.success}
			<p class="feedback success" role="status">Display name updated.</p>
		{/if}

		<button type="submit">Save</button>
	</form>
</section>

<section class="panel">
	<h2>Password</h2>

	<form method="POST" action="?/changePassword" use:enhance>
		<label for="current_password">Current password</label>
		<input
			id="current_password"
			name="current_password"
			type="password"
			required
			autocomplete="current-password"
		/>

		<label for="new_password">New password</label>
		<input
			id="new_password"
			name="new_password"
			type="password"
			required
			minlength="8"
			autocomplete="new-password"
		/>

		<label for="confirm_password">Confirm new password</label>
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
			<p class="feedback success" role="status">
				Password changed. Other devices have been signed out.
			</p>
		{/if}

		<button type="submit">Change password</button>
	</form>
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

	input {
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

	button:hover {
		background: var(--slate-700);
	}
</style>
