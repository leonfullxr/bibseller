/** Listing form helpers shared by the create (/sell) and edit flows. */
import type { MessageKey } from '$lib/i18n/en';

export interface ParsedPrice {
	priceCents: number | null;
	originalCents: number | null;
}

/**
 * Parses the optional euro amounts from a listing form into integer cents,
 * mirroring the API's rules: amounts are non-negative and the asking price may
 * not exceed the face value (D2). The API stays the authority; this is UX.
 */
export function parseListingPrice(
	priceRaw: string,
	originalRaw: string
): { ok: true; value: ParsedPrice } | { ok: false; key: MessageKey } {
	const price = toCents(priceRaw);
	const original = toCents(originalRaw);
	if (price === 'invalid' || original === 'invalid') {
		return { ok: false, key: 'formError.invalidAmount' };
	}
	if (price != null && original != null && price > original) {
		return { ok: false, key: 'formError.priceExceedsFace' };
	}
	return { ok: true, value: { priceCents: price, originalCents: original } };
}

/**
 * Echoes the listing form's entered values back into a `fail()` payload so an
 * invalid submit re-renders with the user's input preserved. Shared by the
 * create (/sell) and edit actions.
 */
export function listingFormSnapshot(form: FormData) {
	return {
		price: String(form.get('price') ?? ''),
		original_price: String(form.get('original_price') ?? ''),
		description: String(form.get('description') ?? '')
	};
}

// toCents returns null for empty (the field is optional), 'invalid' for a
// non-numeric or negative amount, or the rounded integer cents otherwise.
function toCents(raw: string): number | null | 'invalid' {
	const s = raw.trim();
	if (s === '') return null;
	const n = Number(s);
	if (!Number.isFinite(n) || n < 0) return 'invalid';
	return Math.round(n * 100);
}
