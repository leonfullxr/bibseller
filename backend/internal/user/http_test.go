package user_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
	"github.com/leonfullxr/bibseller/backend/internal/user"
)

// seedUser inserts a login-less account directly — used as the public-profile
// subject and as the "victim" an attacker tries (and fails) to rename.
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

// handler mounts the user routes plus auth — the PATCH gate resolves the
// session, so tests need a real registration flow to obtain a token.
func handler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)}, user.Routes(q), auth.Routes(q))
}

type session struct {
	token string
	user  auth.Account
}

// register creates a throwaway account (cleaned up via ON DELETE CASCADE) and
// returns its session token, the way a signed-in caller would hold one.
func register(t *testing.T, h http.Handler, pool *pgxpool.Pool) session {
	t.Helper()
	body := `{"email":"u-` + ids.New().String() + `@test.local",` +
		`"password":"correct horse battery staple","display_name":"Caller"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Token string       `json:"token"`
		User  auth.Account `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("register: bad JSON: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, resp.User.ID)
	})
	return session{token: resp.Token, user: resp.User}
}

// patch PATCHes a user, presenting token (if any) as the session cookie.
func patch(t *testing.T, h http.Handler, id uuid.UUID, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+id.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

func TestGetUserIsPublic(t *testing.T) {
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

func TestUpdateOwnDisplayName(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	me := register(t, h, pool)

	rec := patch(t, h, me.user.ID, `{"display_name": "  New Name  "}`, me.token)
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
	stored, err := sqlcgen.New(pool).GetUserByID(context.Background(), me.user.ID)
	if err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if stored.DisplayName != "New Name" {
		t.Errorf("stored display_name = %q, want %q", stored.DisplayName, "New Name")
	}
}

// The gate: no session is 401, and a valid session may not rename anyone else.
func TestUpdateDisplayNameIsOwnerOnly(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	me := register(t, h, pool)
	victim := seedUser(t, pool)

	if rec := patch(t, h, me.user.ID, `{"display_name":"Anyone"}`, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no session: status = %d, want 401", rec.Code)
	}
	if rec := patch(t, h, victim.ID, `{"display_name":"Pwned"}`, me.token); rec.Code != http.StatusForbidden {
		t.Errorf("cross-user: status = %d, want 403", rec.Code)
	}

	// The 403 must mean the write never happened.
	stored, err := sqlcgen.New(pool).GetUserByID(context.Background(), victim.ID)
	if err != nil {
		t.Fatalf("reload victim: %v", err)
	}
	if stored.DisplayName != "Original Name" {
		t.Errorf("victim renamed across the gate: %q", stored.DisplayName)
	}
}

func TestUpdateDisplayNameRejectsBadInput(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	me := register(t, h, pool)

	for name, body := range map[string]string{
		"not json":        `not json`,
		"missing field":   `{}`,
		"too short":       `{"display_name": "x"}`,
		"whitespace only": `{"display_name": "   "}`,
		"too long":        `{"display_name": "` + strings.Repeat("a", 51) + `"}`,
	} {
		if rec := patch(t, h, me.user.ID, body, me.token); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", name, rec.Code)
		}
	}

	// Nothing above may have touched the stored name.
	stored, err := sqlcgen.New(pool).GetUserByID(context.Background(), me.user.ID)
	if err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if stored.DisplayName != "Caller" {
		t.Errorf("stored display_name = %q, want untouched %q", stored.DisplayName, "Caller")
	}
}
