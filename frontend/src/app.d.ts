// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { SessionUser } from '$lib/api/types';
import type { Locale } from '$lib/i18n/locale';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: SessionUser | null;
			locale: Locale;
			// When set, the English page shows a dismissible "switch to Spanish"
			// banner (geo/Accept-Language suggestion; never an automatic redirect).
			suggestLocale: Locale | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
