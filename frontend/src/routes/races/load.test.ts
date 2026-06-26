// The /races map is decorative: it is built from a second, unfiltered count
// fetch that exists only to colour the choropleth. That fetch must never be able
// to fail the whole browse page - if it errors, the grid still renders and the
// map is simply omitted (the +page.svelte gates the map on a non-empty
// countryCounts). apiGet is mocked so we can fail the count call in isolation.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { apiGet } from '$lib/api/server';
import { load } from './+page.server';

vi.mock('$lib/api/server', () => ({ apiGet: vi.fn() }));

const mockedApiGet = vi.mocked(apiGet);

// load() is typed via PageServerLoad to allow a void return; this route always
// returns the page payload, so narrow it for the assertions below.
type LoadData = Exclude<Awaited<ReturnType<typeof load>>, void>;

function event() {
	return {
		url: new URL('http://localhost/races'),
		fetch: vi.fn(),
		setHeaders: vi.fn(),
		locals: {}
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
	} as any;
}

const race = (over: Record<string, unknown> = {}) => ({
	id: '1',
	name: 'Paris Marathon',
	slug: 'paris-marathon',
	country: 'FR',
	city: 'Paris',
	...over
});

describe('races load', () => {
	beforeEach(() => mockedApiGet.mockReset());

	it('still renders the grid when the map count fetch fails', async () => {
		mockedApiGet
			.mockResolvedValueOnce({ items: [race()], next_cursor: null }) // grid
			.mockRejectedValueOnce(new Error('count endpoint down')); // map counts

		const data = (await load(event())) as LoadData;

		expect(data.races).toHaveLength(1);
		expect(data.countryCounts).toEqual({});
		expect(data.cities).toEqual([]);
	});

	it('builds country and city counts when the map fetch succeeds', async () => {
		mockedApiGet
			.mockResolvedValueOnce({ items: [race()], next_cursor: null }) // grid
			.mockResolvedValueOnce({
				items: [
					race(),
					race({ id: '2', name: 'Berlin', slug: 'berlin', country: 'DE', city: 'Berlin' })
				],
				next_cursor: null
			}); // map counts

		const data = (await load(event())) as LoadData;

		expect(data.countryCounts).toEqual({ FR: 1, DE: 1 });
		expect(data.cities).toHaveLength(2);
	});
});
