package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

// noopMailer satisfies auth.Mailer without sending anything - the verification
// flow is exercised by seeding a token row directly, not by parsing email.
type noopMailer struct{}

func (noopMailer) SendVerification(_, _ string) error { return nil }

func handler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)}, auth.Routes(q, noopMailer{}, "http://test.local"))
}

type sessionResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	User      auth.Account `json:"user"`
}

// post sends a JSON body; token (if non-empty) is presented the way real
// callers present it: as the __Host-session cookie.
func post(t *testing.T, h http.Handler, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

func getMe(t *testing.T, h http.Handler, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

// register creates a throwaway account (unique email per call) and cleans it
// up; sessions go with it via ON DELETE CASCADE.
func register(t *testing.T, h http.Handler, pool *pgxpool.Pool, password string) sessionResponse {
	t.Helper()
	email := "t-" + ids.New().String() + "@test.local"
	body := `{"email":"` + email + `","password":"` + password + `","display_name":"Auth Tester"}`
	rec := post(t, h, "/api/v1/auth/register", body, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp sessionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("register: bad JSON: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, resp.User.ID)
	})
	return resp
}

func TestSessionLifecycle(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)

	// Register: account + first session.
	reg := register(t, h, pool, "correct horse battery staple")
	if reg.Token == "" || reg.ExpiresAt == "" {
		t.Fatal("register response missing token or expires_at")
	}

	// The database must hold SHA-256(token), never the token itself.
	sum := sha256.Sum256([]byte(reg.Token))
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM sessions WHERE token_hash = $1`, sum[:]).Scan(&n); err != nil || n != 1 {
		t.Fatalf("hashed session row: n = %d, err = %v", n, err)
	}
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM sessions WHERE token_hash = $1`, []byte(reg.Token)).Scan(&n); err != nil || n != 0 {
		t.Fatalf("raw token stored in DB: n = %d, err = %v", n, err)
	}

	// The session authenticates /auth/me.
	rec := getMe(t, h, reg.Token)
	if rec.Code != http.StatusOK {
		t.Fatalf("me: status = %d, body = %s", rec.Code, rec.Body)
	}
	var me auth.Account
	if err := json.Unmarshal(rec.Body.Bytes(), &me); err != nil || me.ID != reg.User.ID {
		t.Fatalf("me: got %+v, err = %v", me, err)
	}

	// Login presenting the old session: new token issued, old one rotated out.
	body := `{"email":"` + reg.User.Email + `","password":"correct horse battery staple"}`
	rec = post(t, h, "/api/v1/auth/login", body, reg.Token)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: status = %d, body = %s", rec.Code, rec.Body)
	}
	var login sessionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &login); err != nil {
		t.Fatalf("login: bad JSON: %v", err)
	}
	if login.Token == reg.Token {
		t.Fatal("login reused the old session token (no rotation)")
	}
	if rec := getMe(t, h, reg.Token); rec.Code != http.StatusUnauthorized {
		t.Errorf("rotated-out token still works: status = %d, want 401", rec.Code)
	}
	if rec := getMe(t, h, login.Token); rec.Code != http.StatusOK {
		t.Errorf("fresh login token rejected: status = %d", rec.Code)
	}

	// Logout kills the session; a second logout is still a success (idempotent).
	if rec := post(t, h, "/api/v1/auth/logout", "", login.Token); rec.Code != http.StatusNoContent {
		t.Fatalf("logout: status = %d", rec.Code)
	}
	if rec := getMe(t, h, login.Token); rec.Code != http.StatusUnauthorized {
		t.Errorf("token survives logout: status = %d, want 401", rec.Code)
	}
	if rec := post(t, h, "/api/v1/auth/logout", "", login.Token); rec.Code != http.StatusNoContent {
		t.Errorf("repeated logout: status = %d, want 204", rec.Code)
	}
}

func TestEmailVerification(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	q := sqlcgen.New(pool)
	reg := register(t, h, pool, "correct horse battery staple")

	// A fresh account is unverified.
	if me := getMe(t, h, reg.Token); me.Code == http.StatusOK {
		var acc auth.Account
		_ = json.Unmarshal(me.Body.Bytes(), &acc)
		if acc.EmailVerified {
			t.Fatal("freshly registered account is already email_verified")
		}
	}

	// A bogus token is rejected.
	if rec := post(t, h, "/api/v1/auth/verify", `{"token":"not-a-real-token"}`, ""); rec.Code != http.StatusBadRequest {
		t.Errorf("bogus token: status = %d, want 400", rec.Code)
	}

	// Seed a known token for this user (we control the raw value; the DB only
	// ever holds its hash), then consume it.
	token := "verify-" + ids.New().String()
	sum := sha256.Sum256([]byte(token))
	if _, err := q.CreateEmailVerification(context.Background(), sqlcgen.CreateEmailVerificationParams{
		TokenHash: sum[:], UserID: reg.User.ID,
	}); err != nil {
		t.Fatalf("seed verification token: %v", err)
	}

	if rec := post(t, h, "/api/v1/auth/verify", `{"token":"`+token+`"}`, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("verify: status = %d, body = %s", rec.Code, rec.Body)
	}

	// The account is now verified, and /auth/me reflects it.
	me := getMe(t, h, reg.Token)
	var acc auth.Account
	if err := json.Unmarshal(me.Body.Bytes(), &acc); err != nil || !acc.EmailVerified {
		t.Errorf("me.email_verified = %v after verify (err=%v)", acc.EmailVerified, err)
	}

	// The token is single-use: replaying it fails.
	if rec := post(t, h, "/api/v1/auth/verify", `{"token":"`+token+`"}`, ""); rec.Code != http.StatusBadRequest {
		t.Errorf("replayed token: status = %d, want 400", rec.Code)
	}
}

func TestLoginFailuresAreIndistinguishable(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)
	reg := register(t, h, pool, "correct horse battery staple")

	wrongPw := post(t, h, "/api/v1/auth/login",
		`{"email":"`+reg.User.Email+`","password":"not the password"}`, "")
	unknown := post(t, h, "/api/v1/auth/login",
		`{"email":"nobody-`+ids.New().String()+`@test.local","password":"whatever123"}`, "")

	for name, rec := range map[string]*httptest.ResponseRecorder{
		"wrong password": wrongPw, "unknown email": unknown,
	} {
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s: status = %d, want 401", name, rec.Code)
		}
	}
	// Same body for both: the response may not reveal which part was wrong.
	if wrongPw.Body.String() != unknown.Body.String() {
		t.Errorf("login failures distinguishable:\n%s\n%s", wrongPw.Body, unknown.Body)
	}
}

func TestRegisterValidation(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)

	for name, body := range map[string]string{
		"bad email":       `{"email":"not-an-email","password":"longenough","display_name":"Tester"}`,
		"short password":  `{"email":"a@test.local","password":"short","display_name":"Tester"}`,
		"short name":      `{"email":"a@test.local","password":"longenough","display_name":"x"}`,
		"not json":        `nope`,
		"address w/ name": `{"email":"Bob <bob@test.local>","password":"longenough","display_name":"Tester"}`,
	} {
		if rec := post(t, h, "/api/v1/auth/register", body, ""); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", name, rec.Code)
		}
	}

	// Duplicate email -> 409, case-insensitively (citext).
	reg := register(t, h, pool, "correct horse battery staple")
	dup := `{"email":"` + strings.ToUpper(reg.User.Email) + `","password":"longenough","display_name":"Copycat"}`
	if rec := post(t, h, "/api/v1/auth/register", dup, ""); rec.Code != http.StatusConflict {
		t.Errorf("duplicate email: status = %d, want 409", rec.Code)
	}
}

func TestCrossSiteMutationBlocked(t *testing.T) {
	pool := testdb.Pool(t)
	h := handler(pool)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login",
		strings.NewReader(`{"email":"a@test.local","password":"whatever123"}`))
	req.Header.Set("Sec-Fetch-Site", "cross-site") // what a browser sends for another origin's request
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("cross-site POST: status = %d, want 403", rec.Code)
	}
}
