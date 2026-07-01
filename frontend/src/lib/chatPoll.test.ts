import { describe, expect, it } from 'vitest';
import { ACTIVE_POLL_MS, HIDDEN_POLL_MS, pollInterval } from './chatPoll';

describe('pollInterval', () => {
	it('polls fast while the tab is visible', () => {
		expect(pollInterval(false)).toBe(ACTIVE_POLL_MS);
	});

	it('backs off while the tab is hidden', () => {
		expect(pollInterval(true)).toBe(HIDDEN_POLL_MS);
		expect(HIDDEN_POLL_MS).toBeGreaterThanOrEqual(ACTIVE_POLL_MS * 5);
	});
});
