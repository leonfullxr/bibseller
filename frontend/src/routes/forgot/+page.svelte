<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import type { PageProps } from './$types';

	let { form }: PageProps = $props();
</script>

<svelte:head>
	<title>Reset password - Bibseller</title>
</svelte:head>

<section class="panel">
	<h1>Reset your password</h1>

	{#if form?.sent}
		<p class="feedback ok" role="status">
			If an account exists for that address, we've sent a link to reset your password. Check your
			inbox.
		</p>
		<p class="alt"><a href={resolve('/login')}>Back to log in</a></p>
	{:else}
		<p class="lede">Enter your email and we'll send you a reset link.</p>

		<form method="POST" use:enhance>
			<label for="email">Email</label>
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

			<button type="submit">Send reset link</button>
		</form>

		<p class="alt"><a href={resolve('/login')}>Back to log in</a></p>
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
