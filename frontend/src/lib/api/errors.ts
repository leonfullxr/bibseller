// Maps a Go API error-envelope `code` to a translatable message key. Server
// actions read the stable code (never the English message) and translate via
// t(apiErrorKey(code)). Any unknown or missing code falls back to
// apiError.unknown, so no English API string is ever fed to the translator (#49).
// Mirrors the `key in en` guard used by sportLabel.
import { en, type MessageKey } from '$lib/i18n/en';

export function apiErrorKey(code: string | null | undefined): MessageKey {
	const key = `apiError.${code}` as MessageKey;
	return code && key in en ? key : 'apiError.unknown';
}
