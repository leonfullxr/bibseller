package auth

import (
	"encoding/json"
	"net/http"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// changePassword updates the signed-in user's password. It re-checks the
// current password first - a stolen-but-live session must not be enough to
// lock the real owner out - then keeps the caller signed in while revoking
// every other session, so changing the password signs out all other devices.
func (h *Handler) changePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	if err := validatePassword(req.NewPassword); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	match, err := verifyPassword(req.CurrentPassword, user.PasswordHash)
	if err != nil || !match {
		httpx.Error(w, http.StatusUnauthorized, "invalid_credentials", "current password is incorrect")
		return
	}

	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not change password")
		return
	}
	if err := h.q.UpdateUserPassword(r.Context(), sqlcgen.UpdateUserPasswordParams{
		ID: user.ID, PasswordHash: hash,
	}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not change password")
		return
	}

	// Keep this session, revoke the rest. The request authenticated, so a token
	// is present; if it somehow is not, fail safe by revoking all (never fewer).
	if token, ok := requestToken(r); ok {
		_ = h.q.DeleteSessionsForUserExcept(r.Context(), sqlcgen.DeleteSessionsForUserExceptParams{
			UserID: user.ID, TokenHash: hashToken(token),
		})
	} else {
		_ = h.q.DeleteAllSessionsForUser(r.Context(), user.ID)
	}
	w.WriteHeader(http.StatusNoContent)
}
