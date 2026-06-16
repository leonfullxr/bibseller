<script lang="ts">
	import { resolve } from '$app/paths';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head><title>Verify email — Bibseller</title></svelte:head>

<section class="verify">
	{#if data.status === 'ok'}
		<h1>Email verified</h1>
		<p>Your email address is confirmed — you're all set.</p>
		<a href={resolve('/')}>Continue</a>
	{:else if data.status === 'invalid'}
		<h1>Link invalid or expired</h1>
		<p>This verification link is no longer valid. Sign in and request a fresh one.</p>
		<a href={resolve('/login')}>Sign in</a>
	{:else if data.status === 'missing'}
		<h1>Nothing to verify</h1>
		<p>Open the verification link from your email to confirm your address.</p>
		<a href={resolve('/')}>Home</a>
	{:else}
		<h1>Something went wrong</h1>
		<p>We couldn't verify your email right now. Please try again in a moment.</p>
		<a href={resolve('/')}>Home</a>
	{/if}
</section>

<style>
	.verify {
		max-width: 28rem;
		margin-inline: auto;
		text-align: center;
	}

	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	p {
		margin-top: 0.5rem;
		color: var(--slate-600);
	}

	a {
		display: inline-block;
		margin-top: 1.5rem;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	a:hover {
		background: var(--slate-700);
	}
</style>
