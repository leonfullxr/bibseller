<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
</script>

<svelte:head>
	<title>Set a new password - Bibseller</title>
</svelte:head>

<section class="panel">
	<h1>Set a new password</h1>

	{#if form?.done}
		<p class="feedback ok" role="status">
			Your password has been updated. You've been signed out everywhere - sign in with your new
			password.
		</p>
		<p class="alt"><a href={resolve('/login')}>Log in</a></p>
	{:else if !data.token}
		<p class="feedback" role="alert">This reset link is missing its token. Request a new one.</p>
		<p class="alt"><a href={resolve('/forgot')}>Request a reset link</a></p>
	{:else}
		<form method="POST" use:enhance>
			<input type="hidden" name="token" value={data.token} />

			<label for="password">New password</label>
			<input
				id="password"
				name="password"
				type="password"
				required
				minlength="8"
				autocomplete="new-password"
			/>

			<label for="confirm">Confirm password</label>
			<input
				id="confirm"
				name="confirm"
				type="password"
				required
				minlength="8"
				autocomplete="new-password"
			/>

			{#if form?.error}
				<p class="feedback" role="alert">{form.error}</p>
			{/if}

			<button type="submit">Update password</button>
		</form>
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
