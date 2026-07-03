import { defineConfig } from 'vitest/config';
import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
	plugins: [
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter(),
			// Root-relative paths (/races, not ./races). The default (relative: true)
			// breaks the i18n /es prefix (D17): a relative ./races prepended with /es
			// yields /es./races. We deploy adapter-node at the domain root, so
			// absolute paths are correct here.
			//
			// assets: in prod, emit absolute (origin-rooted) asset URLs so Vite's
			// __vitePreload deps are not resolved relative to the entry module URL,
			// which doubled the path (/_app/immutable/entry/_app/immutable/...) ->
			// 404 -> no CSS/JS (#66). Empty in dev; set per-origin via the
			// PUBLIC_ORIGIN build arg (the image is built for one origin).
			// Cast: paths.assets is typed as an absolute http(s) URL or ''; the env
			// var is a plain string (empty in dev, the origin in prod).
			paths: {
				relative: false,
				assets: (process.env.PUBLIC_ORIGIN ?? '') as '' | `http://${string}` | `https://${string}`
			},
			// Content-Security-Policy (#12 security pass). 'auto' makes SvelteKit
			// nonce its own inline hydration <script> and inline <style> elements
			// per request. The app loads nothing external (build-time SVG map, no
			// fonts/CDN/tiles, same-origin /api), so everything else is strict
			// 'self'. style-src-attr allows the inline style="" attributes the app
			// genuinely uses (app.html's body wrapper, the marquee's
			// --marquee-duration, RaceMap popover/dot positioning) - those can't be
			// nonced. No inline event handlers exist, so script-src carries no
			// 'unsafe-inline' - just 'self' (SvelteKit's JS bundles) plus the
			// per-request nonce SvelteKit adds for the inline hydration script.
			// frame-ancestors backs up Caddy's X-Frame-Options; base-uri is 'none'
			// (the app renders no <base> tag) to block base-tag injection.
			csp: {
				mode: 'auto',
				directives: {
					'default-src': ['self'],
					'script-src': ['self'],
					'style-src': ['self'],
					'style-src-attr': ['unsafe-inline'],
					'img-src': ['self', 'data:'],
					'object-src': ['none'],
					'base-uri': ['none'],
					'form-action': ['self'],
					'frame-ancestors': ['none']
				}
			}
		})
	],
	server: {
		// Same-origin API in dev: the browser hits /api/* here and Vite
		// forwards to the Go server - no CORS anywhere (docs/ARCHITECTURE.md).
		proxy: {
			'/api': 'http://localhost:8080'
		}
	},
	test: {
		expect: { requireAssertions: true },
		projects: [
			{
				extends: './vite.config.ts',
				test: {
					name: 'server',
					environment: 'node',
					include: ['src/**/*.{test,spec}.{js,ts}'],
					exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
				}
			}
		]
	}
});
