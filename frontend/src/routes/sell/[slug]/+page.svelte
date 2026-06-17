<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import ListingFields from '$lib/components/ListingFields.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate } from '$lib/format';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const race = $derived(data.race);
</script>

<svelte:head>
	<title>List your bib for {race.name} - Bibseller</title>
</svelte:head>

<nav><a href={resolve('/sell')}>Back to race search</a></nav>

<header>
	<div class="title-row">
		<h1>List your bib</h1>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="meta">{race.name} - {formatDate(race.event_date)} - {race.city}, {race.country}</p>
</header>

<div class="callout-wrap">
	<PolicyCallout
		policy={race.transfer_policy}
		officialUrl={race.official_transfer_url}
		notes={race.policy_notes}
	/>
</div>

{#if !data.verified}
	<p class="notice">
		Verify your email to publish a listing.
		<a href={resolve('/settings')}>Account settings</a>
	</p>
{:else}
	<form method="POST" use:enhance class="panel">
		<input type="hidden" name="race_id" value={race.id} />

		<ListingFields
			price={form?.values?.price ?? ''}
			original={form?.values?.original_price ?? ''}
			description={form?.values?.description ?? ''}
		/>

		{#if form?.error}
			<p class="feedback error" role="alert">{form.error}</p>
		{/if}

		<button type="submit">Publish listing</button>
	</form>
{/if}

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

	.title-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
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

	.callout-wrap {
		margin-top: 1rem;
	}

	.notice {
		margin-top: 1rem;
		border-radius: 0.375rem;
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		color: var(--amber-900);
	}

	.notice a {
		color: var(--amber-900);
		text-decoration: underline;
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
		background: var(--emerald-600);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	button:hover {
		background: var(--emerald-700);
	}
</style>
