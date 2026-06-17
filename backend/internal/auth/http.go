// Package auth owns password hashing, the session lifecycle, and the
// /auth/* endpoints (docs/ARCHITECTURE.md -> Auth & sessions, decision D12).
//
// Follow-ups tracked for M3 completion: per-IP rate limiting on these
// endpoints and email verification (which gates listing/chat, not browsing).
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

const (
	minNameLen = 2
	maxNameLen = 50

	// NIST/OWASP guidance: length is the strength factor - minimum 8, no
	// composition rules (they push users toward predictable patterns). The
	// upper bound exists only so an attacker cannot post megabytes into a
	// memory-hard hash function.
	minPasswordLen = 8
	maxPasswordLen = 512

	maxEmailLen = 254 // RFC 5321 path limit
)

type Handler struct {
	q      *sqlcgen.Queries
	mailer Mailer
	appURL string // frontend base, for building verification links
	// loginLimiter throttles login attempts per account (by email), complementing
	// the per-IP RateLimit middleware so a distributed guessing attack can't dodge
	// the cap by rotating source addresses.
	loginLimiter *rateLimiter
}

func Routes(q *sqlcgen.Queries, mailer Mailer, appURL string) func(*http.ServeMux) {
	h := &Handler{q: q, mailer: mailer, appURL: appURL, loginLimiter: newRateLimiter(rateLimitMax, rateLimitWindow)}
	go h.loginLimiter.sweep(rateLimitWindow)
	return func(mux *http.ServeMux) {
		mux.HandleFunc("POST /auth/register", h.register)
		mux.HandleFunc("POST /auth/login", h.login)
		mux.HandleFunc("POST /auth/logout", h.logout)
		mux.HandleFunc("POST /auth/logout/all", h.logoutAll)
		mux.HandleFunc("GET /auth/me", h.me)
		mux.HandleFunc("POST /auth/verify", h.verify)
		mux.HandleFunc("POST /auth/verify/resend", h.resendVerification)
		mux.HandleFunc("POST /auth/password/reset/request", h.requestPasswordReset)
		mux.HandleFunc("POST /auth/password/reset", h.resetPassword)
		mux.HandleFunc("POST /auth/password", h.changePassword)
	}
}

// Account is the authenticated user's own view - email is fine here (it is
// their account), unlike the public user.Profile DTO.
type Account struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	DisplayName   string    `json:"display_name"`
	EmailVerified bool      `json:"email_verified"`
}

// sessionResponse is returned to the SvelteKit server, which translates
// Token into the __Host-session cookie. The raw token exists exactly twice:
// in this response body and in the browser's cookie jar - never in our
// database, which holds only its SHA-256.
type sessionResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      Account   `json:"user"`
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}

	email := strings.TrimSpace(req.Email)
	if len(email) > maxEmailLen || !validEmail(email) {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "email is not valid")
		return
	}
	name, err := ValidateDisplayName(req.DisplayName)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	if err := validatePassword(req.Password); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not create account")
		return
	}

	row, err := h.q.CreateUser(r.Context(), sqlcgen.CreateUserParams{
		ID: ids.New(), Email: email, PasswordHash: hash,
		DisplayName: name, Locale: "en",
	})
	// The citext UNIQUE constraint is the source of truth for "email taken" -
	// a prior SELECT would be a race. 409 admits the account exists, which is
	// an enumeration tradeoff every register endpoint makes; the mitigation
	// is rate limiting, not a lie the user can't act on.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		httpx.Error(w, http.StatusConflict, "email_taken", "an account with this email already exists")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not create account")
		return
	}

	// Fire off the verification email; never blocks or fails registration
	// (the user can resend). A fresh account is always unverified.
	h.startEmailVerification(r.Context(), row.ID, row.Email)
	h.respondWithSession(w, r, http.StatusCreated, row.ID, row.Email, row.DisplayName, false)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}

	email := strings.TrimSpace(req.Email)
	// Per-account throttle, counted before any DB or argon2 work. Keyed by the
	// lowercased email (matching citext), so it caps attempts on one account
	// regardless of which IP they come from. Applied to all attempts, not just
	// failures, to stay stateless and simple.
	if ok, retry := h.loginLimiter.allow("login:"+strings.ToLower(email), time.Now()); !ok {
		w.Header().Set("Retry-After", strconv.Itoa(retry))
		httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "too many attempts for this account, try again later")
		return
	}

	user, err := h.q.GetUserByEmail(r.Context(), email)
	if errors.Is(err, pgx.ErrNoRows) {
		// Unknown email burns the same argon2id work as a wrong password,
		// and returns the same code+message: neither timing nor wording may
		// reveal whether the account exists.
		_, _ = verifyPassword(req.Password, dummyHash)
		httpx.Error(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not log in")
		return
	}

	match, err := verifyPassword(req.Password, user.PasswordHash)
	if err != nil || !match {
		httpx.Error(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}

	// Rotation on login (spec): if the caller presented an existing session,
	// it dies here and a fresh token is minted. A token that predates this
	// authentication can never ride on it (session fixation defense).
	if old, ok := requestToken(r); ok {
		_ = h.q.DeleteSession(r.Context(), hashToken(old))
	}

	h.respondWithSession(w, r, http.StatusOK, user.ID, user.Email, user.DisplayName, user.EmailVerifiedAt != nil)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// Idempotent: logging out without (or with a dead) session is success,
	// not an error - the end state "no session" is what was asked for.
	if token, ok := requestToken(r); ok {
		if err := h.q.DeleteSession(r.Context(), hashToken(token)); err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal", "could not log out")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// logoutAll revokes every session for the signed-in user - the "log out all
// devices" control. The caller's own session goes too; the SvelteKit action
// then clears the cookie and sends them to the login page.
func (h *Handler) logoutAll(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if err := h.q.DeleteAllSessionsForUser(r.Context(), user.ID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not log out")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// me resolves the presented session cookie to the account - the endpoint the
// SvelteKit server hook will call to populate locals.user on each request.
func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	row, ok := UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, Account{
		ID: row.ID, Email: row.Email, DisplayName: row.DisplayName,
		EmailVerified: row.EmailVerifiedAt != nil,
	})
}

func (h *Handler) respondWithSession(w http.ResponseWriter, r *http.Request, status int, id uuid.UUID, email, name string, emailVerified bool) {
	token, expiresAt, err := issueSession(r.Context(), h.q, id, r)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not create session")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, status, sessionResponse{
		Token: token, ExpiresAt: expiresAt,
		User: Account{ID: id, Email: email, DisplayName: name, EmailVerified: emailVerified},
	})
}

// ValidateDisplayName trims and bounds-checks a display name, returning the
// cleaned value. The single source of the 2..50 rule, shared by registration
// here and profile updates in internal/user (which imports auth).
func ValidateDisplayName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if n := utf8.RuneCountInString(name); n < minNameLen || n > maxNameLen {
		return "", fmt.Errorf("display_name must be %d..%d characters", minNameLen, maxNameLen)
	}
	return name, nil
}

// validatePassword enforces the length policy shared by registration, reset,
// and change-password. Length is the only rule (see the minPasswordLen note).
func validatePassword(pw string) error {
	if n := len(pw); n < minPasswordLen || n > maxPasswordLen {
		return fmt.Errorf("password must be at least %d characters", minPasswordLen)
	}
	return nil
}

// validEmail accepts what net/mail parses as a single bare address. This is
// deliberately shallow - the only proof of a deliverable address is the
// verification email (M3 follow-up), not a cleverer regex.
func validEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}
