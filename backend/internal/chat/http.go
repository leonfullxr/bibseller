// Package chat owns buyer<->seller conversations: the policy acknowledgment, the
// one-per-(listing,buyer) thread, message send and poll (HTTP polling, D13), and
// the inbox with unread counts. Every read and write is gated to the thread's
// two participants; in connect_only/unknown modes the buyer must acknowledge the
// race's terms before the first message (PRODUCT policy matrix).
package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/race"
)

const (
	maxMessageLen   = 4000       // matches the messages_content_check body bound
	messagePageSize = 100        // schema note: a poll fetches up to 100 messages
	maxImageBytes   = 5 << 20    // 5 MiB cap on an uploaded image
	maxImagePixels  = 40_000_000 // decode-bomb guard: ~40 MP (high-res phone photos still fit)
)

// Mailer sends the new-message notification. Declared here (the consumer) so
// chat stays independent of the SMTP implementation: cmd/api injects the shared
// email.SMTPMailer, tests inject a no-op.
type Mailer interface {
	SendNewMessage(to, link, locale string) error
}

// Storage holds and serves the private image objects. Declared here (the
// consumer) so chat depends on the capability, not the S3 implementation;
// cmd/api injects *storage.Client, tests inject the real client (skipped when
// the object store is unreachable).
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	IsNotFound(err error) bool
}

type Handler struct {
	pool       *pgxpool.Pool // for the send transaction (insert + thread touch)
	q          *sqlcgen.Queries
	mailer     Mailer
	storage    Storage
	appURL     string       // frontend base, for the inbox link in the email
	msgLimiter *rateLimiter // per-account cap on message sends
	ipLimiter  *rateLimiter // per-source-IP cap on message sends
}

func Routes(pool *pgxpool.Pool, mailer Mailer, store Storage, appURL string) func(*http.ServeMux) {
	h := &Handler{
		pool: pool, q: sqlcgen.New(pool), mailer: mailer, storage: store, appURL: appURL,
		msgLimiter: newRateLimiter(msgRateMax, msgRateWindow),
		ipLimiter:  newRateLimiter(ipRateMax, msgRateWindow),
	}
	go h.msgLimiter.sweep(msgRateWindow)
	go h.ipLimiter.sweep(msgRateWindow)
	return func(mux *http.ServeMux) {
		mux.HandleFunc("POST /races/{slug}/ack", h.ack)
		mux.HandleFunc("POST /listings/{id}/threads", h.startThread)
		mux.HandleFunc("POST /threads/{id}/messages", h.postMessage)
		mux.HandleFunc("GET /threads/{id}/messages", h.listMessages)
		mux.HandleFunc("GET /threads/{id}/messages/{mid}/image", h.getMessageImage)
		mux.HandleFunc("GET /threads/{id}", h.getThreadHeader)
		mux.HandleFunc("GET /threads", h.listThreads)
	}
}

type messageDTO struct {
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Body      *string   `json:"body"` // null for an image-only message
	HasImage  bool      `json:"has_image"`
	CreatedAt time.Time `json:"created_at"`
}

func toMessageDTO(m sqlcgen.Message) messageDTO {
	return messageDTO{
		ID: m.ID, SenderID: m.SenderID, Body: m.Body,
		HasImage: m.ImageKey != nil, CreatedAt: m.CreatedAt,
	}
}

type threadSummary struct {
	ID            uuid.UUID  `json:"id"`
	ListingID     uuid.UUID  `json:"listing_id"`
	RaceName      string     `json:"race_name"`
	RaceSlug      string     `json:"race_slug"`
	Role          string     `json:"role"`           // the caller's role: "buyer" | "seller"
	OtherParty    string     `json:"other_party"`    // display name of the other participant
	OtherPartyID  uuid.UUID  `json:"other_party_id"` // the other participant, for block/unblock
	LastMessageAt *time.Time `json:"last_message_at"`
	UnreadCount   int        `json:"unread_count"`
}

// ack records the buyer's acknowledgment of a race's transfer terms. Required
// before the first message in connect_only/unknown modes; idempotent and
// harmless in the others, so the frontend can call it without branching.
func (h *Handler) ack(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	rc, err := h.q.GetRaceBySlug(r.Context(), r.PathValue("slug"))
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !race.IsPublic(rc.Status)) {
		httpx.Error(w, http.StatusNotFound, "not_found", "race not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load race")
		return
	}
	if err := h.q.CreatePolicyAck(r.Context(), sqlcgen.CreatePolicyAckParams{
		ID: ids.New(), UserID: caller.ID, RaceID: rc.ID, Policy: rc.TransferPolicy,
	}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not record acknowledgment")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type sendRequest struct {
	Body string `json:"body"`
}

type startThreadResponse struct {
	ThreadID uuid.UUID  `json:"thread_id"`
	Message  messageDTO `json:"message"`
}

// startThread opens (or reuses) the caller's thread for a listing and posts the
// first message in one transaction - how a buyer contacts a seller. Buyers only:
// verified, not the listing's own seller, on an active listing, and past the
// policy-ack gate in connect_only/unknown modes.
func (h *Handler) startThread(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	listingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	lst, err := h.q.GetListingByID(r.Context(), listingID)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !race.IsPublic(lst.Race.Status)) {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load listing")
		return
	}
	if lst.Listing.SellerID == caller.ID {
		httpx.Error(w, http.StatusForbidden, "own_listing", "you cannot start a thread on your own listing")
		return
	}
	if lst.Listing.Status != "active" {
		httpx.Error(w, http.StatusConflict, "listing_not_active", "this listing is no longer active")
		return
	}
	if gatedPolicy(lst.Race.TransferPolicy) {
		_, err := h.q.GetPolicyAck(r.Context(), sqlcgen.GetPolicyAckParams{UserID: caller.ID, RaceID: lst.Race.ID})
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.Error(w, http.StatusForbidden, "ack_required", "acknowledge the race's terms before contacting the seller")
			return
		}
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal", "could not check acknowledgment")
			return
		}
	}

	if blocked, err := h.q.IsBlocked(r.Context(), sqlcgen.IsBlockedParams{BlockerID: caller.ID, BlockedID: lst.Listing.SellerID}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not check block status")
		return
	} else if blocked {
		httpx.Error(w, http.StatusForbidden, "blocked", "you cannot contact this seller")
		return
	}

	body, ok := h.decodeBody(w, r, caller.ID)
	if !ok {
		return
	}

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}
	defer func() { _ = tx.Rollback(r.Context()) }() // no-op once committed
	qtx := h.q.WithTx(tx)

	thread, err := qtx.CreateThread(r.Context(), sqlcgen.CreateThreadParams{
		ID: ids.New(), ListingID: listingID, BuyerID: caller.ID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// The listing left 'active' between the check above and this guarded
		// insert: no new thread is opened on a non-active listing.
		httpx.Error(w, http.StatusConflict, "listing_not_active", "this listing is no longer active")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not open thread")
		return
	}
	// A brand-new thread comes back with no last_message_at; that is the one
	// time the seller is emailed.
	firstMessage := thread.LastMessageAt == nil
	msg, ok := insertAndTouch(w, r, qtx, thread.ID, caller.ID, &body, nil)
	if !ok {
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}

	if firstMessage {
		h.notifySeller(lst.Listing.SellerID)
	}
	httpx.JSON(w, http.StatusCreated, startThreadResponse{ThreadID: thread.ID, Message: toMessageDTO(msg)})
}

// postMessage appends a message to an existing thread. Either participant may
// post; the listing's status no longer matters - an open conversation continues
// after the listing leaves the catalog.
func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	threadID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	tc, err := h.q.GetThreadParticipants(r.Context(), threadID)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load thread")
		return
	}
	if !isParticipant(caller.ID, tc.BuyerID, tc.SellerID) {
		httpx.Error(w, http.StatusForbidden, "forbidden", "you are not a participant in this thread")
		return
	}

	other := tc.SellerID
	if caller.ID == tc.SellerID {
		other = tc.BuyerID
	}
	if blocked, err := h.q.IsBlocked(r.Context(), sqlcgen.IsBlockedParams{BlockerID: caller.ID, BlockedID: other}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not check block status")
		return
	} else if blocked {
		httpx.Error(w, http.StatusForbidden, "blocked", "you cannot message in this thread")
		return
	}

	// A multipart body is an image upload (with an optional caption); JSON is a
	// plain text message.
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		h.postImageMessage(w, r, threadID, caller.ID)
		return
	}

	body, ok := h.decodeBody(w, r, caller.ID)
	if !ok {
		return
	}

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}
	defer func() { _ = tx.Rollback(r.Context()) }() // no-op once committed
	qtx := h.q.WithTx(tx)

	msg, ok := insertAndTouch(w, r, qtx, threadID, caller.ID, &body, nil)
	if !ok {
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}
	httpx.JSON(w, http.StatusCreated, toMessageDTO(msg))
}

type messageListResponse struct {
	Items      []messageDTO `json:"items"`
	NextCursor *string      `json:"next_cursor"`
}

// listMessages returns a thread's messages in id order, optionally only those
// newer than ?since=<cursor> (the polling path). Fetching also advances the
// caller's last-read mark to the newest message returned.
func (h *Handler) listMessages(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	threadID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	tc, err := h.q.GetThreadParticipants(r.Context(), threadID)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load thread")
		return
	}
	if !isParticipant(caller.ID, tc.BuyerID, tc.SellerID) {
		httpx.Error(w, http.StatusForbidden, "forbidden", "you are not a participant in this thread")
		return
	}

	params := sqlcgen.ListMessagesParams{ThreadID: threadID, PageSize: messagePageSize}
	if v := r.URL.Query().Get("since"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "malformed cursor")
			return
		}
		params.Cursor = &id
	}
	rows, err := h.q.ListMessages(r.Context(), params)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load messages")
		return
	}

	items := make([]messageDTO, len(rows))
	for i, row := range rows {
		items[i] = toMessageDTO(row)
	}
	// Reading advances the caller's last-read to the newest message fetched, so
	// the inbox unread count reflects what they have actually seen.
	if len(rows) > 0 {
		newest := rows[len(rows)-1].CreatedAt
		if err := h.q.MarkThreadRead(r.Context(), sqlcgen.MarkThreadReadParams{
			Reader: caller.ID, ReadAt: &newest, ID: threadID,
		}); err != nil {
			slog.Error("chat: mark-read failed", "err", err, "thread_id", threadID)
		}
	}

	resp := messageListResponse{Items: items}
	if len(rows) == messagePageSize {
		c := rows[len(rows)-1].ID.String()
		resp.NextCursor = &c
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, resp)
}

type threadListResponse struct {
	Items      []threadSummary `json:"items"`
	NextCursor *string         `json:"next_cursor"`
}

// listThreads is the caller's inbox: their threads as buyer or seller, newest
// activity first, each with the other party and the caller's unread count.
// Keyset-paginated on (last_message_at, id) - #97 - so a power-seller's inbox
// runs a bounded number of per-thread unread subqueries per call.
func (h *Handler) listThreads(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	limit, err := httpx.ParseLimit(r.URL.Query())
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	params := sqlcgen.ListThreadsForUserParams{Caller: caller.ID, PageSize: limit}
	if v := r.URL.Query().Get("cursor"); v != "" {
		at, id, err := parseThreadCursor(v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "malformed cursor")
			return
		}
		params.CursorLastMessageAt = &at
		params.CursorID = &id
	}
	rows, err := h.q.ListThreadsForUser(r.Context(), params)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load inbox")
		return
	}
	items := make([]threadSummary, len(rows))
	for i, row := range rows {
		role, other, otherID := "seller", row.BuyerName, row.BuyerID
		if row.BuyerID == caller.ID {
			role, other, otherID = "buyer", row.SellerName, row.SellerID
		}
		items[i] = threadSummary{
			ID: row.ID, ListingID: row.ListingID,
			RaceName: row.RaceName, RaceSlug: row.RaceSlug,
			Role: role, OtherParty: other, OtherPartyID: otherID,
			LastMessageAt: row.LastMessageAt, UnreadCount: int(row.UnreadCount),
		}
	}
	resp := threadListResponse{Items: items}
	if len(rows) == int(params.PageSize) {
		last := rows[len(rows)-1]
		if last.LastMessageAt != nil {
			c := formatThreadCursor(*last.LastMessageAt, last.ID)
			resp.NextCursor = &c
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, resp)
}

// getThreadHeader returns one thread's header (listing/race context, the
// other participant, unread count) so the single-thread page can render
// without fetching the whole inbox (#97).
func (h *Handler) getThreadHeader(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	threadID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	row, err := h.q.GetThreadHeader(r.Context(), sqlcgen.GetThreadHeaderParams{ID: threadID, Caller: caller.ID})
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load thread")
		return
	}
	if !isParticipant(caller.ID, row.BuyerID, row.SellerID) {
		httpx.Error(w, http.StatusForbidden, "forbidden", "you are not a participant in this thread")
		return
	}
	role, other, otherID := "seller", row.BuyerName, row.BuyerID
	if row.BuyerID == caller.ID {
		role, other, otherID = "buyer", row.SellerName, row.SellerID
	}
	w.Header().Set("Cache-Control", "no-store")
	httpx.JSON(w, http.StatusOK, threadSummary{
		ID: row.ID, ListingID: row.ListingID,
		RaceName: row.RaceName, RaceSlug: row.RaceSlug,
		Role: role, OtherParty: other, OtherPartyID: otherID,
		LastMessageAt: row.LastMessageAt, UnreadCount: int(row.UnreadCount),
	})
}

// Cursor format: "<RFC3339Nano>~<uuid>" - keyset position on
// (last_message_at, id) DESC, mirroring race.formatCursor's convention. UTC,
// per this repo's timestamptz convention - it also keeps the cursor free of a
// "+" offset, which a caller that doesn't URL-encode its query would corrupt.
func formatThreadCursor(at time.Time, id uuid.UUID) string {
	return at.UTC().Format(time.RFC3339Nano) + "~" + id.String()
}

func parseThreadCursor(s string) (time.Time, uuid.UUID, error) {
	atPart, idPart, ok := strings.Cut(s, "~")
	if !ok {
		return time.Time{}, uuid.Nil, errors.New("missing separator")
	}
	at, err := time.Parse(time.RFC3339Nano, atPart)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	id, err := uuid.Parse(idPart)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	return at, id, nil
}

// decodeBody parses and validates the message body and applies the per-account
// send budget. It writes the error response itself and returns ok=false so the
// caller just returns.
func (h *Handler) decodeBody(w http.ResponseWriter, r *http.Request, sender uuid.UUID) (string, bool) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_body", "request body must be JSON")
		return "", false
	}
	body, err := validateBody(req.Body)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return "", false
	}
	if !h.allowSend(w, r, sender) {
		return "", false
	}
	return body, true
}

// allowSend applies the message-send budgets. Per-account is the primary control;
// the generous per-IP cap stops one source flooding writes across many accounts
// without tripping a shared NAT. It writes the 429 itself and returns false.
func (h *Handler) allowSend(w http.ResponseWriter, r *http.Request, sender uuid.UUID) bool {
	now := time.Now()
	if allowed, retry := h.msgLimiter.allow("acct:"+sender.String(), now); !allowed {
		tooManyMessages(w, retry)
		return false
	}
	if allowed, retry := h.ipLimiter.allow(clientIP(r), now); !allowed {
		tooManyMessages(w, retry)
		return false
	}
	return true
}

func tooManyMessages(w http.ResponseWriter, retry int) {
	w.Header().Set("Retry-After", strconv.Itoa(retry))
	httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "too many messages, slow down")
}

// insertAndTouch writes the message and bumps the thread's last_message_at (and
// the sender's read mark) inside the caller's transaction. It writes the error
// response and returns ok=false on failure.
func insertAndTouch(w http.ResponseWriter, r *http.Request, qtx *sqlcgen.Queries, threadID, sender uuid.UUID, body, imageKey *string) (sqlcgen.Message, bool) {
	msg, err := qtx.InsertMessage(r.Context(), sqlcgen.InsertMessageParams{
		ID: ids.New(), ThreadID: threadID, SenderID: sender, Body: body, ImageKey: imageKey,
	})
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return sqlcgen.Message{}, false
	}
	// BuyerID here is the sqlc-generated name for the $2 "sender" parameter; the
	// query marks whichever side the sender is on as read.
	if err := qtx.TouchThreadOnMessage(r.Context(), sqlcgen.TouchThreadOnMessageParams{
		ID: threadID, BuyerID: sender,
	}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return sqlcgen.Message{}, false
	}
	return msg, true
}

// notifySeller emails the listing's seller that a buyer started a conversation.
// Fully backgrounded with its own short-lived context: this is best-effort, so
// it must not ride (and die with) the request context after the send commits.
// A failure only means no email; the seller still sees the thread in their inbox.
func (h *Handler) notifySeller(sellerID uuid.UUID) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		seller, err := h.q.GetUserByID(ctx, sellerID)
		if err != nil {
			slog.Error("chat: seller lookup for notification failed", "err", err, "seller_id", sellerID)
			return
		}
		if err := h.mailer.SendNewMessage(seller.Email, h.appURL+"/account/inbox", seller.Locale); err != nil {
			slog.Error("chat: new-message email send failed", "err", err, "seller_id", sellerID)
		}
	}()
}

func isParticipant(userID, buyerID, sellerID uuid.UUID) bool {
	return userID == buyerID || userID == sellerID
}

// gatedPolicy reports whether a race's policy requires a buyer acknowledgment
// before the first message (PRODUCT policy matrix): the venue-only modes.
func gatedPolicy(policy string) bool {
	return policy == "connect_only" || policy == "unknown"
}

// validateBody enforces the 1..4000 character bound (matching the DB CHECK) and
// rejects whitespace-only bodies. The text itself is not trimmed; only the
// emptiness test ignores surrounding whitespace.
func validateBody(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", errors.New("message body is required")
	}
	if utf8.RuneCountInString(raw) > maxMessageLen {
		return "", fmt.Errorf("message must be at most %d characters", maxMessageLen)
	}
	return raw, nil
}

// postImageMessage handles a multipart image upload on an existing thread; the
// caller is already a verified participant. The image is decoded and re-encoded
// (which drops any EXIF), only JPEG and PNG are accepted, the object key is
// random, and the bucket is private - the image is reachable only via
// getMessageImage. An optional text caption rides along as the message body.
func (h *Handler) postImageMessage(w http.ResponseWriter, r *http.Request, threadID, sender uuid.UUID) {
	if !h.allowSend(w, r, sender) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxImageBytes+(1<<20)) // image cap + multipart slack
	if err := r.ParseMultipartForm(maxImageBytes); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			httpx.Error(w, http.StatusRequestEntityTooLarge, "image_too_large", "the image is too large")
		} else {
			httpx.Error(w, http.StatusBadRequest, "invalid_body", "could not read the upload")
		}
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll() // remove any parts that spilled to temp files
		}
	}()

	// Validate the optional caption before storing anything, so a bad caption can
	// never orphan an uploaded object.
	var caption *string
	if c := strings.TrimSpace(r.FormValue("body")); c != "" {
		if utf8.RuneCountInString(c) > maxMessageLen {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter",
				fmt.Sprintf("caption must be at most %d characters", maxMessageLen))
			return
		}
		caption = &c
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "an image file is required")
		return
	}
	defer func() { _ = file.Close() }()

	// Guard against a decompression bomb before fully decoding: read just the
	// header for the dimensions, reject an absurd pixel count, then rewind.
	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_image", "could not read the image")
		return
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxImagePixels {
		httpx.Error(w, http.StatusBadRequest, "image_too_large", "the image dimensions are too large")
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not read the image")
		return
	}

	// Decoding proves it really is an image; re-encoding strips any EXIF metadata.
	img, format, err := image.Decode(file)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_image", "could not read the image")
		return
	}
	var buf bytes.Buffer
	var contentType, ext string
	switch format {
	case "jpeg":
		contentType, ext = "image/jpeg", "jpg"
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	case "png":
		contentType, ext = "image/png", "png"
		err = png.Encode(&buf, img)
	default:
		httpx.Error(w, http.StatusBadRequest, "unsupported_image", "only JPEG and PNG images are allowed")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not process the image")
		return
	}

	key := "threads/" + threadID.String() + "/" + ids.New().String() + "." + ext
	if err := h.storage.Put(r.Context(), key, &buf, int64(buf.Len()), contentType); err != nil {
		slog.Error("chat: image store failed", "err", err, "thread_id", threadID)
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not store the image")
		return
	}

	// The object now exists; any failure below must delete it to avoid an orphan.
	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.discardObject(key)
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}
	defer func() { _ = tx.Rollback(r.Context()) }() // no-op once committed
	qtx := h.q.WithTx(tx)

	msg, ok := insertAndTouch(w, r, qtx, threadID, sender, caption, &key)
	if !ok {
		h.discardObject(key)
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		h.discardObject(key)
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not send message")
		return
	}
	httpx.JSON(w, http.StatusCreated, toMessageDTO(msg))
}

// discardObject best-effort deletes an uploaded image whose message row did not
// commit, so a failed send leaves no orphan in the bucket.
func (h *Handler) discardObject(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.storage.Delete(ctx, key); err != nil {
		slog.Error("chat: orphaned image cleanup failed", "err", err, "key", key)
	}
}

// getMessageImage streams a message's private image to a thread participant.
func (h *Handler) getMessageImage(w http.ResponseWriter, r *http.Request) {
	caller, ok := auth.UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "not signed in")
		return
	}
	if caller.EmailVerifiedAt == nil {
		httpx.Error(w, http.StatusForbidden, "email_unverified", "verify your email before chatting")
		return
	}
	threadID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "image not found")
		return
	}
	msgID, err := uuid.Parse(r.PathValue("mid"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "image not found")
		return
	}
	row, err := h.q.GetMessageImage(r.Context(), sqlcgen.GetMessageImageParams{ID: msgID, ThreadID: threadID})
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, http.StatusNotFound, "not_found", "image not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load the image")
		return
	}
	if !isParticipant(caller.ID, row.BuyerID, row.SellerID) {
		httpx.Error(w, http.StatusForbidden, "forbidden", "you are not a participant in this thread")
		return
	}
	if row.ImageKey == nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "this message has no image")
		return
	}

	obj, err := h.storage.Get(r.Context(), *row.ImageKey)
	if err != nil {
		if h.storage.IsNotFound(err) {
			httpx.Error(w, http.StatusNotFound, "not_found", "image not found")
			return
		}
		slog.Error("chat: image fetch failed", "err", err, "message_id", msgID)
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load the image")
		return
	}
	defer func() { _ = obj.Close() }()
	w.Header().Set("Content-Type", contentTypeForKey(*row.ImageKey))
	w.Header().Set("Cache-Control", "private, no-store")
	_, _ = io.Copy(w, obj) // a mid-stream client disconnect is not actionable here
}

func contentTypeForKey(key string) string {
	if strings.HasSuffix(key, ".png") {
		return "image/png"
	}
	return "image/jpeg"
}
