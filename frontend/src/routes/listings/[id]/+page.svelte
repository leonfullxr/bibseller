<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import ListingCTA from '$lib/components/ListingCTA.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate, formatPrice } from '$lib/format';
	import { requiresAck } from '$lib/policy';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const listing = $derived(data.listing);
	const race = $derived(data.listing.race);

	const price = $derived(formatPrice(listing.price_cents, listing.currency));
	const original = $derived(formatPrice(listing.original_price_cents, listing.currency));
	const belowFace = $derived(
		listing.price_cents != null &&
			listing.original_price_cents != null &&
			listing.price_cents < listing.original_price_cents
	);
	const available = $derived(listing.status === 'active');
	const isOwn = $derived(listing.is_own_listing);
	const needsAck = $derived(requiresAck(race.transfer_policy));

	let reportReason = $state('scam');
	let reportDetails = $state('');
	let reportStatus = $state('');
	let reporting = $state(false);

	async function reportListing(e: SubmitEvent) {
		e.preventDefault();
		if (reporting) return;
		reporting = true;
		reportStatus = '';
		try {
			const res = await fetch('/api/v1/reports', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'same-origin',
				body: JSON.stringify({
					subject_type: 'listing',
					subject_id: listing.id,
					reason: reportReason,
					details: reportDetails.trim() || undefined
				})
			});
			reportStatus = res.ok
				? 'Thanks - this listing has been reported.'
				: 'Could not file the report. Try again.';
			if (res.ok) reportDetails = '';
		} catch {
			reportStatus = 'Network error - try again.';
		} finally {
			reporting = false;
		}
	}
</script>

<svelte:head>
	<title>Bib for {race.name} - Bibseller</title>
</svelte:head>

<nav>
	<a href={resolve('/races/[slug]', { slug: race.slug })}>
		Back to {race.name}
	</a>
</nav>

<div class="panel" class:unavailable={!available}>
	<div class="head">
		<div>
			<h1>Bib for {race.name}</h1>
			<p class="meta">
				{formatDate(race.event_date)} - {race.city}, {race.country}
				{#if race.distance}
					- {race.distance}{/if}
			</p>
		</div>
		<PolicyBadge policy={race.transfer_policy} />
	</div>

	{#if !available}
		<div class="gone">
			This listing is no longer available ({listing.status}).
		</div>
	{/if}

	<div class="price-row">
		<span class="price">{price ?? 'Price on request'}</span>
		{#if belowFace && original}
			<span class="original">{original}</span>
			<span class="deal">below face value</span>
		{/if}
	</div>

	{#if listing.description}
		<p class="desc">{listing.description}</p>
	{/if}
	<p class="listed-by">
		Listed by {listing.seller_name} on {formatDate(listing.created_at.slice(0, 10))}
	</p>

	{#if available}
		<div class="cta-wrap">
			<ListingCTA policy={race.transfer_policy} officialUrl={race.official_transfer_url} />
		</div>
	{/if}
</div>

<div class="callout-wrap">
	<PolicyCallout policy={race.transfer_policy} officialUrl={race.official_transfer_url} />
</div>

{#if available}
	<section class="contact">
		<h2>Contact the seller</h2>

		{#if !data.user}
			<p class="hint">
				<a href={resolve('/login')}>Log in</a> to message the seller.
			</p>
		{:else if !data.user.email_verified}
			<p class="hint">
				Verify your email to message the seller.
				<a href={resolve('/settings')}>Account settings</a>
			</p>
		{:else if isOwn}
			<p class="hint">
				This is your listing - manage it from
				<a href={resolve('/account/listings')}>your listings</a>.
			</p>
		{:else}
			<form method="POST" action="?/contact" use:enhance class="composer">
				<textarea
					name="body"
					rows="4"
					required
					maxlength="4000"
					aria-label="Message to the seller"
					placeholder="Hi - is this bib still available?">{form?.body ?? ''}</textarea
				>

				{#if needsAck}
					<label class="ack">
						<input type="checkbox" name="ack" required />
						<span>
							I understand the platform handles no money and takes no responsibility for this
							transfer - the race's own rules apply.
						</span>
					</label>
				{/if}

				{#if form?.error}
					<p class="feedback error" role="alert">{form.error}</p>
				{/if}

				<button type="submit">Send message</button>
			</form>
		{/if}
	</section>
{/if}

{#if data.user}
	<details class="report">
		<summary>Report this listing</summary>
		<form class="report-form" onsubmit={reportListing}>
			<select bind:value={reportReason} aria-label="Reason for report">
				<option value="forbidden_transfer">Forbidden transfer</option>
				<option value="scam">Scam</option>
				<option value="offensive">Offensive</option>
				<option value="other">Other</option>
			</select>
			<textarea
				bind:value={reportDetails}
				rows="2"
				maxlength="2000"
				placeholder="Details (optional)"
			></textarea>
			{#if reportStatus}<p class="report-status" role="status">{reportStatus}</p>{/if}
			<button type="submit" disabled={reporting}
				>{reporting ? 'Reporting...' : 'Submit report'}</button
			>
		</form>
	</details>
{/if}

<style>
	nav {
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	nav a {
		color: var(--slate-500);
	}

	nav a:hover {
		color: var(--slate-900);
	}

	.panel {
		margin-top: 1rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem;
	}

	.panel.unavailable {
		opacity: 0.75;
	}

	.head {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
	}

	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.meta {
		margin-top: 0.25rem;
		color: var(--slate-600);
	}

	.gone {
		margin-top: 1rem;
		border-radius: 0.375rem;
		background: var(--slate-100);
		padding: 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: var(--slate-700);
	}

	.price-row {
		margin-top: 1.5rem;
		display: flex;
		align-items: baseline;
		gap: 0.75rem;
	}

	.price {
		font-size: 2.25rem;
		line-height: 2.5rem;
		font-weight: 800;
		letter-spacing: -0.025em;
	}

	.original {
		font-size: 1.125rem;
		line-height: 1.75rem;
		color: var(--slate-400);
		text-decoration: line-through;
	}

	.deal {
		border-radius: 9999px;
		background: var(--emerald-100);
		padding: 0.125rem 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 600;
		color: var(--emerald-800);
	}

	.desc {
		margin-top: 1rem;
		max-width: 65ch;
		color: var(--slate-700);
	}

	.listed-by {
		margin-top: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	.cta-wrap {
		margin-top: 1.5rem;
	}

	.callout-wrap {
		margin-top: 1.5rem;
	}

	.contact {
		margin-top: 1.5rem;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-200);
		background: white;
		padding: 1.5rem;
	}

	.contact h2 {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.hint {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.hint a {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.composer {
		margin-top: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		max-width: 40rem;
	}

	textarea {
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font: inherit;
		font-size: 0.875rem;
		resize: vertical;
	}

	.ack {
		display: flex;
		gap: 0.5rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-700);
	}

	.ack input {
		margin-top: 0.2rem;
		flex-shrink: 0;
	}

	.feedback {
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

	.report {
		margin-top: 1.5rem;
		font-size: 0.875rem;
	}

	.report summary {
		cursor: pointer;
		color: var(--slate-500);
	}

	.report-form {
		margin-top: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-width: 28rem;
	}

	select {
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.375rem 0.5rem;
		font: inherit;
		font-size: 0.875rem;
	}

	.report-status {
		font-size: 0.8125rem;
		color: var(--slate-700);
	}
</style>
