import type { SubmitFunction } from '@sveltejs/kit';

/**
 * Pending state for use:enhance forms:
 *   const { busy, submit } = pendingForm();
 *   <form method="POST" use:enhance={submit}> ... <button disabled={busy.value}>
 */
// ponytail: one flag per pendingForm() call; forms sharing an instance disable together.
export function pendingForm() {
	const busy = $state({ value: false });
	const submit: SubmitFunction = () => {
		busy.value = true;
		return async ({ update }) => {
			await update();
			busy.value = false;
		};
	};
	return { busy, submit };
}
