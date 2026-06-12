<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
</script>

<div class="py-20 text-center">
	<p class="text-6xl font-extrabold text-slate-300">{page.status}</p>
	<h1 class="mt-4 text-xl font-semibold">
		{page.status === 404 ? "That page doesn't exist." : 'Something went wrong.'}
	</h1>
	<p class="mt-2 text-slate-600">{page.error?.message}</p>
	<a
		href={resolve('/')}
		class="mt-6 inline-block rounded-md bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-700"
	>
		Back home
	</a>
</div>
