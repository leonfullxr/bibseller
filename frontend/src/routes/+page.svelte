<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const modes = [
		{
			name: 'Platform sale',
			desc: 'The race allows resale: list, chat, and pay securely through the platform.'
		},
		{
			name: 'Official process',
			desc: 'The race runs its own name change: we connect you and link the official procedure.'
		},
		{
			name: 'Connect only',
			desc: 'Restricted or unverified races: we provide the chat, the rest stays between you two.'
		}
	];
</script>

<svelte:head>
	<title>Bibseller — race bibs find new runners</title>
	<meta
		name="description"
		content="Non-profit, EU-wide marketplace connecting runners who can't start with runners who missed registration."
	/>
</svelte:head>

<section class="py-8 text-center">
	<h1 class="text-4xl font-extrabold tracking-tight sm:text-5xl">
		Race bibs find <span class="text-emerald-600">new runners</span>
	</h1>
	<p class="mx-auto mt-4 max-w-2xl text-lg text-slate-600">
		Injured? Plans changed? Missed registration? A zero-commission marketplace that connects sellers
		and buyers of race bibs — always within each race's own rules.
	</p>

	<div class="mt-6 flex items-center justify-center gap-2 text-sm">
		{#if data.apiStatus === 'ok'}
			<span class="inline-block h-2 w-2 rounded-full bg-emerald-500"></span>
			<span class="text-slate-600">API connected</span>
		{:else}
			<span class="inline-block h-2 w-2 rounded-full bg-amber-500"></span>
			<span class="text-slate-600">
				API unreachable — run <code class="rounded bg-slate-200 px-1 py-0.5">make dev</code>
			</span>
		{/if}
	</div>
</section>

<section class="grid gap-4 py-8 sm:grid-cols-3">
	{#each modes as mode (mode.name)}
		<div class="rounded-lg border border-slate-200 bg-white p-5">
			<h2 class="font-semibold">{mode.name}</h2>
			<p class="mt-2 text-sm text-slate-600">{mode.desc}</p>
		</div>
	{/each}
</section>

<p class="py-4 text-center text-sm text-slate-400">
	Under construction — follow the
	<a
		href="https://github.com/leonfullxr/bibseller/issues/13"
		class="underline hover:text-slate-600"
		rel="external">roadmap</a
	>.
</p>
