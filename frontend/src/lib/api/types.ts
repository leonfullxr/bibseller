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
	seller_id: string;
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
	email_verified: boolean;
	locale: string;
	country: string | null;
}

export interface SessionResponse {
	token: string;
	expires_at: string;
	user: SessionUser;
}

/** A chat message, as returned by the chat message endpoints (M5.1). */
export interface ChatMessage {
	id: string;
	sender_id: string;
	body: string;
	created_at: string;
}

/** An inbox thread row, as returned by GET /threads. */
export interface ChatThreadSummary {
	id: string;
	listing_id: string;
	race_name: string;
	race_slug: string;
	role: 'buyer' | 'seller'; // the caller's role in this thread
	other_party: string; // display name of the other participant
	last_message_at: string | null;
	unread_count: number;
}

/** A seller's own listing, as returned by GET /me/listings. */
export interface OwnedListing {
	id: string;
	status: string;
	price_cents: number | null;
	currency: string;
	original_price_cents: number | null;
	description: string | null;
	created_at: string;
	race_name: string;
	race_slug: string;
	event_date: string;
}
