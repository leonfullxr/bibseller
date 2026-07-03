<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { ResolvedPathname } from '$app/types';
	import ListingCTA from '$lib/components/ListingCTA.svelte';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import PolicyCallout from '$lib/components/PolicyCallout.svelte';
	import { formatDate, formatPrice } from '$lib/format';
	import { pendingForm } from '$lib/forms.svelte';
	import { getI18n, listingStatusLabel } from '$lib/i18n';
	import { policyView, requiresAck } from '$lib/policy';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const { t, locale, link } = getI18n();
	const { busy, submit } = pendingForm();
	// Send the visitor back here after logging in. The login action validates
	// ?next= server-side (same-site paths only). Cast: a resolved path plus a
	// query string is still one, which is what no-navigation-without-resolve checks.
	const loginHref = $derived(
		`${link(resolve('/login'))}?next=${encodeURIComponent(page.url.pathname + page.url.search)}` as ResolvedPathname
	);
	const listing = $derived(data.listing);
	const race = $derived(data.listing.race);

	const price = $derived(formatPrice(listing.price_cents, listing.currency, locale));
	const original = $derived(formatPrice(listing.original_price_cents, listing.currency, locale));
	const belowFace = $derived(
		listing.price_cents != null &&
			listing.original_price_cents != null &&
			listing.price_cents < listing.original_price_cents
	);
	const available = $derived(listing.status === 'active');
	const isOwn = $derived(listing.is_own_listing);
	const needsAck = $derived(requiresAck(race.transfer_policy));
	// Whether ListingCTA will actually render a control - mirrors its own
	// branches, so no empty wrapper reserves space in the other modes.
	const hasCta = $derived.by(() => {
		const action = policyView[race.transfer_policy].primaryAction;
		return action === 'buy' || (action === 'official' && !!race.official_transfer_url);
	});

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
			reportStatus = res.ok ? t('report.success') : t('report.failed');
			if (res.ok) reportDetails = '';
		} catch {
			reportStatus = t('report.networkError');
		} finally {
			reporting = false;
		}
	}
</script>

<svelte:head>
	<title>{t('listingDetail.title', { name: race.name })}</title>
</svelte:head>

<nav>
	<a href={link(resolve('/races/[slug]', { slug: race.slug }))}>
		{t('listingDetail.back', { name: race.name })}
	</a>
</nav>

<div class="panel" class:unavailable={!available}>
	<div class="head">
		<div>
			<h1>{t('listingDetail.heading', { name: race.name })}</h1>
			<p class="meta">
				{formatDate(race.event_date, locale)} - {race.city}, {race.country}
				{#if race.distance}
					- {race.distance}{/if}
			</p>
		</div>
		<PolicyBadge policy={race.transfer_policy} />
	</div>

	{#if !available}
		<div class="gone">
			{t('listingDetail.unavailable', { status: listingStatusLabel(t, listing.status) })}
		</div>
	{/if}

	<div class="price-row">
		<span class="price">{price ?? t('listingCard.priceOnRequest')}</span>
		{#if belowFace && original}
			<span class="original">{original}</span>
			<span class="pill deal">{t('listingCard.belowFace')}</span>
		{/if}
	</div>

	{#if listing.description}
		<p class="desc">{listing.description}</p>
	{/if}
	<p class="listed-by">
		{t('listingDetail.listedByOn', {
			name: listing.seller_name,
			date: formatDate(listing.created_at.slice(0, 10), locale)
		})}
	</p>

	{#if available && hasCta}
		<div class="cta-wrap">
			<ListingCTA policy={race.transfer_policy} officialUrl={race.official_transfer_url} />
		</div>
	{/if}
</div>

<!-- No officialUrl here: the CTA above owns the official-transfer action, so
     the callout stays informational (one button, not two). -->
<div class="callout-wrap">
	<PolicyCallout policy={race.transfer_policy} />
</div>

{#if available}
	<section class="panel contact">
		<h2>{t('listingDetail.contact')}</h2>

		{#if !data.user}
			<p class="hint">
				<a href={loginHref}>{t('nav.login')}</a>
				{t('listingDetail.toMessageSeller')}
			</p>
		{:else if !data.user.email_verified}
			<p class="hint">
				{t('listingDetail.verifyToMessage')}
				<a href={link(resolve('/settings'))}>{t('listingDetail.accountSettings')}</a>
			</p>
		{:else if isOwn}
			<p class="hint">
				{t('listingDetail.ownPre')}
				<a href={link(resolve('/account/listings'))}>{t('listingDetail.yourListings')}</a>.
			</p>
		{:else}
			<form method="POST" action="?/contact" use:enhance={submit} class="composer">
				<textarea
					name="body"
					class="field"
					rows="4"
					required
					maxlength="4000"
					aria-label={t('listingDetail.messageAria')}
					placeholder={t('listingDetail.messagePlaceholder')}>{form?.body ?? ''}</textarea
				>

				{#if needsAck}
					<label class="ack">
						<input type="checkbox" name="ack" required />
						<span>{t('listingDetail.ackText')}</span>
					</label>
				{/if}

				{#if form?.error}
					<p class="alert" role="alert">{form.error}</p>
				{/if}

				<button type="submit" class="btn btn-primary" disabled={busy.value}
					>{t('listingDetail.send')}</button
				>
			</form>
		{/if}
	</section>
{/if}

{#if data.user}
	<details class="report">
		<summary>{t('report.summary')}</summary>
		<form class="report-form" onsubmit={reportListing}>
			<select class="field" bind:value={reportReason} aria-label={t('report.reasonAria')}>
				<option value="forbidden_transfer">{t('report.reason.forbidden_transfer')}</option>
				<option value="scam">{t('report.reason.scam')}</option>
				<option value="offensive">{t('report.reason.offensive')}</option>
				<option value="other">{t('report.reason.other')}</option>
			</select>
			<textarea
				class="field"
				bind:value={reportDetails}
				rows="2"
				maxlength="2000"
				aria-label={t('report.detailsAria')}
				placeholder={t('report.detailsPlaceholder')}
			></textarea>
			{#if reportStatus}<p class="report-status" role="status">{reportStatus}</p>{/if}
			<button type="submit" class="btn btn-primary" disabled={reporting}
				>{reporting ? t('report.submitting') : t('report.submit')}</button
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
		background: var(--emerald-100);
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

	button {
		align-self: flex-start;
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
		padding: 0.375rem 0.5rem;
		font: inherit;
	}

	.report-status {
		font-size: 0.8125rem;
		color: var(--slate-700);
	}
</style>
