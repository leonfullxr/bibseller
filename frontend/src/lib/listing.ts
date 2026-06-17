/** Listing form helpers shared by the create (/sell) and edit flows. */

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
): { ok: true; value: ParsedPrice } | { ok: false; error: string } {
	const price = toCents(priceRaw);
	const original = toCents(originalRaw);
	if (price === 'invalid' || original === 'invalid') {
		return { ok: false, error: 'Enter a valid amount, e.g. 45 or 45.00.' };
	}
	if (price != null && original != null && price > original) {
		return { ok: false, error: 'Asking price cannot exceed the original face value.' };
	}
	return { ok: true, value: { priceCents: price, originalCents: original } };
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
