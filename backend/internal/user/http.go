// Package user owns user profile reads and updates.
//
// GET /users/{id} is public: it returns only id + display_name, the same
// non-PII the catalog already exposes as a listing's seller_name. PATCH is
// gated to the signed-in user - you may rename yourself, no one else
// (401/403 per docs/ARCHITECTURE.md).
package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

const (
	minNameLen = 2
	maxNameLen = 50
)

type Handler struct {
	q *sqlcgen.Queries
}

func Routes(q *sqlcgen.Queries) func(*http.ServeMux) {
	h := &Handler{q: q}
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /users/{id}", h.get)
		mux.HandleFunc("PATCH /users/{id}", h.updateDisplayName)
	}
}

// Profile deliberately excludes email and other PII: the endpoint is
// unauthenticated until M3.
type Profile struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}

	row, err := h.q.GetUserByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load user")
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, Profile{ID: row.ID, DisplayName: row.DisplayName})
}

type updateRequest struct {
	DisplayName string `json:"display_name"`
}

func (h *Handler) updateDisplayName(w http.ResponseWriter, r *http.Request) {
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

	name := strings.TrimSpace(req.DisplayName)
	if n := utf8.RuneCountInString(name); n < minNameLen || n > maxNameLen {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter",
			fmt.Sprintf("display_name must be %d..%d characters", minNameLen, maxNameLen))
		return
	}

	row, err := h.q.UpdateUserDisplayName(r.Context(), sqlcgen.UpdateUserDisplayNameParams{
		ID: id, DisplayName: name,
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
