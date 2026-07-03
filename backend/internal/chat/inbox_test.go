package chat_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// seedInboxTestRace mirrors seedRace but with its own country - see #96/#97's
// shared note: the race package's TestMapCounts treats "ZX" as exclusively its
// own, so concurrent packages must not add more "ZX" rows.
func seedInboxTestRace(t *testing.T, pool *pgxpool.Pool) sqlcgen.Race {
	t.Helper()
	id := ids.New()
	src := "https://example.org/policy"
	race, err := sqlcgen.New(pool).CreateRace(context.Background(), sqlcgen.CreateRaceParams{
		ID: id, Slug: "tr-" + id.String(), Name: "Test Race", Sport: "running",
		EventDate: future(), City: "Testville", Country: "QW",
		TransferPolicy: "platform_sale", Status: "published", PolicySourceUrl: &src,
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM messages WHERE thread_id IN
			(SELECT t.id FROM chat_threads t JOIN listings l ON l.id = t.listing_id WHERE l.race_id = $1)`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM chat_threads WHERE listing_id IN (SELECT id FROM listings WHERE race_id = $1)`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM policy_acks WHERE race_id = $1`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE race_id = $1`, race.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
	})
	return race
}

// TestInboxPagination proves #97: a page smaller than the caller's thread
// count returns exactly page_size items plus a next_cursor, and paging
// through with it covers every thread exactly once.
func TestInboxPagination(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedInboxTestRace(t, pool)

	threadIDs := map[string]bool{}
	for i := 0; i < 3; i++ {
		sellerTok, _ := registerUser(t, h, pool, "Seller", true)
		listingID := createListing(t, h, race.ID, sellerTok)
		threadIDs[startThread(t, h, listingID, buyerTok, "hi")] = true
	}

	rec := doJSON(t, h, http.MethodGet, "/api/v1/threads?limit=2", "", buyerTok)
	if rec.Code != http.StatusOK {
		t.Fatalf("page1: status = %d, body = %s", rec.Code, rec.Body)
	}
	var page1 struct {
		Items      []inboxThread `json:"items"`
		NextCursor *string       `json:"next_cursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page1); err != nil {
		t.Fatalf("page1: bad JSON: %v", err)
	}
	if len(page1.Items) != 2 || page1.NextCursor == nil {
		t.Fatalf("page1: items = %d, cursor = %v, want 2 items + a cursor", len(page1.Items), page1.NextCursor)
	}

	rec = doJSON(t, h, http.MethodGet, "/api/v1/threads?limit=2&cursor="+*page1.NextCursor, "", buyerTok)
	if rec.Code != http.StatusOK {
		t.Fatalf("page2: status = %d, body = %s", rec.Code, rec.Body)
	}
	var page2 struct {
		Items      []inboxThread `json:"items"`
		NextCursor *string       `json:"next_cursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page2); err != nil {
		t.Fatalf("page2: bad JSON: %v", err)
	}
	if len(page2.Items) != 1 || page2.NextCursor != nil {
		t.Fatalf("page2: items = %d, cursor = %v, want the final 1 item + no cursor", len(page2.Items), page2.NextCursor)
	}

	seen := map[string]bool{}
	for _, it := range append(page1.Items, page2.Items...) {
		if seen[it.ID] {
			t.Errorf("thread %s returned on both pages", it.ID)
		}
		seen[it.ID] = true
	}
	if len(seen) != len(threadIDs) {
		t.Fatalf("paged through %d distinct threads, want %d", len(seen), len(threadIDs))
	}
	for id := range threadIDs {
		if !seen[id] {
			t.Errorf("thread %s missing from paged results", id)
		}
	}
}

// TestInboxDefaultLimitIsGenerous is a regression check (Copilot review on
// #111): the frontend inbox list calls /api/v1/threads with no limit/cursor
// at all, so defaulting this endpoint to the catalog's DefaultPageSize (24)
// would silently truncate any caller with more threads than that, with no
// paging UI to reach the rest. It must default to something bigger.
func TestInboxDefaultLimitIsGenerous(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedInboxTestRace(t, pool)

	const threadCount = 26 // > httpx.DefaultPageSize (24)
	for i := 0; i < threadCount; i++ {
		sellerTok, _ := registerUser(t, h, pool, "Seller", true)
		listingID := createListing(t, h, race.ID, sellerTok)
		startThread(t, h, listingID, buyerTok, "hi")
	}

	rec := doJSON(t, h, http.MethodGet, "/api/v1/threads", "", buyerTok)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Items []inboxThread `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if len(resp.Items) != threadCount {
		t.Fatalf("no-limit fetch: items = %d, want %d (all of them, untruncated)", len(resp.Items), threadCount)
	}
}

// TestInboxExcludesThreadWithNoMessage proves the last_message_at IS NOT NULL
// guard (Copilot review on #111): a thread inserted directly, bypassing
// CreateThread/TouchThreadOnMessage (as at least one existing test fixture
// already does - internal/chat/retention_test.go's seedThreadMessage), has a
// NULL last_message_at. Without the guard this would sort first under
// ORDER BY ... DESC and never satisfy the keyset cursor comparison, breaking
// pagination for whoever else is in the caller's inbox. The fetch must
// succeed and simply omit it - there's no message to show for it anyway.
func TestInboxExcludesThreadWithNoMessage(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	buyerTok, buyerID := registerUser(t, h, pool, "Buyer", true)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	race := seedInboxTestRace(t, pool)
	listingID := createListing(t, h, race.ID, sellerTok)
	realThreadID := startThread(t, h, listingID, buyerTok, "hi")

	listing2ID := createListing(t, h, race.ID, sellerTok)
	listing2, err := uuid.Parse(listing2ID)
	if err != nil {
		t.Fatalf("parse listing id: %v", err)
	}
	noMsgThreadID := ids.New()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO chat_threads (id, listing_id, buyer_id) VALUES ($1, $2, $3)`,
		noMsgThreadID, listing2, buyerID); err != nil {
		t.Fatalf("seed thread with no message: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM chat_threads WHERE id = $1`, noMsgThreadID)
	})

	items := inbox(t, h, buyerTok)
	found := map[string]bool{}
	for _, it := range items {
		found[it.ID] = true
	}
	if !found[realThreadID] {
		t.Errorf("real thread missing from inbox: %+v", items)
	}
	if found[noMsgThreadID.String()] {
		t.Errorf("no-message thread should be excluded, got it in: %+v", items)
	}
}

// TestGetThreadHeader proves #97: the single-thread page's header endpoint
// works for either participant, 403s a stranger, and 404s a missing thread.
func TestGetThreadHeader(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, _ := registerUser(t, h, pool, "Seller", true)
	buyerTok, buyerID := registerUser(t, h, pool, "Buyer", true)
	strangerTok, _ := registerUser(t, h, pool, "Stranger", true)
	race := seedInboxTestRace(t, pool)
	listingID := createListing(t, h, race.ID, sellerTok)
	threadID := startThread(t, h, listingID, buyerTok, "hi")

	rec := doJSON(t, h, http.MethodGet, "/api/v1/threads/"+threadID, "", sellerTok)
	if rec.Code != http.StatusOK {
		t.Fatalf("seller header: status = %d, body = %s", rec.Code, rec.Body)
	}
	var header struct {
		ID             string `json:"id"`
		Role           string `json:"role"`
		OtherParty     string `json:"other_party"`
		OtherPartyID   string `json:"other_party_id"`
		RaceName       string `json:"race_name"`
		TransferPolicy string `json:"transfer_policy"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &header); err != nil {
		t.Fatalf("seller header: bad JSON: %v", err)
	}
	if header.ID != threadID || header.Role != "seller" || header.OtherParty != "Buyer" || header.OtherPartyID != buyerID.String() {
		t.Fatalf("seller header: %+v", header)
	}
	// The race's transfer policy rides along so the thread page can render the
	// in-chat policy reminder (connect_only/unknown) without a second fetch.
	if header.TransferPolicy != race.TransferPolicy {
		t.Fatalf("seller header: transfer_policy = %q, want %q", header.TransferPolicy, race.TransferPolicy)
	}

	if rec := doJSON(t, h, http.MethodGet, "/api/v1/threads/"+threadID, "", buyerTok); rec.Code != http.StatusOK {
		t.Fatalf("buyer header: status = %d, body = %s", rec.Code, rec.Body)
	}

	if rec := doJSON(t, h, http.MethodGet, "/api/v1/threads/"+threadID, "", strangerTok); rec.Code != http.StatusForbidden {
		t.Fatalf("stranger header: status = %d, want 403", rec.Code)
	}

	if rec := doJSON(t, h, http.MethodGet, "/api/v1/threads/"+ids.New().String(), "", buyerTok); rec.Code != http.StatusNotFound {
		t.Fatalf("missing thread: status = %d, want 404", rec.Code)
	}
}
