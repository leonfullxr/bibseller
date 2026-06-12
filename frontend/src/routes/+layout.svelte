<script lang="ts">
	import { resolve } from '$app/paths';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="shell">
	<header>
		<div class="bar">
			<a href={resolve('/')} class="brand">bib<span>seller</span></a>
			<nav>
				<a href={resolve('/races')}>Races</a>
			</nav>
		</div>
	</header>

	<main>
		{@render children()}
	</main>

	<footer>
		<div class="bar foot">
			<span>Zero commission, EU-wide. Non-profit by design.</span>
			<a href="https://github.com/leonfullxr/bibseller" rel="external">GitHub</a>
		</div>
	</footer>
</div>

<style>
	.shell {
		display: flex;
		min-height: 100vh;
		flex-direction: column;
		background: var(--slate-50);
		color: var(--slate-900);
	}

	header {
		border-bottom: 1px solid var(--slate-200);
		background: white;
	}

	.bar {
		margin-inline: auto;
		display: flex;
		width: 100%;
		max-width: 64rem;
		align-items: center;
		justify-content: space-between;
		padding-inline: 1rem;
	}

	header .bar {
		height: 3.5rem;
	}

	.brand {
		font-size: 1.125rem;
		line-height: 1.75rem;
		font-weight: 700;
		letter-spacing: -0.025em;
	}

	.brand span {
		color: var(--emerald-600);
	}

	nav {
		font-size: 0.875rem;
		line-height: 1.25rem;
	}

	nav a {
		font-weight: 500;
		color: var(--slate-600);
	}

	nav a:hover {
		color: var(--slate-900);
	}

	main {
		margin-inline: auto;
		width: 100%;
		max-width: 64rem;
		flex: 1;
		padding: 2.5rem 1rem;
	}

	footer {
		border-top: 1px solid var(--slate-200);
		background: white;
	}

	.foot {
		padding-block: 1rem;
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}

	.foot a:hover {
		color: var(--slate-900);
	}
</style>
