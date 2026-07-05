// Package user owns user profile updates.
//
// PATCH /users/{id} is gated to the signed-in user - you may rename
// yourself, no one else (401/403 per docs/ARCHITECTURE.md).
package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

// allowedLocales is the v1 language set (docs/CONTEXT.md -> D4: en first, es next).
var allowedLocales = map[string]bool{"en": true, "es": true}

// allowedCountries mirrors the catalog's country filter (ISO 3166-1 alpha-2).
var allowedCountries = map[string]bool{
	"AT": true, "BE": true, "DE": true, "ES": true, "FR": true,
	"IT": true, "NL": true, "PL": true, "PT": true,
}

type Handler struct {
	q *sqlcgen.Queries
}

func Routes(q *sqlcgen.Queries) func(*http.ServeMux) {
	h := &Handler{q: q}
	return func(mux *http.ServeMux) {
		mux.HandleFunc("PATCH /users/{id}", h.updateProfile)
	}
}

// Profile deliberately excludes email and other PII.
type Profile struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
}

// updateRequest replaces the editable profile fields. A signed-in user saves
// their whole profile at once (the /settings form sends all three), so this is
// a replace, not a partial patch.
type updateRequest struct {
	DisplayName string  `json:"display_name"`
	Locale      string  `json:"locale"`
	Country     *string `json:"country"` // null or "" clears it (stored NULL)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	if caller.ID != id {
		httpx.Error(w, http.StatusForbidden, "forbidden", "cannot modify another user")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}

	name, err := auth.ValidateDisplayName(req.DisplayName)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	if !allowedLocales[req.Locale] {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "unsupported locale")
		return
	}
	country, err := normalizeCountry(req.Country)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	row, err := h.q.UpdateUserProfile(r.Context(), sqlcgen.UpdateUserProfileParams{
		ID: id, DisplayName: name, Locale: req.Locale, Country: country,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not update user")
		return
	}

	httpx.JSON(w, http.StatusOK, Profile{ID: row.ID, DisplayName: row.DisplayName})
}

// normalizeCountry upper-cases and validates an optional country code. A nil or
// empty value clears it (stored NULL); anything outside the allowlist is an error.
func normalizeCountry(in *string) (*string, error) {
	if in == nil {
		return nil, nil
	}
	c := strings.ToUpper(strings.TrimSpace(*in))
	if c == "" {
		return nil, nil
	}
	if !allowedCountries[c] {
		return nil, errors.New("unsupported country")
	}
	return &c, nil
}
