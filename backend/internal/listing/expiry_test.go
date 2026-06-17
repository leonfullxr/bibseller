package listing

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

// seedRaceListing creates a race on eventDate with one listing, and (optionally)
// transitions the listing out of 'active'. Returns the listing id.
func seedRaceListing(t *testing.T, pool *pgxpool.Pool, sellerID uuid.UUID, eventDate time.Time, cancel bool) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	q := sqlcgen.New(pool)
	src := "https://example.org/source"
	raceID := ids.New()
	race, err := q.CreateRace(ctx, sqlcgen.CreateRaceParams{
		ID: raceID, Slug: "ex-" + raceID.String(), Name: "Expiry Race", Sport: "running",
		EventDate: eventDate, City: "Testville", Country: "ZX",
		TransferPolicy: "platform_sale", PolicySourceUrl: &src, Status: "published",
	})
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	price, orig := int32(5000), int32(6000)
	l, err := q.CreateListing(ctx, sqlcgen.CreateListingParams{
		ID: ids.New(), RaceID: race.ID, SellerID: sellerID,
		PriceCents: &price, Currency: "EUR", OriginalPriceCents: &orig,
		ExpiresAt: race.EventDate,
	})
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
	if cancel {
		if _, err := q.UpdateListingStatus(ctx, sqlcgen.UpdateListingStatusParams{
			ID: l.ID, FromStatus: "active", ToStatus: "cancelled",
		}); err != nil {
			t.Fatalf("cancel seed listing: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id = $1`, l.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM races WHERE id = $1`, race.ID)
	})
	return l.ID
}

func statusOf(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM listings WHERE id = $1`, id).Scan(&s); err != nil {
		t.Fatalf("status of %s: %v", id, err)
	}
	return s
}

func seedSeller(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	id := ids.New()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, 'x', 'Expiry Seller')`,
		id, id.String()+"@test.local"); err != nil {
		t.Fatalf("seed seller: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	})
	return id
}

func TestExpirePastRaceListings(t *testing.T) {
	pool := testdb.Pool(t)
	seller := seedSeller(t, pool)

	pastActive := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 0, -3), false)
	futureActive := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 1, 0), false)
	pastCancelled := seedRaceListing(t, pool, seller, time.Now().UTC().AddDate(0, 0, -3), true)

	n, ran, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC())
	if err != nil {
		t.Fatalf("expire: %v", err)
	}
	if !ran {
		t.Fatal("expiry did not run (lock not held)")
	}
	if n < 1 {
		t.Errorf("expired count = %d, want >= 1", n)
	}

	if s := statusOf(t, pool, pastActive); s != "expired" {
		t.Errorf("past active listing = %q, want expired", s)
	}
	if s := statusOf(t, pool, futureActive); s != "active" {
		t.Errorf("future active listing = %q, want active", s)
	}
	if s := statusOf(t, pool, pastCancelled); s != "cancelled" {
		t.Errorf("past cancelled listing = %q, want cancelled (untouched)", s)
	}

	// A race happening today is not past: re-running expires nothing new from
	// the today boundary. (Sanity: the function is idempotent for this set.)
	if _, _, err := expirePastRaceListings(context.Background(), pool, time.Now().UTC()); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if s := statusOf(t, pool, futureActive); s != "active" {
		t.Errorf("future listing changed on second run: %q", s)
	}
}

func TestExpiryAdvisoryLockSerializes(t *testing.T) {
	pool := testdb.Pool(t)
	ctx := context.Background()

	// Hold the expiry lock on a separate connection/transaction.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var held bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, expiryLockKey).Scan(&held); err != nil {
		t.Fatalf("acquire lock: %v", err)
	}
	if !held {
		t.Fatal("could not acquire the lock to set up the test")
	}

	// The job must see the lock as taken and skip, not block or run.
	n, ran, err := expirePastRaceListings(ctx, pool, time.Now().UTC())
	if err != nil {
		t.Fatalf("expire while locked: %v", err)
	}
	if ran {
		t.Error("expiry ran while another holder had the lock")
	}
	if n != 0 {
		t.Errorf("count = %d, want 0 when lock not acquired", n)
	}
}
