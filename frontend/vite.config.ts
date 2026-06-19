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
			paths: { relative: false }
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
