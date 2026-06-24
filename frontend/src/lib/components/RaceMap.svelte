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
	import {
		cityCoords,
		COUNTRY_VIEWBOX,
		EUROPE_VIEWBOX,
		fitViewBox,
		project
	} from '$lib/geo/cities';
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

	// The box keeps a constant aspect ratio (a constant on-page size): fitViewBox
	// pads a country's frame to that ratio so zooming never changes the page
	// height, and dot positions map linearly onto the box.
	const box = $derived(fitViewBox(zoomed ? COUNTRY_VIEWBOX[country] : EUROPE_VIEWBOX));
	const viewBox = $derived(box.join(' '));
	const unit = $derived(box[2] / 100); // marker sizes scale with the frame

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

	type Marker = {
		city: string;
		country: string;
		races: string[];
		x: number;
		y: number;
		left: number;
		top: number;
	};

	// City dots: only those we have coordinates for; when zoomed, only the focused
	// country's cities (the list below is already filtered to that country). left/top
	// are the dot's position as a percent of the box, for the hover popover.
	const markers = $derived<Marker[]>(
		cities
			.filter((c) => cityCoords(c.city) && (!zoomed || c.country === country))
			.map((c) => {
				const [x, y] = project(...cityCoords(c.city)!);
				return {
					...c,
					x,
					y,
					left: ((x - box[0]) / box[2]) * 100,
					top: ((y - box[1]) / box[3]) * 100
				};
			})
	);

	let hovered = $state<Marker | null>(null);
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
						onmouseenter={() => (hovered = m)}
						onmouseleave={() => (hovered = null)}
						onfocus={() => (hovered = m)}
						onblur={() => (hovered = null)}
					>
						<circle cx={m.x} cy={m.y} r={unit * 1.4} />
						<text x={m.x + unit * 2.4} y={m.y} font-size={unit * 2.8}>{m.city}</text>
					</a>
				{:else}
					<circle
						cx={m.x}
						cy={m.y}
						r={unit * 1.2}
						role="presentation"
						onmouseenter={() => (hovered = m)}
						onmouseleave={() => (hovered = null)}
					/>
				{/if}
			{/each}
		</svg>
		{#if hovered}
			<div class="popover" style="left:{hovered.left}%; top:{hovered.top}%">
				<span class="popover-city">{hovered.city}</span>
				{#each hovered.races as r (r)}
					<span class="race-box">{r}</span>
				{/each}
			</div>
		{/if}
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

	/* Fixed aspect-ratio box: the map is always the same on-page size, so picking a
	   country never grows the page or forces scrolling. */
	.map-wrap {
		position: relative;
		max-width: 46rem;
		margin-inline: auto;
		aspect-ratio: 105 / 87;
	}

	/* Base map and marker overlay both fill the box and share the same viewBox. */
	.map-wrap :global(svg) {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		display: block;
	}

	/* Overlay passes clicks through to the country links beneath; only the dots
	   themselves stay interactive (so countries are selectable from the map). */
	.markers {
		pointer-events: none;
	}

	.markers a,
	.markers circle {
		pointer-events: auto;
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
		pointer-events: none;
	}

	.map-wrap.zoomed :global(.has-races path) {
		fill: var(--emerald-600);
		stroke: white;
		stroke-width: 0.4;
		pointer-events: auto;
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

	/* Hover popover: one box per race, anchored above the city dot. */
	.popover {
		position: absolute;
		z-index: 2;
		transform: translate(-50%, calc(-100% - 0.5rem));
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		pointer-events: none;
	}

	.popover-city {
		text-align: center;
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--slate-900);
	}

	.race-box {
		white-space: nowrap;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.2rem 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-700);
		box-shadow: 0 1px 2px rgb(0 0 0 / 0.08);
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
