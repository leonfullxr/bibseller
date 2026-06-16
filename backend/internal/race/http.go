// Package race owns the public race catalog: browse with filters and
// cursor pagination, and race detail by slug. Read-only in M2.
package race

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

var (
	validSports = map[string]bool{
		"running": true, "trail": true, "triathlon": true,
		"cycling": true, "obstacle": true, "other": true,
	}
	validPolicies = map[string]bool{
		"platform_sale": true, "official_only": true,
		"connect_only": true, "unknown": true,
	}
)

// IsPublic reports whether a race with this status is visible in the public
// catalog. "published" is the only public status; drafts and anything else are
// hidden. Centralized here so race detail and listing reads share one
// definition instead of each testing the string.
func IsPublic(status string) bool {
	return status == "published"
}

type Handler struct {
	q *sqlcgen.Queries
}

func Routes(q *sqlcgen.Queries) func(*http.ServeMux) {
	h := &Handler{q: q}
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /races", h.list)
		mux.HandleFunc("GET /races/{slug}", h.get)
	}
}

type Summary struct {
	ID             uuid.UUID `json:"id"`
	Slug           string    `json:"slug"`
	Name           string    `json:"name"`
	Series         *string   `json:"series"`
	Sport          string    `json:"sport"`
	Distance       *string   `json:"distance"`
	EventDate      string    `json:"event_date"` // YYYY-MM-DD
	City           string    `json:"city"`
	Country        string    `json:"country"`
	TransferPolicy string    `json:"transfer_policy"`
	ActiveListings int64     `json:"active_listings"`
}

type Detail struct {
	Summary
	WebsiteURL          *string    `json:"website_url"`
	OfficialTransferURL *string    `json:"official_transfer_url"`
	PolicyNotes         *string    `json:"policy_notes"`
	PolicyVerifiedAt    *time.Time `json:"policy_verified_at"`
}

type listResponse struct {
	Items      []Summary `json:"items"`
	NextCursor *string   `json:"next_cursor"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	limit, err := httpx.ParseLimit(qp)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	params := sqlcgen.ListRacesParams{PageSize: limit}

	if v := qp.Get("country"); v != "" {
		c := strings.ToUpper(v)
		if len(c) != 2 {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "country must be a 2-letter ISO code")
			return
		}
		params.Country = &c
	}
	if v := qp.Get("sport"); v != "" {
		if !validSports[v] {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "unknown sport")
			return
		}
		params.Sport = &v
	}
	if v := qp.Get("policy"); v != "" {
		if !validPolicies[v] {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "unknown transfer policy")
			return
		}
		params.TransferPolicy = &v
	}
	if v := qp.Get("date_from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "date_from must be YYYY-MM-DD")
			return
		}
		params.DateFrom = &t
	}
	if v := qp.Get("date_to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "date_to must be YYYY-MM-DD")
			return
		}
		params.DateTo = &t
	}
	if v := strings.TrimSpace(qp.Get("q")); v != "" {
		if len(v) > 100 {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "search query too long")
			return
		}
		params.Search = &v
	}
	if v := qp.Get("cursor"); v != "" {
		date, id, err := parseCursor(v)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid_parameter", "malformed cursor")
			return
		}
		params.CursorDate = &date
		params.CursorID = &id
	}

	rows, err := h.q.ListRaces(r.Context(), params)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not list races")
		return
	}

	items := make([]Summary, len(rows))
	for i, row := range rows {
		items[i] = Summary{
			ID: row.ID, Slug: row.Slug, Name: row.Name, Series: row.Series,
			Sport: row.Sport, Distance: row.Distance,
			EventDate: row.EventDate.Format("2006-01-02"),
			City:      row.City, Country: row.Country,
			TransferPolicy: row.TransferPolicy, ActiveListings: row.ActiveListings,
		}
	}

	resp := listResponse{Items: items}
	if len(rows) == int(params.PageSize) {
		last := rows[len(rows)-1]
		c := formatCursor(last.EventDate, last.ID)
		resp.NextCursor = &c
	}

	w.Header().Set("Cache-Control", httpx.CatalogCacheControl)
	httpx.JSON(w, http.StatusOK, resp)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	row, err := h.q.GetRaceBySlug(r.Context(), r.PathValue("slug"))
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Status != "published") {
		httpx.Error(w, http.StatusNotFound, "not_found", "race not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal", "could not load race")
		return
	}

	w.Header().Set("Cache-Control", httpx.CatalogCacheControl)
	httpx.JSON(w, http.StatusOK, Detail{
		Summary: Summary{
			ID: row.ID, Slug: row.Slug, Name: row.Name, Series: row.Series,
			Sport: row.Sport, Distance: row.Distance,
			EventDate: row.EventDate.Format("2006-01-02"),
			City:      row.City, Country: row.Country,
			TransferPolicy: row.TransferPolicy, ActiveListings: row.ActiveListings,
		},
		WebsiteURL:          row.WebsiteUrl,
		OfficialTransferURL: row.OfficialTransferUrl,
		PolicyNotes:         row.PolicyNotes,
		PolicyVerifiedAt:    row.PolicyVerifiedAt,
	})
}

// Cursor format: "<YYYY-MM-DD>~<uuid>" - keyset position on (event_date, id).
func formatCursor(date time.Time, id uuid.UUID) string {
	return date.Format("2006-01-02") + "~" + id.String()
}

func parseCursor(s string) (time.Time, uuid.UUID, error) {
	datePart, idPart, ok := strings.Cut(s, "~")
	if !ok {
		return time.Time{}, uuid.Nil, errors.New("missing separator")
	}
	date, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	id, err := uuid.Parse(idPart)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	return date, id, nil
}
