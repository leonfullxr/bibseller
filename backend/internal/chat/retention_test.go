package chat

import (
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/config"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/platform/storage"
)

func retentionStorage(t *testing.T) *storage.Client {
	t.Helper()
	cfg := config.Load()
	store, err := storage.New(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	if err != nil {
		t.Fatalf("storage: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := store.EnsureBucket(ctx); err != nil {
		t.Skipf("object storage not reachable, skipping retention test: %v", err)
	}
	return store
}

func seedUser(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	id := ids.New()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, 'x', 'Retention User')`,
		id, id.String()+"@test.local"); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id) })
	return id
}

// seedThreadMessage creates a race on eventDate with a listing, a thread, and a
// message carrying an image object. Returns the message id and image key.
func seedThreadMessage(t *testing.T, pool *pgxpool.Pool, store *storage.Client, seller, buyer uuid.UUID, eventDate time.Time) (uuid.UUID, string) {
	t.Helper()
	ctx := context.Background()
	q := sqlcgen.New(pool)
	raceID := ids.New()
	race, err := q.CreateRace(ctx, sqlcgen.CreateRaceParams{
		ID: raceID, Slug: "ret-" + raceID.String(), Name: "Retention Race", Sport: "running",
		EventDate: eventDate, City: "Testville", Country: "ZX",
		TransferPolicy: "connect_only", Status: "published",
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	l, err := q.CreateListing(ctx, sqlcgen.CreateListingParams{
		ID: ids.New(), RaceID: race.ID, SellerID: seller, Currency: "EUR", ExpiresAt: race.EventDate,
	})
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
	thread, err := q.CreateThread(ctx, sqlcgen.CreateThreadParams{ID: ids.New(), ListingID: l.ID, BuyerID: buyer})
	if err != nil {
		t.Fatalf("seed thread: %v", err)
	}
	key := "threads/" + thread.ID.String() + "/" + ids.New().String() + ".jpg"
	if err := store.Put(ctx, key, strings.NewReader("img"), 3, "image/jpeg"); err != nil {
		t.Fatalf("put object: %v", err)
	}
	body := "old message"
	msg, err := q.InsertMessage(ctx, sqlcgen.InsertMessageParams{
		ID: ids.New(), ThreadID: thread.ID, SenderID: buyer, Body: &body, ImageKey: &key,
	})
	if err != nil {
		t.Fatalf("seed message: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM messages WHERE thread_id = $1`, thread.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM chat_threads WHERE id = $1`, thread.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id = $1`, l.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
		_ = store.Delete(ctx, key)
	})
	return msg.ID, key
}

func messageExists(t *testing.T, pool *pgxpool.Pool, msgID uuid.UUID) bool {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM messages WHERE id = $1`, msgID).Scan(&n); err != nil {
		t.Fatalf("message exists: %v", err)
	}
	return n > 0
}

// TestPurgeExpiredMessages proves the retention job deletes messages (and their
// image objects) past the 12-month horizon while leaving recent ones intact.
func TestPurgeExpiredMessages(t *testing.T) {
	pool := testdb.Pool(t)
	store := retentionStorage(t)
	seller := seedUser(t, pool)
	buyer := seedUser(t, pool)

	now := time.Now().UTC()
	expired, expiredKey := seedThreadMessage(t, pool, store, seller, buyer, now.AddDate(0, -13, 0)) // past horizon
	recent, _ := seedThreadMessage(t, pool, store, seller, buyer, now.AddDate(0, -1, 0))            // within horizon

	n, ran, err := purgeExpiredMessages(context.Background(), pool, store, slog.New(slog.DiscardHandler), now)
	if err != nil {
		t.Fatalf("purge: %v", err)
	}
	if !ran {
		t.Fatal("purge did not run (lock not held)")
	}
	if n < 1 {
		t.Errorf("purged count = %d, want >= 1", n)
	}

	if messageExists(t, pool, expired) {
		t.Error("expired message still present after purge")
	}
	if !messageExists(t, pool, recent) {
		t.Error("recent (in-horizon) message was deleted")
	}
	if _, err := store.Get(context.Background(), expiredKey); !store.IsNotFound(err) {
		t.Errorf("expired image object still present: err = %v", err)
	}
}
