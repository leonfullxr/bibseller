<script lang="ts">
	import { resolve } from '$app/paths';
	import PolicyBadge from '$lib/components/PolicyBadge.svelte';
	import { formatDate } from '$lib/format';
	import { getI18n } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const { t, locale, link } = getI18n();
</script>

<svelte:head>
	<title>{t('sell.title')}</title>
</svelte:head>

<h1>{t('sell.heading')}</h1>
<p class="lede">{t('sell.lede')}</p>

{#if !data.verified}
	<p class="notice">
		{t('sell.verifyNotice')}
		<a href={link(resolve('/settings'))}>{t('listingDetail.accountSettings')}</a>
	</p>
{/if}

<form method="GET" action={link(resolve('/sell'))} class="search">
	<input
		type="search"
		name="q"
		value={data.q}
		placeholder={t('races.filter.searchPlaceholder')}
		aria-label={t('sell.searchAria')}
	/>
	<button type="submit">{t('home.search')}</button>
</form>

{#if data.races.length === 0}
	<p class="empty">
		{t('sell.emptyPre')}
		<a href={link(resolve('/races'))}>{t('sell.browseAllLink')}</a>.
	</p>
{:else}
	<ul class="races">
		{#each data.races as race (race.id)}
			<li>
				<div class="info">
					<a href={link(resolve('/sell/[slug]', { slug: race.slug }))} class="name">{race.name}</a>
					<p class="meta">{formatDate(race.event_date, locale)} - {race.city}, {race.country}</p>
				</div>
				<div class="right">
					<PolicyBadge policy={race.transfer_policy} />
					<a href={link(resolve('/sell/[slug]', { slug: race.slug }))} class="sell"
						>{t('sell.sellHere')}</a
					>
				</div>
			</li>
		{/each}
	</ul>
{/if}

<style>
	h1 {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.lede {
		margin-top: 0.25rem;
		color: var(--slate-600);
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

	.search {
		margin-top: 1.5rem;
		display: flex;
		gap: 0.5rem;
	}

	.search input {
		flex: 1;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
	}

	.search button {
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	.empty {
		margin-top: 2rem;
		color: var(--slate-600);
	}

	.empty a,
	.races .name {
		color: var(--emerald-700);
		text-decoration: underline;
	}

	.races {
		margin-top: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		list-style: none;
		padding: 0;
	}

	.races li {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		border: 1px solid var(--slate-200);
		border-radius: 0.5rem;
		background: white;
		padding: 0.75rem 1rem;
	}

	.name {
		font-weight: 600;
	}

	.meta {
		margin-top: 0.125rem;
		font-size: 0.75rem;
		color: var(--slate-500);
	}

	.right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		white-space: nowrap;
	}

	.sell {
		border-radius: 0.375rem;
		background: var(--emerald-600);
		padding: 0.375rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: white;
	}

	.sell:hover {
		background: var(--emerald-700);
	}
</style>
