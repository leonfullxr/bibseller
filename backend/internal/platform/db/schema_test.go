// Integration tests for the schema's load-bearing constraints: they verify
// the database itself rejects illegal data, independent of service-layer
// validation. They need a migrated Postgres (CI provides one; locally run
// `make migrate` first) and skip when none is reachable.
package db_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// tx returns a rolled-back-on-cleanup transaction against the migrated test
// database, or skips the test when the database isn't available.
func tx(t *testing.T) pgx.Tx {
	t.Helper()
	ctx := context.Background()
	pool := testdb.Pool(t)

	transaction, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	t.Cleanup(func() { _ = transaction.Rollback(ctx) })
	return transaction
}

func wantConstraint(t *testing.T, err error, constraint string) {
	t.Helper()
	if err == nil {
		t.Fatalf("insert succeeded, want %q violation", constraint)
	}
	if !strings.Contains(err.Error(), constraint) {
		t.Fatalf("error = %v, want %q violation", err, constraint)
	}
}

// --- fixtures (raw SQL on purpose: these tests are about the schema) ---

func seedUser(t *testing.T, tx pgx.Tx) uuid.UUID {
	t.Helper()
	id := ids.New()
	_, err := tx.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, display_name)
		 VALUES ($1, $2, 'x', 'Test User')`, id, id.String()+"@test.local")
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return id
}

func seedRace(t *testing.T, tx pgx.Tx, policy string) uuid.UUID {
	t.Helper()
	id := ids.New()
	_, err := tx.Exec(context.Background(),
		`INSERT INTO races (id, slug, name, event_date, city, country, transfer_policy,
		                    official_transfer_url, policy_source_url, status)
		 VALUES ($1, $2, 'Test Race', '2027-01-01', 'Testville', 'ES', $3,
		         'https://example.org/official', 'https://example.org/source', 'published')`,
		id, "test-race-"+id.String(), policy)
	if err != nil {
		t.Fatalf("seed race: %v", err)
	}
	return id
}

func seedListing(t *testing.T, tx pgx.Tx, raceID, sellerID uuid.UUID) uuid.UUID {
	t.Helper()
	id := ids.New()
	_, err := tx.Exec(context.Background(),
		`INSERT INTO listings (id, race_id, seller_id, price_cents, expires_at)
		 VALUES ($1, $2, $3, 5000, '2027-01-01')`, id, raceID, sellerID)
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
	return id
}

func insertOrder(tx pgx.Tx, listingID, buyerID, sellerID uuid.UUID, state string, item, fee, total int) error {
	_, err := tx.Exec(context.Background(),
		`INSERT INTO orders (id, listing_id, buyer_id, seller_id, state,
		                     item_amount_cents, processing_fee_cents, total_amount_cents)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		ids.New(), listingID, buyerID, sellerID, state, item, fee, total)
	return err
}

// --- the constraints under test ---

func TestRejectsInvalidTransferPolicy(t *testing.T) {
	x := tx(t)
	_, err := x.Exec(context.Background(),
		`INSERT INTO races (id, slug, name, event_date, city, country, transfer_policy)
		 VALUES ($1, 'bad-policy-race', 'Bad', '2027-01-01', 'X', 'ES', 'free_for_all')`,
		ids.New())
	wantConstraint(t, err, "races_policy_check")
}

func TestOfficialOnlyRequiresOfficialURL(t *testing.T) {
	x := tx(t)
	_, err := x.Exec(context.Background(),
		`INSERT INTO races (id, slug, name, event_date, city, country, transfer_policy)
		 VALUES ($1, 'official-no-url', 'X', '2027-01-01', 'X', 'ES', 'official_only')`,
		ids.New())
	wantConstraint(t, err, "races_official_url_required")
}

func TestPlatformSaleRequiresPolicySource(t *testing.T) {
	x := tx(t)
	_, err := x.Exec(context.Background(),
		`INSERT INTO races (id, slug, name, event_date, city, country, transfer_policy)
		 VALUES ($1, 'platform-no-source', 'X', '2027-01-01', 'X', 'ES', 'platform_sale')`,
		ids.New())
	wantConstraint(t, err, "races_platform_source_required")
}

func TestRejectsNegativeListingPrice(t *testing.T) {
	x := tx(t)
	race := seedRace(t, x, "platform_sale")
	seller := seedUser(t, x)
	_, err := x.Exec(context.Background(),
		`INSERT INTO listings (id, race_id, seller_id, price_cents, expires_at)
		 VALUES ($1, $2, $3, -1, '2027-01-01')`, ids.New(), race, seller)
	wantConstraint(t, err, "listings_price_nonnegative")
}

func TestRejectsInvalidOrderState(t *testing.T) {
	x := tx(t)
	seller, buyer := seedUser(t, x), seedUser(t, x)
	listing := seedListing(t, x, seedRace(t, x, "platform_sale"), seller)
	err := insertOrder(x, listing, buyer, seller, "money_teleported", 5000, 0, 5000)
	wantConstraint(t, err, "orders_state_check")
}

func TestOrderTotalMustEqualItemPlusFee(t *testing.T) {
	x := tx(t)
	seller, buyer := seedUser(t, x), seedUser(t, x)
	listing := seedListing(t, x, seedRace(t, x, "platform_sale"), seller)
	err := insertOrder(x, listing, buyer, seller, "pending_payment", 5000, 100, 9999)
	wantConstraint(t, err, "orders_total_consistent")
}

func TestOneLiveOrderPerListing(t *testing.T) {
	x := tx(t)
	seller, buyer := seedUser(t, x), seedUser(t, x)
	listing := seedListing(t, x, seedRace(t, x, "platform_sale"), seller)

	if err := insertOrder(x, listing, buyer, seller, "pending_payment", 5000, 0, 5000); err != nil {
		t.Fatalf("first order: %v", err)
	}
	err := insertOrder(x, listing, seedUser(t, x), seller, "paid_held", 5000, 0, 5000)
	wantConstraint(t, err, "orders_one_live_per_listing")
}

func TestCompletedOrderDoesNotBlockNewOrder(t *testing.T) {
	x := tx(t)
	seller, buyer := seedUser(t, x), seedUser(t, x)
	listing := seedListing(t, x, seedRace(t, x, "platform_sale"), seller)

	if err := insertOrder(x, listing, buyer, seller, "refunded", 5000, 0, 5000); err != nil {
		t.Fatalf("terminal order: %v", err)
	}
	// Terminal states are excluded from the partial unique index.
	if err := insertOrder(x, listing, seedUser(t, x), seller, "pending_payment", 5000, 0, 5000); err != nil {
		t.Fatalf("new order after terminal one should be allowed: %v", err)
	}
}

func TestPolicyAckOncePerUserAndRace(t *testing.T) {
	x := tx(t)
	user := seedUser(t, x)
	race := seedRace(t, x, "connect_only")

	ack := func() error {
		_, err := x.Exec(context.Background(),
			`INSERT INTO policy_acks (id, user_id, race_id, policy)
			 VALUES ($1, $2, $3, 'connect_only')`, ids.New(), user, race)
		return err
	}
	if err := ack(); err != nil {
		t.Fatalf("first ack: %v", err)
	}
	wantConstraint(t, ack(), "policy_acks_once_per_race")
}

func TestEmailIsCaseInsensitivelyUnique(t *testing.T) {
	x := tx(t)
	insert := func(email string) error {
		_, err := x.Exec(context.Background(),
			`INSERT INTO users (id, email, password_hash, display_name)
			 VALUES ($1, $2, 'x', 'X')`, ids.New(), email)
		return err
	}
	if err := insert("Case@Test.local"); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	// citext: same address, different case must collide
	wantConstraint(t, insert("case@test.local"), "users_email_key")
}
