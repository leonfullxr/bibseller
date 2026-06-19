<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { createTranslator, defaultLocale, getI18n } from '$lib/i18n';

	// The error boundary can render above the root layout (e.g. a failure in
	// hooks), where the i18n context is not set - fall back to English so the
	// error page itself never crashes on a missing context.
	const i18n = getI18n();
	const t = i18n?.t ?? createTranslator(defaultLocale);
	const homeHref = i18n ? i18n.link(resolve('/')) : resolve('/');
</script>

<div class="error">
	<p class="status">{page.status}</p>
	<h1>
		{page.status === 404 ? t('error.notFound') : t('error.generic')}
	</h1>
	<p class="msg">{page.error?.message}</p>
	<a href={homeHref} class="home">{t('error.backHome')}</a>
</div>

<style>
	.error {
		padding-block: 5rem;
		text-align: center;
	}

	.status {
		font-size: 3.75rem;
		line-height: 1;
		font-weight: 800;
		color: var(--slate-300);
	}

	h1 {
		margin-top: 1rem;
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.msg {
		margin-top: 0.5rem;
		color: var(--slate-600);
	}

	.home {
		margin-top: 1.5rem;
		display: inline-block;
		border-radius: 0.375rem;
		background: var(--slate-900);
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		font-weight: 600;
		color: white;
	}

	.home:hover {
		background: var(--slate-700);
	}
</style>
