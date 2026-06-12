// Package listing owns public listing reads: active listings per race and
// listing detail. Read-only in M2; seller mutations arrive in M4.
package listing

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

const (
	defaultPageSize = 24
	maxPageSize     = 100
	cacheControl    = "public, max-age=60, stale-while-revalidate=300"
)

type Handler struct {
	q *sqlcgen.Queries
}

func Routes(q *sqlcgen.Queries) func(chi.Router) {
	h := &Handler{q: q}
	return func(r chi.Router) {
		r.Get("/races/{slug}/listings", h.listByRace)
		r.Get("/listings/{id}", h.get)
	}
}

type Summary struct {
	ID                 uuid.UUID `json:"id"`
	Status             string    `json:"status"`
	PriceCents         *int32    `json:"price_cents"`
	Currency           string    `json:"currency"`
	OriginalPriceCents *int32    `json:"original_price_cents"`
	Description        *string   `json:"description"`
	SellerName         string    `json:"seller_name"`
	CreatedAt          time.Time `json:"created_at"`
}

type raceContext struct {
	Slug                string  `json:"slug"`
	Name                string  `json:"name"`
	Distance            *string `json:"distance"`
	EventDate           string  `json:"event_date"`
	City                string  `json:"city"`
	Country             string  `json:"country"`
	TransferPolicy      string  `json:"transfer_policy"`
	OfficialTransferURL *string `json:"official_transfer_url"`
}

type Detail struct {
	Summary
	Race raceContext `json:"race"`
}

type listResponse struct {
	Items      []Summary `json:"items"`
	NextCursor *string   `json:"next_cursor"`
}

func (h *Handler) listByRace(w http.ResponseWriter, r *http.Request) {
	race, err := h.q.GetRaceBySlug(r.Context(), chi.URLParam(r, "slug"))
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && race.Status != "published") {
		httpx.Error(w, http.StatusNotFound, "not_found", "race not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load race")
		return
	}

	params := sqlcgen.ListActiveListingsByRaceParams{
		RaceID:   race.ID,
		PageSize: defaultPageSize,
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > maxPageSize {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "invalid limit")
			return
		}
		params.PageSize = int32(n)
	}
	if v := r.URL.Query().Get("cursor"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "malformed cursor")
			return
		}
		params.CursorID = &id
	}

	rows, err := h.q.ListActiveListingsByRace(r.Context(), params)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not list listings")
		return
	}

	items := make([]Summary, len(rows))
	for i, row := range rows {
		items[i] = Summary{
			ID: row.ID, Status: row.Status,
			PriceCents: row.PriceCents, Currency: row.Currency,
			OriginalPriceCents: row.OriginalPriceCents,
			Description:        row.Description,
			SellerName:         row.SellerName, CreatedAt: row.CreatedAt,
		}
	}

	resp := listResponse{Items: items}
	if len(rows) == int(params.PageSize) {
		c := rows[len(rows)-1].ID.String()
		resp.NextCursor = &c
	}

	w.Header().Set("Cache-Control", cacheControl)
	httpx.JSON(w, http.StatusOK, resp)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}

	row, err := h.q.GetListingByID(r.Context(), id)
	// Listings on unpublished races are not public.
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Race.Status != "published") {
		httpx.Error(w, http.StatusNotFound, "not_found", "listing not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load listing")
		return
	}

	w.Header().Set("Cache-Control", cacheControl)
	httpx.JSON(w, http.StatusOK, Detail{
		Summary: Summary{
			ID: row.Listing.ID, Status: row.Listing.Status,
			PriceCents: row.Listing.PriceCents, Currency: row.Listing.Currency,
			OriginalPriceCents: row.Listing.OriginalPriceCents,
			Description:        row.Listing.Description,
			SellerName:         row.SellerName, CreatedAt: row.Listing.CreatedAt,
		},
		Race: raceContext{
			Slug: row.Race.Slug, Name: row.Race.Name, Distance: row.Race.Distance,
			EventDate: row.Race.EventDate.Format("2006-01-02"),
			City:      row.Race.City, Country: row.Race.Country,
			TransferPolicy:      row.Race.TransferPolicy,
			OfficialTransferURL: row.Race.OfficialTransferUrl,
		},
	})
}
