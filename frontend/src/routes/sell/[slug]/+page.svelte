<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import ListingFields from '$lib/components/ListingFields.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate } from '$lib/format';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, locale, link } = getI18n();
	const race = $derived(data.race);
	const { busy, submit } = pendingForm();
</script>

<svelte:head>
	<title>{t('sellForm.title', { name: race.name })}</title>
</svelte:head>

<nav><a href={link(resolve('/sell'))}>{t('sellForm.back')}</a></nav>

<header>
	<div class="title-row">
		<h1>{t('sellForm.heading')}</h1>
		<PolicyBadge policy={race.transfer_policy} />
	</div>
	<p class="meta">
		{race.name} - {formatDate(race.event_date, locale)} - {race.city}, {race.country}
	</p>
</header>

<div class="callout-wrap">
	<PolicyCallout
		policy={race.transfer_policy}
		officialUrl={race.official_transfer_url}
		notes={race.policy_notes}
	/>
</div>

{#if !data.verified}
	<p class="alert notice">
		{t('sellForm.verifyNotice')}
		<a href={link(resolve('/settings'))}>{t('listingDetail.accountSettings')}</a>
	</p>
{:else}
	<form method="POST" use:enhance={submit} class="panel">
		<input type="hidden" name="race_id" value={race.id} />

		<ListingFields
			price={form?.values?.price ?? ''}
			original={form?.values?.original_price ?? ''}
			description={form?.values?.description ?? ''}
		/>

		{#if form?.error}
			<p class="alert feedback" role="alert">{form.error}</p>
		{/if}

		<button type="submit" class="btn btn-primary" disabled={busy.value}
			>{t('sellForm.publish')}</button
		>
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
	}

	.feedback {
		margin-top: 0.5rem;
	}

	button {
		margin-top: 1rem;
		align-self: flex-start;
	}
</style>
