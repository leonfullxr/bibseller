<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import ListingFields from '$lib/components/ListingFields.svelte';
	import { formatDate } from '$lib/format';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, locale, link } = getI18n();
	const l = $derived(data.listing);
	const { busy, submit } = pendingForm();

	function centsToInput(c: number | null): string {
		return c == null ? '' : String(c / 100);
	}
</script>

<svelte:head>
	<title>{t('editListing.title')}</title>
</svelte:head>

<nav><a href={link(resolve('/account/listings'))}>{t('editListing.back')}</a></nav>

<header>
	<h1>{t('editListing.heading')}</h1>
	<p class="meta">{l.race.name} - {formatDate(l.race.event_date, locale)}</p>
</header>

<form method="POST" use:enhance={submit} class="panel">
	<ListingFields
		price={form?.values?.price ?? centsToInput(l.price_cents)}
		original={form?.values?.original_price ?? centsToInput(l.original_price_cents)}
		description={form?.values?.description ?? l.description ?? ''}
	/>

	{#if form?.error}
		<p class="alert feedback" role="alert">{form.error}</p>
	{/if}

	<button type="submit" class="btn btn-primary" disabled={busy.value}
		>{t('editListing.save')}</button
	>
</form>

<style>
	nav {
		font-size: 0.875rem;
	}

	nav a {
		color: var(--brand-700);
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
	}

	.feedback {
		margin-top: 0.5rem;
	}

	button {
		margin-top: 1rem;
		align-self: flex-start;
	}
</style>
