<script lang="ts">
	// Prototype (Option A): a dependency-free Europe choropleth. The base map is a
	// static SVG inlined at build time - no map library, no tile service.
	// Base map: flekschas/simple-world-map, CC BY-SA 3.0 (cropped to Europe).
	// Each race-country is wrapped in an <a> and highlighted server-side, so the
	// map renders, highlights, and navigates even without client JS.
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import europeMap from '$lib/assets/europe-map.svg?raw';

	let { counts }: { counts: Record<string, number> } = $props();
	const { t, link } = getI18n();
	const racesHref = $derived(link(resolve('/races')));

	// Wrap the <path> (single landmass) or <g> (with islands) for this country in
	// a link to its filtered list. CSS colours `.has-races path`.
	function linkify(svg: string, cc: string, n: number, href: string): string {
		const id = cc.toLowerCase();
		const label = t('races.mapCountry', { country: cc, n: String(n) });
		const open = `<a class="has-races" href="${href}?country=${cc}" aria-label="${label}">`;
		const pathRe = new RegExp(`<path\\b[^>]*\\bid="${id}"[^>]*/>`);
		const gRe = new RegExp(`<g\\b[^>]*\\bid="${id}"[^>]*>[\\s\\S]*?</g>`);
		if (pathRe.test(svg)) return svg.replace(pathRe, (m) => open + m + '</a>');
		if (gRe.test(svg)) return svg.replace(gRe, (m) => open + m + '</a>');
		return svg;
	}

	const svgHtml = $derived(
		Object.entries(counts).reduce((svg, [cc, n]) => linkify(svg, cc, n, racesHref), europeMap)
	);
</script>

<section class="race-map" aria-label={t('races.mapHeading')}>
	<div class="map-wrap">
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html svgHtml}
	</div>
	<p class="map-hint">{t('races.mapHint')}</p>
</section>

<style>
	.race-map {
		margin-top: 1.25rem;
	}

	.map-wrap {
		max-width: 30rem;
		margin-inline: auto;
	}

	.map-wrap :global(svg) {
		width: 100%;
		height: auto;
		display: block;
	}

	/* Base countries (no races). */
	.map-wrap :global(path) {
		fill: var(--slate-200);
		stroke: white;
		stroke-width: 0.4;
	}

	/* Race-countries: fill the linked country (and its island paths). */
	.map-wrap :global(.has-races) {
		cursor: pointer;
		outline: none;
	}

	.map-wrap :global(.has-races path) {
		fill: var(--emerald-600);
	}

	.map-wrap :global(.has-races:hover path),
	.map-wrap :global(.has-races:focus-visible path) {
		fill: var(--emerald-700);
	}

	.map-hint {
		margin-top: 0.75rem;
		text-align: center;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}
</style>
