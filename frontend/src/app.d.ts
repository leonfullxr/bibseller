// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { SessionUser } from '$lib/api/types';
import type { Locale } from '$lib/i18n/locale';
import type { MessageKey } from '$lib/i18n/en';

declare global {
	namespace App {
		interface Error {
			message: string;
			// Optional i18n key for the detail line. When set, +error.svelte renders
			// t(key) instead of the English message (#49); the message stays as the
			// fallback for boundaries rendered without an i18n context.
			key?: MessageKey;
		}
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
