<script lang="ts">
	// Small inline-SVG icon set for the landing page (no icon dependency, D14).
	// Stroke uses currentColor and the box is 1em, so size/colour come from CSS.
	let { name }: { name: string } = $props();
</script>

<svg
	class="icon"
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.75"
	stroke-linecap="round"
	stroke-linejoin="round"
	aria-hidden="true"
>
	{#if name === 'list'}
		<rect x="5" y="4" width="14" height="17" rx="2" />
		<path d="M9 4h6v3H9z" />
		<path d="M8.5 11h7M8.5 15h4.5" />
	{:else if name === 'search'}
		<circle cx="11" cy="11" r="6" />
		<path d="M20 20l-4.3-4.3" />
	{:else if name === 'chat'}
		<rect x="4" y="5" width="16" height="11" rx="2.5" />
		<path d="M9 16v3.2L12.5 16" />
	{:else if name === 'transfer'}
		<path d="M5 9h11l-3-3M19 15H8l3 3" />
	{:else if name === 'check'}
		<circle cx="12" cy="12" r="9" />
		<path d="M8 12.5l2.8 2.8L16 9.5" />
	{:else if name === 'handover'}
		<path d="M12 3l8 4.2v9.6L12 21l-8-4.2V7.2z" />
		<path d="M4 7.2l8 4.2 8-4.2M12 11.4V21" />
	{:else if name === 'medal'}
		<circle cx="12" cy="14" r="6" />
		<path d="M8.5 8L6 3.5M15.5 8L18 3.5" />
		<path d="M10 14l1.5 1.5 3-3" />
	{:else if name === 'person'}
		<circle cx="12" cy="8" r="3.5" />
		<path d="M5.5 20a6.5 6.5 0 0 1 13 0" />
	{:else if name === 'arrow'}
		<path d="M4 12h13" />
		<path d="M13 7l5 5-5 5" />
	{/if}
</svg>

<style>
	.icon {
		width: 1em;
		height: 1em;
		display: block;
	}
</style>
