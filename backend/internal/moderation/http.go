// Package moderation owns the beta's minimal trust-and-safety write path:
// reporting a listing/message/user and blocking another user. The admin review
// queue and actioning are M7 (#11); blocks are enforced in internal/chat.
package moderation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ratelimit"
)

const maxDetailsLen = 2000

// Allowed values mirror the reports table CHECK constraints; validating here
// turns a bad value into a clean 400 instead of a constraint-violation 500.
var (
	validSubjectType = map[string]bool{"listing": true, "message": true, "user": true}
	validReason      = map[string]bool{"forbidden_transfer": true, "scam": true, "offensive": true, "other": true}
)

type Handler struct {
	q       *sqlcgen.Queries
	limiter *ratelimit.Limiter
}

func Routes(q *sqlcgen.Queries) func(*http.ServeMux) {
	h := &Handler{q: q, limiter: ratelimit.New(reportRateMax, reportRateWindow)}
	go h.limiter.Sweep(reportRateWindow)
	return func(mux *http.ServeMux) {
		mux.HandleFunc("POST /reports", h.createReport)
		mux.HandleFunc("POST /blocks", h.createBlock)
		mux.HandleFunc("DELETE /blocks/{id}", h.deleteBlock)
	}
}

type reportRequest struct {
	SubjectType string `json:"subject_type"`
	SubjectID   string `json:"subject_id"`
	Reason      string `json:"reason"`
	Details     string `json:"details"`
}

// createReport files a report against a listing, message, or user. Verified
// users only; rate-limited. The admin queue that reviews these is M7.
func (h *Handler) createReport(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before reporting")
		return
	}
	if !h.allow(w, caller.ID) {
		return
	}

	var req reportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	if !validSubjectType[req.SubjectType] {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "subject_type must be listing, message, or user")
		return
	}
	if !validReason[req.Reason] {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "reason is not valid")
		return
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "subject_id is not a valid id")
		return
	}
	var details *string
	if d := strings.TrimSpace(req.Details); d != "" {
		if utf8.RuneCountInString(d) > maxDetailsLen {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "details is too long")
			return
		}
		details = &d
	}

	reporterID := caller.ID
	row, err := h.q.CreateReport(r.Context(), sqlcgen.CreateReportParams{
		ID: ids.New(), ReporterID: &reporterID, SubjectType: req.SubjectType,
		SubjectID: subjectID, Reason: req.Reason, Details: details,
	})
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not file the report")
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]string{"id": row.ID.String()})
}

type blockRequest struct {
	BlockedID string `json:"blocked_id"`
}

// createBlock blocks another user; a block silences any conversation between
// them both ways (enforced in internal/chat). Idempotent.
func (h *Handler) createBlock(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if !h.allow(w, caller.ID) {
		return
	}
	var req blockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return
	}
	blockedID, err := uuid.Parse(req.BlockedID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "blocked_id is not a valid id")
		return
	}
	if blockedID == caller.ID {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "you cannot block yourself")
		return
	}
	err = h.q.CreateBlock(r.Context(), sqlcgen.CreateBlockParams{BlockerID: caller.ID, BlockedID: blockedID})
	// A foreign-key violation means the target user does not exist.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not block the user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteBlock unblocks a user. Idempotent: unblocking who you never blocked is success.
func (h *Handler) deleteBlock(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	blockedID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "not a valid user id")
		return
	}
	if err := h.q.DeleteBlock(r.Context(), sqlcgen.DeleteBlockParams{BlockerID: caller.ID, BlockedID: blockedID}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not unblock the user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) allow(w http.ResponseWriter, userID uuid.UUID) bool {
	if ok, retry := h.limiter.Allow(userID.String(), time.Now()); !ok {
		w.Header().Set("Retry-After", strconv.Itoa(retry))
		httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "too many requests, slow down")
		return false
	}
	return true
}
