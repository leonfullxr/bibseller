<script lang="ts">
	// Dependency-free Europe choropleth. The base map is a static SVG inlined at
	// build time - no map library, no tile service.
	// Base map: flekschas/simple-world-map, CC BY-SA 3.0.
	// Each race-country is wrapped in an <a> and highlighted server-side, so the
	// map renders, highlights, and navigates even without client JS. The viewBox
	// is set per request: the whole of Europe by default, or one country's frame
	// when that country filter is active (the URL filter is the zoom state).
	import { resolve } from '$app/paths';
	import { getI18n } from '$lib/i18n';
	import { cityCoords, COUNTRY_VIEWBOX, EUROPE_VIEWBOX, project } from '$lib/geo/cities';
	import europeMap from '$lib/assets/europe-map.svg?raw';

	let {
		counts,
		cities,
		country
	}: {
		counts: Record<string, number>;
		cities: { city: string; country: string; races: string[] }[];
		country: string;
	} = $props();
	const { t, link } = getI18n();
	const racesHref = $derived(link(resolve('/races')));

	const zoomed = $derived(Boolean(country && COUNTRY_VIEWBOX[country]));
	const viewBox = $derived(zoomed ? COUNTRY_VIEWBOX[country] : EUROPE_VIEWBOX);
	// Marker sizes scale with the viewBox so dots/labels stay a constant on-screen
	// size whether we show all of Europe or a single zoomed-in country.
	const unit = $derived(Number(viewBox.split(' ')[2]) / 100);

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

	// When zoomed, highlight only the focused country; CSS (.zoomed) hides the
	// rest so the map shows that country alone.
	const baseHtml = $derived(
		(zoomed
			? linkify(europeMap, country, counts[country] ?? 0, racesHref)
			: Object.entries(counts).reduce((svg, [cc, n]) => linkify(svg, cc, n, racesHref), europeMap)
		).replace(/viewBox="[^"]*"/, `viewBox="${viewBox}"`)
	);

	// City dots: only those we have coordinates for; when zoomed, only the focused
	// country's cities (the list below is already filtered to that country). The
	// hover title lists the city's races (CONTEXT: a city can hold several).
	const markers = $derived(
		cities
			.filter((c) => cityCoords(c.city) && (!zoomed || c.country === country))
			.map((c) => {
				const [x, y] = project(...cityCoords(c.city)!);
				return { ...c, x, y, title: [c.city, ...c.races].join('\n') };
			})
	);
</script>

<section class="race-map" aria-label={t('races.mapHeading')}>
	<div class="map-wrap" class:zoomed>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html baseHtml}
		<svg class="markers" {viewBox} role="presentation" aria-hidden="true">
			{#each markers as m (m.country + m.city)}
				{#if zoomed}
					<a
						href="{racesHref}?country={m.country}&q={encodeURIComponent(m.city)}"
						aria-label={t('races.mapCity', { city: m.city, n: String(m.races.length) })}
					>
						<title>{m.title}</title>
						<circle cx={m.x} cy={m.y} r={unit * 1.4} />
						<text x={m.x + unit * 2.4} y={m.y} font-size={unit * 2.8}>{m.city}</text>
					</a>
				{:else}
					<circle cx={m.x} cy={m.y} r={unit * 1.2}>
						<title>{m.title}</title>
					</circle>
				{/if}
			{/each}
		</svg>
	</div>
	{#if zoomed}
		<p class="map-hint">
			<a href={racesHref}>{t('races.mapBack')}</a>
		</p>
	{:else}
		<p class="map-hint">{t('races.mapHint')}</p>
	{/if}
</section>

<style>
	.race-map {
		margin-top: 1.25rem;
	}

	.map-wrap {
		position: relative;
		max-width: 46rem;
		margin-inline: auto;
	}

	.map-wrap :global(svg) {
		width: 100%;
		height: auto;
		display: block;
	}

	/* Marker overlay sits exactly on top of the base map (same viewBox). */
	.markers {
		position: absolute;
		inset: 0;
		height: 100%;
	}

	/* Base countries (no races). */
	.map-wrap :global(path) {
		fill: var(--slate-200);
		stroke: white;
		stroke-width: 0.4;
	}

	/* Zoomed into one country: hide every other country so it shows alone. */
	.map-wrap.zoomed :global(path) {
		fill: transparent;
		stroke: none;
	}

	.map-wrap.zoomed :global(.has-races path) {
		fill: var(--emerald-600);
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

	/* City dots. */
	.markers circle {
		fill: var(--emerald-700);
		stroke: white;
		stroke-width: 0.4;
	}

	.markers a {
		cursor: pointer;
	}

	.markers a:hover circle,
	.markers a:focus-visible circle {
		fill: var(--slate-900);
	}

	.markers text {
		fill: var(--slate-700);
		stroke: white;
		stroke-width: 0.8;
		paint-order: stroke;
		dominant-baseline: middle;
		font-weight: 600;
	}

	.map-hint {
		margin-top: 0.75rem;
		text-align: center;
		font-size: 0.8125rem;
		line-height: 1.25rem;
		color: var(--slate-500);
	}

	.map-hint a {
		color: var(--emerald-700);
		text-decoration: underline;
	}
</style>
