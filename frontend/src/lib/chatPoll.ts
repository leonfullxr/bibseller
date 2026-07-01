// D13: 3-5s active polling. A hidden/backgrounded tab backs off ~10x (#96) -
// there's nothing new to show a tab nobody is looking at.
export const ACTIVE_POLL_MS = 4000;
export const HIDDEN_POLL_MS = 30000;

export function pollInterval(hidden: boolean): number {
	return hidden ? HIDDEN_POLL_MS : ACTIVE_POLL_MS;
}
