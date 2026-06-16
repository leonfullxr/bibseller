/** Mirrors the Go API DTOs (backend/internal/race, backend/internal/listing). */

export type TransferPolicy = 'platform_sale' | 'official_only' | 'connect_only' | 'unknown';

export interface RaceSummary {
	id: string;
	slug: string;
	name: string;
	series: string | null;
	sport: string;
	distance: string | null;
	event_date: string; // YYYY-MM-DD
	city: string;
	country: string;
	transfer_policy: TransferPolicy;
	active_listings: number;
}

export interface RaceDetail extends RaceSummary {
	website_url: string | null;
	official_transfer_url: string | null;
	policy_notes: string | null;
	policy_verified_at: string | null;
}

export interface ListingSummary {
	id: string;
	status: string;
	price_cents: number | null;
	currency: string;
	original_price_cents: number | null;
	description: string | null;
	seller_name: string;
	created_at: string;
}

export interface ListingDetail extends ListingSummary {
	race: {
		slug: string;
		name: string;
		distance: string | null;
		event_date: string;
		city: string;
		country: string;
		transfer_policy: TransferPolicy;
		official_transfer_url: string | null;
	};
}

export interface Page<T> {
	items: T[];
	next_cursor: string | null;
}

/** The signed-in account, as returned by GET /auth/me and the session response. */
export interface SessionUser {
	id: string;
	email: string;
	display_name: string;
}

export interface SessionResponse {
	token: string;
	expires_at: string;
	user: SessionUser;
}
