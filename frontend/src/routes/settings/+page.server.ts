import { fail } from '@sveltejs/kit';
import { apiFetch, apiGet } from '$lib/api/server';
import type { Actions, PageServerLoad } from './$types';

/**
 * TEMP until sessions (M3): the settings page edits the seeded demo user
 * (marta@example.com), whose id is fixed in backend/cmd/seed. Replace with
 * the signed-in user from `locals` once auth lands.
 */
const DEMO_USER_ID = '00000000-0000-7000-8000-000000000001';

interface UserProfile {
	id: string;
	display_name: string;
}

export const load: PageServerLoad = async ({ fetch }) => {
	const user = await apiGet<UserProfile>(`/api/v1/users/${DEMO_USER_ID}`, fetch);
	return { user };
};

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const displayName = String(data.get('display_name') ?? '').trim();

		// Server-side mirror of the HTML5 constraints — the browser attributes
		// are UX, not security; never trust the client.
		if (displayName.length < 2 || displayName.length > 50) {
			return fail(400, {
				value: displayName,
				error: 'Display name must be between 2 and 50 characters.'
			});
		}

		let res: Response;
		try {
			res = await apiFetch(`/api/v1/users/${DEMO_USER_ID}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ display_name: displayName })
			});
		} catch {
			return fail(502, { value: displayName, error: 'The API is unreachable.' });
		}

		if (!res.ok) {
			// The Go API's error envelope carries a human message; surface it.
			const body = (await res.json().catch(() => null)) as {
				error?: { message?: string };
			} | null;
			return fail(res.status >= 500 ? 502 : res.status, {
				value: displayName,
				error: body?.error?.message ?? 'The API rejected the update.'
			});
		}

		const user = (await res.json()) as UserProfile;
		return { success: true, value: user.display_name };
	}
};
