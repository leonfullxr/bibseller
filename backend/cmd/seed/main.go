// Dev-only seeder: wipes and repopulates the database with realistic data
// across all four transfer-policy modes. Refuses to run in production.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/config"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

func main() {
	cfg := config.Load()
	if !cfg.IsDev() {
		log.Fatal("seed is dev-only: refusing to run with ENV != development")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	must(err)
	defer pool.Close()
	must(pool.Ping(ctx))

	_, err = pool.Exec(ctx, `TRUNCATE users, sessions, races, listings, chat_threads,
		messages, policy_acks, orders, order_events, stripe_events, reports, audit_log
		RESTART IDENTITY CASCADE`)
	must(err)

	q := sqlcgen.New(pool)

	users := seedUsers(ctx, q, pool)
	races := seedRaces(ctx, q, users["admin"])
	nListings := seedListings(ctx, q, races, users)

	fmt.Printf("seeded: %d users, %d races, %d listings\n", len(users), len(races), nListings)
}

const noLoginHash = "*seeded-account-no-login*"

func seedUsers(ctx context.Context, q *sqlcgen.Queries, pool *pgxpool.Pool) map[string]uuid.UUID {
	out := map[string]uuid.UUID{}
	for _, u := range []struct {
		key, email, name, locale, country string
	}{
		{"admin", "admin@bibseller.dev", "Admin", "en", "ES"},
		{"marta", "marta@example.com", "Marta R.", "es", "ES"},
		{"jonas", "jonas@example.com", "Jonas K.", "en", "DE"},
		{"claire", "claire@example.com", "Claire D.", "fr", "FR"},
		{"luca", "luca@example.com", "Luca B.", "en", "IT"},
	} {
		row, err := q.CreateUser(ctx, sqlcgen.CreateUserParams{
			ID: ids.New(), Email: u.email, PasswordHash: noLoginHash,
			DisplayName: u.name, Locale: u.locale, Country: ptr(u.country),
		})
		must(err)
		out[u.key] = row.ID
	}
	_, err := pool.Exec(ctx,
		`UPDATE users SET role = 'admin', email_verified_at = now() WHERE id = $1`,
		out["admin"])
	must(err)
	return out
}

type raceSeed struct {
	key, slug, name, sport, distance, city, country string
	date                                            string // YYYY-MM-DD
	policy                                          string
	officialURL, sourceURL, notes                   string
	status                                          string
}

func seedRaces(ctx context.Context, q *sqlcgen.Queries, admin uuid.UUID) map[string]sqlcgen.Race {
	seeds := []raceSeed{
		// platform_sale — resale allowed, payments possible (source required)
		{key: "munich", slug: "munich-marathon-2026", name: "Munich Marathon 2026", sport: "running", distance: "marathon", city: "Munich", country: "DE", date: "2026-10-11", policy: "platform_sale", sourceURL: "https://example.org/munich-tos#transfers", notes: "T&C §7: bib transfer permitted until race week."},
		{key: "granada", slug: "granada-half-2027", name: "Granada Media Maratón 2027", sport: "running", distance: "half", city: "Granada", country: "ES", date: "2027-03-28", policy: "platform_sale", sourceURL: "https://example.org/granada-reglamento#cesion"},
		{key: "brussels", slug: "brussels-20km-2027", name: "Brussels 20km 2027", sport: "running", distance: "20k", city: "Brussels", country: "BE", date: "2027-05-30", policy: "platform_sale", sourceURL: "https://example.org/brussels-rules#transfer"},
		{key: "garda", slug: "garda-trail-42k-2026", name: "Garda Trail 42K 2026", sport: "trail", distance: "42k", city: "Riva del Garda", country: "IT", date: "2026-09-19", policy: "platform_sale", sourceURL: "https://example.org/garda-regolamento#cambio"},
		{key: "mallorca", slug: "mallorca-70-3-2027", name: "Mallorca 70.3 Triathlon 2027", sport: "triathlon", distance: "70.3", city: "Alcúdia", country: "ES", date: "2027-05-08", policy: "platform_sale", sourceURL: "https://example.org/mallorca-athlete-guide#slot-transfer"},

		// official_only — the race runs its own name-change process
		{key: "valencia", slug: "valencia-marathon-2026", name: "Valencia Marathon 2026", sport: "running", distance: "marathon", city: "Valencia", country: "ES", date: "2026-12-06", policy: "official_only", officialURL: "https://example.org/valencia/name-change"},
		{key: "valencia-half", slug: "valencia-half-2026", name: "Valencia Half Marathon 2026", sport: "running", distance: "half", city: "Valencia", country: "ES", date: "2026-10-25", policy: "official_only", officialURL: "https://example.org/valencia-half/name-change"},
		{key: "paris", slug: "paris-marathon-2027", name: "Paris Marathon 2027", sport: "running", distance: "marathon", city: "Paris", country: "FR", date: "2027-04-11", policy: "official_only", officialURL: "https://example.org/paris/official-resale"},
		{key: "amsterdam", slug: "amsterdam-marathon-2026", name: "Amsterdam Marathon 2026", sport: "running", distance: "marathon", city: "Amsterdam", country: "NL", date: "2026-10-18", policy: "official_only", officialURL: "https://example.org/amsterdam/transfer"},
		{key: "vienna", slug: "vienna-city-marathon-2027", name: "Vienna City Marathon 2027", sport: "running", distance: "marathon", city: "Vienna", country: "AT", date: "2027-04-18", policy: "official_only", officialURL: "https://example.org/vienna/name-change"},

		// connect_only — transfers restricted; chat only, strong disclaimer
		{key: "berlin", slug: "berlin-marathon-2026", name: "Berlin Marathon 2026", sport: "running", distance: "marathon", city: "Berlin", country: "DE", date: "2026-09-27", policy: "connect_only", notes: "Organizer forbids bib transfers."},
		{key: "porto", slug: "porto-marathon-2026", name: "Porto Marathon 2026", sport: "running", distance: "marathon", city: "Porto", country: "PT", date: "2026-11-08", policy: "connect_only"},
		{key: "frankfurt", slug: "frankfurt-marathon-2026", name: "Frankfurt Marathon 2026", sport: "running", distance: "marathon", city: "Frankfurt", country: "DE", date: "2026-10-25", policy: "connect_only"},
		{key: "milano", slug: "milano-marathon-2027", name: "Milano Marathon 2027", sport: "running", distance: "marathon", city: "Milan", country: "IT", date: "2027-04-04", policy: "connect_only"},
		{key: "rotterdam", slug: "rotterdam-marathon-2027", name: "Rotterdam Marathon 2027", sport: "running", distance: "marathon", city: "Rotterdam", country: "NL", date: "2027-04-11", policy: "connect_only"},

		// unknown — not yet verified; treated as connect_only with a badge
		{key: "sevilla", slug: "sevilla-marathon-2027", name: "Sevilla Marathon 2027", sport: "running", distance: "marathon", city: "Sevilla", country: "ES", date: "2027-02-21", policy: "unknown"},
		{key: "lisbon", slug: "lisbon-half-2027", name: "Lisbon Half Marathon 2027", sport: "running", distance: "half", city: "Lisbon", country: "PT", date: "2027-03-21", policy: "unknown"},
		{key: "krakow", slug: "krakow-marathon-2027", name: "Kraków Marathon 2027", sport: "running", distance: "marathon", city: "Kraków", country: "PL", date: "2027-04-18", policy: "unknown"},
		{key: "bilbao", slug: "bilbao-night-half-2026", name: "Bilbao Night Half 2026", sport: "running", distance: "half", city: "Bilbao", country: "ES", date: "2026-05-16", policy: "unknown"}, // past race: exercises date filters + expired listings
		{key: "madrid", slug: "madrid-marathon-2027", name: "Madrid Marathon 2027", sport: "running", distance: "marathon", city: "Madrid", country: "ES", date: "2027-04-25", policy: "unknown", status: "draft"},
	}

	now := time.Now()
	out := map[string]sqlcgen.Race{}
	for _, s := range seeds {
		date, err := time.Parse("2006-01-02", s.date)
		must(err)
		status := s.status
		if status == "" {
			status = "published"
		}
		p := sqlcgen.CreateRaceParams{
			ID: ids.New(), Slug: s.slug, Name: s.name, Sport: s.sport,
			Distance: ptr(s.distance), EventDate: date, City: s.city, Country: s.country,
			TransferPolicy: s.policy, Status: status, CreatedBy: &admin,
		}
		if s.officialURL != "" {
			p.OfficialTransferUrl = &s.officialURL
		}
		if s.sourceURL != "" {
			p.PolicySourceUrl = &s.sourceURL
		}
		if s.notes != "" {
			p.PolicyNotes = &s.notes
		}
		if s.policy != "unknown" { // verified policies carry their evidence trail
			p.PolicyVerifiedAt = &now
			p.PolicyVerifiedBy = &admin
		}
		race, err := q.CreateRace(ctx, p)
		must(err)
		out[s.key] = race
	}
	return out
}

func seedListings(ctx context.Context, q *sqlcgen.Queries, races map[string]sqlcgen.Race, users map[string]uuid.UUID) int {
	listings := []struct {
		race, seller string
		price, orig  int32
		desc         string
		finalStatus  string // applied via the guarded status update; "" = active
	}{
		{race: "munich", seller: "marta", price: 8500, orig: 9000, desc: "Can't run due to injury. Will help with the transfer paperwork."},
		{race: "munich", seller: "claire", price: 9000, orig: 9000, desc: "Selling at face value."},
		{race: "munich", seller: "luca", price: 8000, orig: 9000, finalStatus: "sold"},
		{race: "garda", seller: "luca", price: 7500, orig: 8000, desc: "Trail entry incl. race pack."},
		{race: "mallorca", seller: "claire", price: 18000, orig: 19500, desc: "70.3 slot, transfer window open until April."},
		{race: "valencia", seller: "marta", price: 6500, orig: 6500, desc: "Name change via the official process; price is the entry fee."},
		{race: "valencia", seller: "claire", price: 6000, orig: 6500},
		{race: "paris", seller: "claire", price: 12000, orig: 12500, desc: "Official resale opens in January."},
		{race: "berlin", seller: "jonas", price: 15000, orig: 15000, desc: "Contact me to discuss options."},
		{race: "berlin", seller: "marta", price: 14000, orig: 15000},
		{race: "porto", seller: "luca", price: 5500, orig: 6000, finalStatus: "cancelled"},
		{race: "sevilla", seller: "marta", price: 5000, orig: 5500, desc: "Policy still unverified — chat first."},
		{race: "bilbao", seller: "claire", price: 3500, orig: 4000, finalStatus: "expired"},
	}

	for _, l := range listings {
		race := races[l.race]
		created, err := q.CreateListing(ctx, sqlcgen.CreateListingParams{
			ID: ids.New(), RaceID: race.ID, SellerID: users[l.seller],
			PriceCents: ptr(l.price), Currency: "EUR", OriginalPriceCents: ptr(l.orig),
			Description: nonEmpty(l.desc), ExpiresAt: race.EventDate,
		})
		must(err)
		if l.finalStatus != "" {
			_, err = q.UpdateListingStatus(ctx, sqlcgen.UpdateListingStatusParams{
				ToStatus: l.finalStatus, ID: created.ID, FromStatus: "active",
			})
			must(err)
		}
	}
	return len(listings)
}

func ptr[T any](v T) *T { return &v }

func nonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
