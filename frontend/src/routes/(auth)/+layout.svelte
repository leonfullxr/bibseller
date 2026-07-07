<script lang="ts">
	// Shared shell for the auth pages (login, register, forgot, reset): the
	// centered panel plus the form styling every page reuses - the same
	// route-group dedup the (legal) layout does for legal prose (#121).
	let { children } = $props();
</script>

<section class="panel">
	{@render children()}
</section>

<style>
	.panel {
		margin-inline: auto;
		max-width: 24rem;
		border: 2px solid var(--ink);
		box-shadow: var(--shadow-hard);
	}

	/* Form styling shared by every auth page. :global so the child page's own
	   h1/form/label (a separate component) inherit it from this layout. */
	.panel :global(h1) {
		font-size: 1.25rem;
		line-height: 1.75rem;
		font-weight: 700;
	}

	.panel :global(form) {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.panel :global(label) {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		color: var(--slate-600);
	}

	.panel :global(.alert) {
		margin-top: 0.75rem;
	}

	.panel :global(.alert.ok) {
		border-color: var(--emerald-200);
		background: var(--emerald-50);
		color: var(--emerald-900);
	}

	.panel :global(button) {
		margin-top: 1rem;
	}

	.panel :global(.alt) {
		margin-top: 1rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: var(--slate-600);
	}

	.panel :global(.alt a) {
		color: var(--brand-700);
		text-decoration: underline;
	}
</style>
