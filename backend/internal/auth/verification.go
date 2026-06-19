package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

// Mailer sends the transactional auth emails. Declared here (the consumer) so
// the auth package stays independent of the SMTP implementation: cmd/api
// injects an email.SMTPMailer, tests inject a no-op.
type Mailer interface {
	SendVerification(to, link, locale string) error
	SendPasswordReset(to, link, locale string) error
}

// startEmailVerification mints a verification token, persists its hash, and
// emails the link. Best-effort by design: a mail (or token) failure is logged
// but never fails the surrounding request - the user can always resend. The
// token row is written synchronously so the link works the instant the email
// lands; the SMTP send, which can block on the network, runs in the background.
func (h *Handler) startEmailVerification(ctx context.Context, userID uuid.UUID, email, locale string) {
	token, tokenHash, err := newToken()
	if err != nil {
		slog.Error("verification token mint failed", "err", err, "user_id", userID)
		return
	}
	if _, err := h.q.CreateEmailVerification(ctx, sqlcgen.CreateEmailVerificationParams{
		TokenHash: tokenHash, UserID: userID,
	}); err != nil {
		slog.Error("verification token persist failed", "err", err, "user_id", userID)
		return
	}
	link := h.appURL + "/verify?token=" + token
	go func() {
		if err := h.mailer.SendVerification(email, link, locale); err != nil {
			slog.Error("verification email send failed", "err", err, "user_id", userID)
		}
	}()
}

type verifyRequest struct {
	Token string `json:"token"`
}

// verify consumes a token from the emailed link and marks the account verified.
// The token itself is the credential, so this endpoint needs no session.
func (h *Handler) verify(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "token is required")
		return
	}

	userID, err := h.q.GetEmailVerificationUser(r.Context(), hashToken(req.Token))
	if errors.Is(err, pgx.ErrNoRows) {
		// Unknown, already-consumed, or expired - indistinguishable to the caller.
		httpx.Error(w, http.StatusBadRequest, "invalid_token", "this verification link is invalid or has expired")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not verify email")
		return
	}

	if err := h.q.MarkEmailVerified(r.Context(), userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not verify email")
		return
	}
	// One-time use: consume every outstanding token for this user.
	_ = h.q.DeleteEmailVerificationsForUser(r.Context(), userID)
	w.WriteHeader(http.StatusNoContent)
}

// resendVerification issues a fresh link for the signed-in user. Idempotent for
// an already-verified account (204, no email). Rate-limited like login/register
// because it sends mail.
func (h *Handler) resendVerification(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if user.EmailVerifiedAt != nil {
		w.WriteHeader(http.StatusNoContent) // nothing to verify
		return
	}
	// Invalidate outstanding tokens before issuing a fresh one.
	_ = h.q.DeleteEmailVerificationsForUser(r.Context(), user.ID)
	h.startEmailVerification(r.Context(), user.ID, user.Email, user.Locale)
	w.WriteHeader(http.StatusNoContent)
}
