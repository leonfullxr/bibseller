<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ResolvedPathname } from '$app/types';
	import { page } from '$app/state';
	import {
		createPlural,
		createTranslator,
		pathForLocale,
		setI18n,
		stripLocale,
		type I18n
	} from '$lib/i18n';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children, data } = $props();

	// Build the bound translator + link helper for the active locale and share
	// them down the tree via context. Locale only changes on a full reload (the
	// first-visit redirect or the switcher's plain form POST), but keeping it
	// reactive costs nothing and keeps the helpers in sync with data.
	const locale = $derived(data.locale);
	const translate = $derived(createTranslator(locale));
	const pluralize = $derived(createPlural(locale));
	const i18n: I18n = {
		get locale() {
			return locale;
		},
		t: (key, params) => translate(key, params),
		plural: (base, n, params) => pluralize(base, n, params),
		link: (path) => pathForLocale(locale, path) as ResolvedPathname
	};
	setI18n(i18n);
	const { t, link } = i18n;

	const switchTo = $derived(data.locale === 'en' ? 'es' : 'en');
	// hreflang alternates for the current page, from its locale-free path (no query).
	const basePath = $derived(stripLocale(page.url.pathname));
	const enHref = $derived(page.url.origin + pathForLocale('en', basePath));
	const esHref = $derived(page.url.origin + pathForLocale('es', basePath));
	// The switcher returns to the same page including its query (e.g. race filters).
	const switchTarget = $derived(basePath + page.url.search);
	// The suggestion banner reads in the suggested locale (Spanish), so it shows
	// English now and Spanish once es.ts is filled (M8.2).
	const suggestT = $derived(createTranslator(data.suggestLocale ?? 'es'));
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="alternate" hreflang="en" href={enHref} />
	<link rel="alternate" hreflang="es" href={esHref} />
	<link rel="alternate" hreflang="x-default" href={enHref} />
</svelte:head>

<div class="shell">
	<header>
		<div class="bar">
			<a href={link(resolve('/'))} class="brand">bib<span>seller</span></a>
			<nav>
				<a href={link(resolve('/races'))}>{t('nav.races')}</a>
				{#if data.user}
					<a href={link(resolve('/sell'))}>{t('nav.sell')}</a>
					<a href={link(resolve('/account/listings'))}>{t('nav.myListings')}</a>
					{#if data.user.email_verified}
						<a href={link(resolve('/account/inbox'))}>{t('nav.inbox')}</a>
					{/if}
					<a href={link(resolve('/settings'))}>{data.user.display_name}</a>
					<form method="POST" action={link(resolve('/logout'))}>
						<button type="submit">{t('nav.logout')}</button>
					</form>
				{:else}
					<a href={link(resolve('/login'))}>{t('nav.login')}</a>
					<a href={link(resolve('/register'))}>{t('nav.register')}</a>
				{/if}
				<form method="POST" action={link(resolve('/locale'))} class="lang">
					<input type="hidden" name="to" value={switchTo} />
					<input type="hidden" name="next" value={switchTarget} />
					<button type="submit">{switchTo === 'es' ? t('lang.es') : t('lang.en')}</button>
				</form>
			</nav>
		</div>
	</header>

	{#if data.suggestLocale === 'es'}
		<div class="suggest-banner">
			<span>{suggestT('banner.suggestText')}</span>
			<form method="POST" action={link(resolve('/locale'))}>
				<input type="hidden" name="next" value={switchTarget} />
				<button type="submit" name="to" value="es">{suggestT('banner.suggestAccept')}</button>
				<button type="submit" name="to" value="en" class="ghost"
					>{suggestT('banner.suggestDismiss')}</button
				>
			</form>
		</div>
	{/if}

	{#if data.user && !data.user.email_verified}
		<div class="verify-banner">
			<span>{t('banner.verifyEmail')}</span>
			<form method="POST" action={link(resolve('/verify/resend'))}>
				<button type="submit">{t('banner.resend')}</button>
			</form>
		</div>
	{/if}

	<main>
		{@render children()}
	</main>

	<footer>
		<div class="bar foot">
			<span>{t('footer.tagline')}</span>
			<nav class="foot-links">
				<a href={link(resolve('/terms'))}>{t('footer.terms')}</a>
				<a href="https://github.com/leonfullxr/bibseller" rel="external">{t('footer.github')}</a>
			</nav>
		</div>
	</footer>
</div>

<style>
	.shell {
		display: flex;
		min-height: 100vh;
		flex-direction: column;
		background: var(--slate-50);
		color: var(--slate-900);
	}

	header {
		border-bottom: 1px solid var(--slate-200);
		background: white;
	}

	.bar {
		margin-inline: auto;
		display: flex;
		width: 100%;
		max-width: 64rem;
		align-items: center;
		justify-content: space-between;
		padding-inline: 1rem;
	}

	header .bar {
		height: 3.5rem;
	}

	.brand {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 700;
		letter-spacing: -0.025em;
	}

	.brand span {
		color: var(--emerald-600);
	}

	nav {
		display: flex;
		align-items: center;
		gap: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	nav a,
	nav button {
		font-weight: 500;
		color: var(--slate-600);
	}

	nav a:hover,
	nav button:hover {
		color: var(--slate-900);
	}

	nav form {
		display: contents;
	}

	nav button {
		cursor: pointer;
		border: none;
		background: none;
		padding: 0;
		font: inherit;
	}

	.verify-banner {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		background: var(--amber-50);
		border-bottom: 1px solid var(--amber-300);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		color: var(--amber-900);
	}

	.verify-banner form {
		display: contents;
	}

	.verify-banner button {
		cursor: pointer;
		border: 1px solid var(--amber-300);
		border-radius: 0.375rem;
		background: white;
		padding: 0.25rem 0.625rem;
		font: inherit;
		font-weight: 600;
		color: var(--amber-900);
	}

	.verify-banner button:hover {
		background: var(--amber-100);
	}

	.suggest-banner {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		background: var(--emerald-50);
		border-bottom: 1px solid var(--emerald-300);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		color: var(--emerald-900);
	}

	.suggest-banner form {
		display: contents;
	}

	.suggest-banner button {
		cursor: pointer;
		border-radius: 0.375rem;
		border: 1px solid var(--emerald-600);
		background: var(--emerald-600);
		padding: 0.25rem 0.625rem;
		font: inherit;
		font-weight: 600;
		color: white;
	}

	.suggest-banner button:hover {
		background: var(--emerald-700);
	}

	.suggest-banner button.ghost {
		border-color: var(--emerald-300);
		background: transparent;
		color: var(--emerald-900);
	}

	.suggest-banner button.ghost:hover {
		background: var(--emerald-100);
	}

	main {
		margin-inline: auto;
		width: 100%;
		max-width: 64rem;
		flex: 1;
		padding: 2.5rem 1rem;
	}

	footer {
		border-top: 1px solid var(--slate-200);
		background: white;
	}

	.foot {
		padding-block: 1rem;
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}

	.foot-links {
		display: flex;
		gap: 1rem;
	}

	.foot a:hover {
		color: var(--slate-900);
	}
</style>
