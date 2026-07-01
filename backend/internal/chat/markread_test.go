package chat_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// seedRaceOwnCountry mirrors seedRace but with its own country instead of the
// shared "ZX" literal - the race package's TestMapCounts treats "ZX" as
// exclusively its own, so another concurrent "ZX" race (this package already
// has several) intermittently tips its count. Not this issue's to fix; this
// test just avoids adding to it.
func seedRaceOwnCountry(t *testing.T, pool *pgxpool.Pool) sqlcgen.Race {
	t.Helper()
	id := ids.New()
	src := "https://example.org/policy"
	race, err := sqlcgen.New(pool).CreateRace(context.Background(), sqlcgen.CreateRaceParams{
		ID: id, Slug: "tr-" + id.String(), Name: "Test Race", Sport: "running",
		EventDate: future(), City: "Testville", Country: "QQ",
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

// TestMarkThreadReadSkipsRedundantWrite proves #96's fix: a poll that reports a
// cursor which doesn't advance the reader's last-read must not write to
// chat_threads (the hot-table churn the issue is about). :execrows lets this be
// asserted directly on the rows the UPDATE actually touched.
func TestMarkThreadReadSkipsRedundantWrite(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	sellerTok, sellerID := registerUser(t, h, pool, "Seller", true)
	buyerTok, _ := registerUser(t, h, pool, "Buyer", true)
	race := seedRaceOwnCountry(t, pool)
	listingID := createListing(t, h, race.ID, sellerTok)
	threadIDStr := startThread(t, h, listingID, buyerTok, "hello")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		t.Fatalf("parse thread id: %v", err)
	}

	q := sqlcgen.New(pool)
	ctx := context.Background()

	t1 := time.Now().UTC().Truncate(time.Millisecond)
	rows, err := q.MarkThreadRead(ctx, sqlcgen.MarkThreadReadParams{Reader: sellerID, ReadAt: &t1, ID: threadID})
	if err != nil || rows != 1 {
		t.Fatalf("first mark-read: rows = %d, err = %v, want 1", rows, err)
	}

	// Same cursor reported again (a redundant poll): must not write.
	rows, err = q.MarkThreadRead(ctx, sqlcgen.MarkThreadReadParams{Reader: sellerID, ReadAt: &t1, ID: threadID})
	if err != nil || rows != 0 {
		t.Fatalf("redundant mark-read (same cursor): rows = %d, err = %v, want 0", rows, err)
	}

	// An older cursor than what's stored: must not regress, must not write.
	earlier := t1.Add(-time.Hour)
	rows, err = q.MarkThreadRead(ctx, sqlcgen.MarkThreadReadParams{Reader: sellerID, ReadAt: &earlier, ID: threadID})
	if err != nil || rows != 0 {
		t.Fatalf("stale mark-read: rows = %d, err = %v, want 0", rows, err)
	}

	// A genuinely newer cursor: must advance and write.
	t2 := t1.Add(time.Minute)
	rows, err = q.MarkThreadRead(ctx, sqlcgen.MarkThreadReadParams{Reader: sellerID, ReadAt: &t2, ID: threadID})
	if err != nil || rows != 1 {
		t.Fatalf("advancing mark-read: rows = %d, err = %v, want 1", rows, err)
	}
}
