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
	user, err := h.q.GetUserByEmail(r.Context(), strings.TrimSpace(req.Email))
	switch {
	case err == nil:
		// Invalidate outstanding tokens before issuing a fresh one. If that
		// clear fails, do NOT issue another token - stacking a new one on top of
		// stale ones would break the single-link model. Log and fall through.
		if err := h.q.DeletePasswordResetsForUser(r.Context(), user.ID); err != nil {
			slog.Error("reset request: clearing old tokens failed", "err", err, "user_id", user.ID)
			break
		}
		h.startPasswordReset(r.Context(), user.ID, user.Email)
	case errors.Is(err, pgx.ErrNoRows):
		// Unknown email - the no-enumeration path; fall through to the same 204.
	default:
		// An operational failure, not "no such user". Log it for visibility but
		// still answer 204 so the response never reveals account existence.
		slog.Error("reset request: user lookup failed", "err", err)
	}
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

	// Hash before opening the transaction: argon2id is CPU-bound and must not
	// hold a DB transaction open.
	hash, err := hashPassword(req.Password)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}

	// One transaction so the reset is all-or-nothing: consuming the token,
	// setting the password, verifying the email, and revoking every session
	// either all commit or all roll back. A partial reset must never leave old
	// sessions valid under a new password.
	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	defer func() { _ = tx.Rollback(r.Context()) }() // no-op once committed
	q := h.q.WithTx(tx)

	// Atomic consume: a replayed or concurrent token finds no row to delete.
	userID, err := q.ConsumePasswordReset(r.Context(), hashToken(req.Token))
	if errors.Is(err, pgx.ErrNoRows) {
		// Unknown, already-consumed, or expired - indistinguishable to the caller.
		httpx.Error(w, http.StatusBadRequest, "invalid_token", "this reset link is invalid or has expired")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}

	if err := q.UpdateUserPassword(r.Context(), sqlcgen.UpdateUserPasswordParams{
		ID: userID, PasswordHash: hash,
	}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	// Using the emailed link proves inbox control, so the reset also verifies
	// the email (no-op if already verified).
	if err := q.MarkEmailVerified(r.Context(), userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	// Drop any sibling tokens and force a fresh login everywhere.
	if err := q.DeletePasswordResetsForUser(r.Context(), userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	if err := q.DeleteAllSessionsForUser(r.Context(), userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not reset password")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
