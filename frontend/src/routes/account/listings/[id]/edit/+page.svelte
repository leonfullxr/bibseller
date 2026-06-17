<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import ListingFields from '$lib/components/ListingFields.svelte';
	import { formatDate } from '$lib/format';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const l = $derived(data.listing);

	function centsToInput(c: number | null): string {
		return c == null ? '' : String(c / 100);
	}
</script>

<svelte:head>
	<title>Edit listing - Bibseller</title>
</svelte:head>

<nav><a href={resolve('/account/listings')}>Back to my listings</a></nav>

<header>
	<h1>Edit listing</h1>
	<p class="meta">{l.race.name} - {formatDate(l.race.event_date)}</p>
</header>

<form method="POST" use:enhance class="panel">
	<ListingFields
		price={form?.values?.price ?? centsToInput(l.price_cents)}
		original={form?.values?.original_price ?? centsToInput(l.original_price_cents)}
		description={form?.values?.description ?? l.description ?? ''}
	/>

	{#if form?.error}
		<p class="feedback error" role="alert">{form.error}</p>
	{/if}

	<button type="submit">Save changes</button>
</form>

<style>
	nav {
		font-size: 0.875rem;
	}

	nav a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	header {
		margin-top: 1rem;
	}

	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.meta {
		margin-top: 0.25rem;
		font-size: 0.875rem;
		color: var(--slate-600);
	}

	.panel {
		margin-top: 1.5rem;
		max-width: 28rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem;
	}

	.feedback {
		margin-top: 0.5rem;
		border-radius: 0.375rem;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.error {
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		color: var(--amber-900);
	}

	button {
		margin-top: 1rem;
		align-self: flex-start;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	button:hover {
		background: var(--slate-700);
	}
</style>
