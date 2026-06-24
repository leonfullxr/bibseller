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
		cities: { city: string; country: string; races: { name: string; slug: string }[] }[];
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
		races: { name: string; slug: string }[];
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

	const RADIUS = 4.6; // rem: radius of the race ring around the hovered city
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
					/>
				{/if}
			{/each}
		</svg>
		{#if hovered}
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				class="popover"
				style="left:{hovered.left}%; top:{hovered.top}%"
				onmouseleave={() => (hovered = null)}
			>
				<span class="popover-city">{hovered.city}</span>
				{#each hovered.races as r, i (r.slug)}
					{@const angle = (i / hovered.races.length) * Math.PI * 2 - Math.PI / 2}
					<span
						class="race-pos"
						style="transform: translate(-50%, -50%) translate({(Math.cos(angle) * RADIUS).toFixed(
							2
						)}rem, {(Math.sin(angle) * RADIUS).toFixed(2)}rem)"
					>
						<a
							class="race-box"
							href={link(resolve('/races/[slug]', { slug: r.slug }))}
							style="animation-delay: {i * 45}ms">{r.name}</a
						>
					</span>
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

	/* Hover popover: the city sits at the centre with its races on a ring around
	   it. The container spans the ring so the popover stays open while the cursor
	   travels from the dot out to the cards. */
	.popover {
		position: absolute;
		z-index: 2;
		width: 15rem;
		height: 15rem;
		transform: translate(-50%, -50%);
		pointer-events: auto;
	}

	.popover-city {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		border-radius: 999px;
		background: var(--emerald-600);
		padding: 0.25rem 0.7rem;
		font-size: 0.7rem;
		font-weight: 700;
		white-space: nowrap;
		color: white;
		box-shadow: 0 2px 6px rgb(0 0 0 / 0.25);
	}

	.race-pos {
		position: absolute;
		left: 50%;
		top: 50%;
	}

	/* Each race is a narrow, tall card; the name wraps into a few lines. */
	.race-box {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 5rem;
		min-height: 3.4rem;
		padding: 0.4rem;
		text-align: center;
		text-wrap: balance;
		text-decoration: none;
		cursor: pointer;
		border-radius: 0.5rem;
		border: 1px solid var(--slate-300);
		background: white;
		font-size: 0.7rem;
		line-height: 1.05rem;
		font-weight: 600;
		color: var(--slate-700);
		box-shadow: 0 2px 6px rgb(0 0 0 / 0.12);
		animation: race-pop 0.22s ease both;
		transition:
			transform 0.15s ease,
			box-shadow 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
	}

	.race-box:hover {
		transform: scale(1.22);
		z-index: 3;
		border-color: var(--emerald-600);
		color: var(--slate-900);
		box-shadow: 0 10px 22px rgb(0 0 0 / 0.22);
	}

	@keyframes race-pop {
		from {
			opacity: 0;
			transform: scale(0.7);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.race-box {
			animation: none;
		}
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
