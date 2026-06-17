package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

type resetRequestRequest struct {
	Email string `json:"email"`
}

// requestPasswordReset starts the reset flow. It ALWAYS returns 204, whether or
// not the email maps to an account: answering differently would turn this into
// an account-enumeration oracle. Rate-limited like the other mail-sending
// endpoints. The token row is written synchronously so the link works the
// instant the email lands; the SMTP send runs in the background.
func (h *Handler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req resetRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	if user, err := h.q.GetUserByEmail(r.Context(), strings.TrimSpace(req.Email)); err == nil {
		// Invalidate outstanding tokens before issuing a fresh one.
		_ = h.q.DeletePasswordResetsForUser(r.Context(), user.ID)
		h.startPasswordReset(r.Context(), user.ID, user.Email)
	}
	// An unknown email (pgx.ErrNoRows) falls through to the same 204 - no oracle.
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) startPasswordReset(ctx context.Context, userID uuid.UUID, email string) {
	token, tokenHash, err := newToken()
	if err != nil {
		slog.Error("reset token mint failed", "err", err, "user_id", userID)
		return
	}
	if _, err := h.q.CreatePasswordReset(ctx, sqlcgen.CreatePasswordResetParams{
		TokenHash: tokenHash, UserID: userID,
	}); err != nil {
		slog.Error("reset token persist failed", "err", err, "user_id", userID)
		return
	}
	link := h.appURL + "/reset?token=" + token
	go func() {
		if err := h.mailer.SendPasswordReset(email, link); err != nil {
			slog.Error("reset email send failed", "err", err, "user_id", userID)
		}
	}()
}

type resetRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// resetPassword consumes a token from the emailed link and sets a new password.
// The token is the credential, so this endpoint needs no session. On success
// every session for the user is revoked (they sign in again with the new
// password) and the email is marked verified - receiving the link proves inbox
// control (founder decision: reset auto-verifies).
func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "token is required")
		return
	}
	if err := validatePassword(req.Password); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	userID, err := h.q.GetPasswordResetUser(r.Context(), hashToken(req.Token))
	if errors.Is(err, pgx.ErrNoRows) {
		// Unknown, already-consumed, or expired - indistinguishable to the caller.
		httpx.Error(w, http.StatusBadRequest, "invalid_token", "this reset link is invalid or has expired")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	if err := h.q.UpdateUserPassword(r.Context(), sqlcgen.UpdateUserPasswordParams{
		ID: userID, PasswordHash: hash,
	}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	// Using the emailed link proves inbox control, so the reset also verifies
	// the email (no-op if already verified).
	_ = h.q.MarkEmailVerified(r.Context(), userID)
	// One-time use, and force a fresh login everywhere under the new password.
	_ = h.q.DeletePasswordResetsForUser(r.Context(), userID)
	_ = h.q.DeleteAllSessionsForUser(r.Context(), userID)
	w.WriteHeader(http.StatusNoContent)
}
