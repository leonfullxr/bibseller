<script lang="ts">
	// Shared shell for the legal pages (M7-lite, #11): a draft banner plus the
	// prose typography every legal page reuses. Copy is DRAFT pending counsel
	// review (D10); the banner is the user-facing marker.
	import { getI18n } from '$lib/i18n';

	let { children } = $props();
	const { t } = getI18n();
</script>

<section class="legal">
	<p class="draft" role="note">{t('legal.draftNotice')}</p>
	{@render children()}
</section>

<style>
	.legal {
		margin-inline: auto;
		max-width: 44rem;
	}

	.draft {
		margin-bottom: 1.5rem;
		border-radius: 0.375rem;
		border: 1px solid var(--amber-300);
		background: var(--amber-50);
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--amber-900);
	}

	/* Prose typography shared by every legal page. :global so the child page's
	   own h1/h2/p/dl (a separate component) inherit it from this layout. */
	.legal :global(h1) {
		font-size: 1.5rem;
		line-height: 2rem;
		font-weight: 700;
	}

	.legal :global(h2) {
		margin-top: 1.75rem;
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 600;
	}

	.legal :global(p) {
		margin-top: 0.75rem;
		font-size: 0.9375rem;
		line-height: 1.6;
		color: var(--slate-700);
	}

	.legal :global(dt) {
		margin-top: 1rem;
		font-weight: 600;
		color: var(--slate-900);
	}

	.legal :global(dd) {
		margin-top: 0.25rem;
		margin-inline-start: 0;
		font-size: 0.9375rem;
		line-height: 1.6;
		color: var(--slate-700);
	}
</style>
