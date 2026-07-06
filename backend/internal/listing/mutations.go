package listing

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/race"
)

type createRequest struct {
	RaceID             string  `json:"race_id"`
	PriceCents         *int32  `json:"price_cents"`
	OriginalPriceCents *int32  `json:"original_price_cents"`
	Description        *string `json:"description"`
}

// create publishes a new listing for a race. Verified users only (the
// listing/chat gate carried over from #5); the race must exist, be published,
// and not be in the past. An asking price is allowed in every policy mode
// (informational outside platform_sale) and capped at face value per D2 when the
// original price is given.
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before creating a listing")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	raceID, err := uuid.Parse(req.RaceID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "race_id is not a valid id")
		return
	}

	rc, err := h.q.GetRaceByID(r.Context(), raceID)
	// A listing can only attach to a public race; drafts and unknown ids are 404.
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !race.IsPublic(rc.Status)) {
		httpx.Error(w, http.StatusNotFound, "not_found", "race not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load race")
		return
	}
	if racePast(rc.EventDate) {
		httpx.Error(w, http.StatusBadRequest, "race_past", "this race has already taken place")
		return
	}
	if err := validateListingPrice(req.PriceCents, req.OriginalPriceCents); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	row, err := h.q.CreateListing(r.Context(), sqlcgen.CreateListingParams{
		ID: ids.New(), RaceID: raceID, SellerID: caller.ID,
		PriceCents: req.PriceCents, Currency: "EUR", // EUR-only v1 (D11)
		OriginalPriceCents: req.OriginalPriceCents, Description: req.Description,
		ExpiresAt: rc.EventDate, // meaningful until race day; M4.2's job flips status
	})
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not create listing")
		return
	}
	httpx.JSON(w, http.StatusCreated, toSummary(row, caller.DisplayName))
}

type updateRequest struct {
	PriceCents         *int32  `json:"price_cents"`
	OriginalPriceCents *int32  `json:"original_price_cents"`
	Description        *string `json:"description"`
}

// update replaces the price, face value, and description of the owner's active
// listing (the edit form sends all three). Only the owner, and only while the
// listing is still active.
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	row, err := h.q.GetListingByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load listing")
		return
	}
	if row.Listing.SellerID != caller.ID {
		httpx.Error(w, http.StatusForbidden, "forbidden", "cannot modify another seller's listing")
		return
	}
	if row.Listing.Status != "active" {
		httpx.Error(w, http.StatusConflict, "not_active", "only an active listing can be edited")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	if err := validateListingPrice(req.PriceCents, req.OriginalPriceCents); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	updated, err := h.q.UpdateListing(r.Context(), sqlcgen.UpdateListingParams{
		ID: id, PriceCents: req.PriceCents,
		OriginalPriceCents: req.OriginalPriceCents, Description: req.Description,
		SellerID: caller.ID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// The read saw it active and owned, but it transitioned before this
		// guarded write - a concurrent change, not a stale client error.
		httpx.Error(w, http.StatusConflict, "not_active", "listing is no longer active")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not update listing")
		return
	}
	httpx.JSON(w, http.StatusOK, toSummary(updated, caller.DisplayName))
}

// cancel moves the owner's active listing to cancelled. The guarded transition
// (WHERE status = 'active') turns a non-active listing into a 409, not a silent
// no-op.
func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	row, err := h.q.GetListingByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load listing")
		return
	}
	if row.Listing.SellerID != caller.ID {
		httpx.Error(w, http.StatusForbidden, "forbidden", "cannot modify another seller's listing")
		return
	}

	updated, err := h.q.UpdateListingStatus(r.Context(), sqlcgen.UpdateListingStatusParams{
		ID: id, FromStatus: "active", ToStatus: "cancelled",
	})
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusConflict, "not_active", "only an active listing can be cancelled")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not cancel listing")
		return
	}
	httpx.JSON(w, http.StatusOK, toSummary(updated, caller.DisplayName))
}

type ownedListing struct {
	ID                 uuid.UUID `json:"id"`
	Status             string    `json:"status"`
	PriceCents         *int32    `json:"price_cents"`
	Currency           string    `json:"currency"`
	OriginalPriceCents *int32    `json:"original_price_cents"`
	Description        *string   `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	RaceName           string    `json:"race_name"`
	RaceSlug           string    `json:"race_slug"`
	EventDate          string    `json:"event_date"`
}

type ownedListResponse struct {
	Items []ownedListing `json:"items"`
}

// listMine returns the caller's own listings for the seller dashboard, newest
// first, across every status.
func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	rows, err := h.q.ListListingsBySeller(r.Context(), caller.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not list listings")
		return
	}
	items := make([]ownedListing, len(rows))
	for i, row := range rows {
		items[i] = ownedListing{
			ID: row.ID, Status: row.Status, PriceCents: row.PriceCents,
			Currency: row.Currency, OriginalPriceCents: row.OriginalPriceCents,
			Description: row.Description, CreatedAt: row.CreatedAt,
			RaceName: row.RaceName, RaceSlug: row.RaceSlug,
			EventDate: row.EventDate.Format("2006-01-02"),
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, ownedListResponse{Items: items})
}

func toSummary(l sqlcgen.Listing, sellerName string) Summary {
	return Summary{
		ID: l.ID, Status: l.Status, PriceCents: l.PriceCents, Currency: l.Currency,
		OriginalPriceCents: l.OriginalPriceCents, Description: l.Description,
		SellerName: sellerName, CreatedAt: l.CreatedAt,
	}
}

// racePast reports whether a race's day is already over (event_date is stored at
// 00:00 UTC of race day; a race happening today is not past).
func racePast(eventDate time.Time) bool {
	return eventDate.Before(time.Now().UTC().Truncate(24 * time.Hour))
}

// validateListingPrice enforces non-negative amounts and the D2 anti-scalping
// cap: when the seller supplies the original face value, the asking price may
// not exceed it. Both are optional - a listing can be contact-to-negotiate.
func validateListingPrice(price, original *int32) error {
	if price != nil && *price < 0 {
		return errors.New("price_cents must be 0 or more")
	}
	if original != nil && *original < 0 {
		return errors.New("original_price_cents must be 0 or more")
	}
	if price != nil && original != nil && *price > *original {
		return errors.New("asking price cannot exceed the original face value")
	}
	return nil
}
