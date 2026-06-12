package user_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/user"
)

func seedUser(t *testing.T, pool *pgxpool.Pool) sqlcgen.User {
	t.Helper()
	ctx := context.Background()
	id := ids.New()
	row, err := sqlcgen.New(pool).CreateUser(ctx, sqlcgen.CreateUserParams{
		ID: id, Email: id.String() + "@test.local", PasswordHash: "x",
		DisplayName: "Original Name", Locale: "en",
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, row.ID)
	})
	return row
}

func handler(pool *pgxpool.Pool) http.Handler {
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool, user.Routes(sqlcgen.New(pool)))
}

func patch(t *testing.T, h http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	return rec
}

func TestGetUser(t *testing.T) {
	pool := testdb.Pool(t)
	u := seedUser(t, pool)
	h := handler(pool)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/users/"+u.ID.String(), nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var body user.Profile
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if body.ID != u.ID || body.DisplayName != "Original Name" {
		t.Errorf("unexpected body: %+v", body)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/users/"+ids.New().String(), nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown user: status = %d, want 404", rec.Code)
	}
}

func TestUpdateDisplayName(t *testing.T) {
	pool := testdb.Pool(t)
	u := seedUser(t, pool)
	h := handler(pool)

	rec := patch(t, h, "/api/v1/users/"+u.ID.String(), `{"display_name": "  New Name  "}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var body user.Profile
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if body.DisplayName != "New Name" {
		t.Errorf("display_name = %q, want trimmed %q", body.DisplayName, "New Name")
	}

	// The mutation must be visible in the database, not just echoed back.
	stored, err := sqlcgen.New(pool).GetUserByID(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if stored.DisplayName != "New Name" {
		t.Errorf("stored display_name = %q, want %q", stored.DisplayName, "New Name")
	}
}

func TestUpdateDisplayNameRejectsBadInput(t *testing.T) {
	pool := testdb.Pool(t)
	u := seedUser(t, pool)
	h := handler(pool)
	path := "/api/v1/users/" + u.ID.String()

	for name, body := range map[string]string{
		"not json":        `not json`,
		"missing field":   `{}`,
		"too short":       `{"display_name": "x"}`,
		"whitespace only": `{"display_name": "   "}`,
		"too long":        `{"display_name": "` + strings.Repeat("a", 51) + `"}`,
	} {
		if rec := patch(t, h, path, body); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", name, rec.Code)
		}
	}

	// Valid body, unknown user.
	if rec := patch(t, h, "/api/v1/users/"+ids.New().String(), `{"display_name": "Valid"}`); rec.Code != http.StatusNotFound {
		t.Errorf("unknown user: status = %d, want 404", rec.Code)
	}

	// Nothing above may have touched the stored name.
	stored, err := sqlcgen.New(pool).GetUserByID(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if stored.DisplayName != "Original Name" {
		t.Errorf("stored display_name = %q, want untouched %q", stored.DisplayName, "Original Name")
	}
}
