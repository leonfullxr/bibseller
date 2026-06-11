<script lang="ts">
	import { resolve } from '$app/paths';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex min-h-screen flex-col bg-slate-50 text-slate-900">
	<header class="border-b border-slate-200 bg-white">
		<div class="mx-auto flex h-14 w-full max-w-5xl items-center justify-between px-4">
			<a href={resolve('/')} class="text-lg font-bold tracking-tight">
				bib<span class="text-emerald-600">seller</span>
			</a>
			<nav class="text-sm text-slate-400">
				<span title="Coming with M2">Races — soon</span>
			</nav>
		</div>
	</header>

	<main class="mx-auto w-full max-w-5xl flex-1 px-4 py-10">
		{@render children()}
	</main>

	<footer class="border-t border-slate-200 bg-white">
		<div
			class="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-4 text-xs text-slate-500"
		>
			<span>Zero commission, EU-wide. Non-profit by design.</span>
			<a href="https://github.com/leonfullxr/bibseller" class="hover:text-slate-900" rel="external"
				>GitHub</a
			>
		</div>
	</footer>
</div>
