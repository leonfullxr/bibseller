// The /races map is built from a second, decorative fetch to the server-side
// aggregate endpoint (/races/map-counts). It must never fail the browse page: if
// it errors, the grid still renders and the map is simply omitted (+page.svelte
// gates the map on a non-empty countryCounts). apiGet is mocked so we can drive
// each call independently.
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

describe('races load', () => {
	beforeEach(() => mockedApiGet.mockReset());

	it('builds map data from the map-counts endpoint', async () => {
		mockedApiGet
			.mockResolvedValueOnce({ items: [{ id: '1' }], next_cursor: null }) // grid
			.mockResolvedValueOnce({
				countries: { FR: 3, DE: 1 },
				cities: [{ city: 'Paris', country: 'FR', count: 3, races: [{ name: 'X', slug: 'x' }] }]
			}); // map-counts

		const data = (await load(event())) as LoadData;

		expect(data.countryCounts).toEqual({ FR: 3, DE: 1 });
		expect(data.cities).toHaveLength(1);
		expect(data.cities[0].count).toBe(3);
	});

	it('still renders the grid when the map-counts fetch fails', async () => {
		mockedApiGet
			.mockResolvedValueOnce({ items: [{ id: '1' }], next_cursor: null }) // grid
			.mockRejectedValueOnce(new Error('map-counts down')); // map-counts

		const data = (await load(event())) as LoadData;

		expect(data.races).toHaveLength(1);
		expect(data.countryCounts).toEqual({});
		expect(data.cities).toEqual([]);
	});
});
