package moderation_test

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
	"github.com/leonfullxr/bibseller/backend/internal/moderation"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/testdb"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/ids"
)

type noopMailer struct{}

func (noopMailer) SendVerification(_, _, _ string) error  { return nil }
func (noopMailer) SendPasswordReset(_, _, _ string) error { return nil }

func authedHandler(pool *pgxpool.Pool) http.Handler {
	q := sqlcgen.New(pool)
	return httpx.NewRouter(slog.New(slog.DiscardHandler), pool,
		[]httpx.Middleware{auth.ResolveUser(q)},
		moderation.Routes(q), auth.Routes(pool, noopMailer{}, "http://test.local"))
}

func doJSON(t *testing.T, h http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	h.ServeHTTP(rec, req)
	return rec
}

func register(t *testing.T, h http.Handler, pool *pgxpool.Pool, verified bool) (string, uuid.UUID) {
	t.Helper()
	email := "m-" + ids.New().String() + "@test.local"
	body := `{"email":"` + email + `","password":"correct horse battery staple","display_name":"Mod User"}`
	rec := doJSON(t, h, http.MethodPost, "/api/v1/auth/register", body, "")
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
	if verified {
		if err := sqlcgen.New(pool).MarkEmailVerified(context.Background(), resp.User.ID); err != nil {
			t.Fatalf("verify: %v", err)
		}
	}
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM reports WHERE reporter_id = $1`, resp.User.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM blocks WHERE blocker_id = $1 OR blocked_id = $1`, resp.User.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, resp.User.ID)
	})
	return resp.Token, resp.User.ID
}

func TestCreateReport(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	tok, uid := register(t, h, pool, true)
	subject := ids.New().String()

	rec := doJSON(t, h, http.MethodPost, "/api/v1/reports",
		`{"subject_type":"listing","subject_id":"`+subject+`","reason":"scam","details":"looks fake"}`, tok)
	if rec.Code != http.StatusCreated {
		t.Fatalf("report: status = %d, body = %s", rec.Code, rec.Body)
	}
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM reports WHERE reporter_id = $1 AND subject_type = 'listing' AND reason = 'scam' AND status = 'open'`,
		uid).Scan(&n); err != nil || n != 1 {
		t.Fatalf("report row: n = %d, err = %v", n, err)
	}

	// Constraint violations are clean 400s, not 500s.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/reports",
		`{"subject_type":"listing","subject_id":"`+subject+`","reason":"nope"}`, tok); rec.Code != http.StatusBadRequest {
		t.Errorf("bad reason: status = %d, want 400", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/reports",
		`{"subject_type":"planet","subject_id":"`+subject+`","reason":"scam"}`, tok); rec.Code != http.StatusBadRequest {
		t.Errorf("bad subject_type: status = %d, want 400", rec.Code)
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/reports",
		`{"subject_type":"listing","subject_id":"`+subject+`","reason":"scam"}`, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no session: status = %d, want 401", rec.Code)
	}
	unv, _ := register(t, h, pool, false)
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/reports",
		`{"subject_type":"listing","subject_id":"`+subject+`","reason":"scam"}`, unv); rec.Code != http.StatusForbidden {
		t.Errorf("unverified: status = %d, want 403", rec.Code)
	}
}

func TestBlockUnblock(t *testing.T) {
	pool := testdb.Pool(t)
	h := authedHandler(pool)
	tok, uid := register(t, h, pool, true)
	_, otherID := register(t, h, pool, true)

	count := func() int {
		var n int
		if err := pool.QueryRow(context.Background(),
			`SELECT count(*) FROM blocks WHERE blocker_id = $1 AND blocked_id = $2`, uid, otherID).Scan(&n); err != nil {
			t.Fatalf("count blocks: %v", err)
		}
		return n
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/blocks", `{"blocked_id":"`+otherID.String()+`"}`, tok); rec.Code != http.StatusNoContent {
		t.Fatalf("block: status = %d, body = %s", rec.Code, rec.Body)
	}
	if count() != 1 {
		t.Fatal("block row not stored")
	}
	// Re-blocking is idempotent.
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/blocks", `{"blocked_id":"`+otherID.String()+`"}`, tok); rec.Code != http.StatusNoContent {
		t.Errorf("re-block: status = %d, want 204", rec.Code)
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/blocks", `{"blocked_id":"`+uid.String()+`"}`, tok); rec.Code != http.StatusBadRequest {
		t.Errorf("block self: status = %d, want 400", rec.Code)
	}
	if rec := doJSON(t, h, http.MethodPost, "/api/v1/blocks", `{"blocked_id":"`+ids.New().String()+`"}`, tok); rec.Code != http.StatusNotFound {
		t.Errorf("block unknown user: status = %d, want 404", rec.Code)
	}

	if rec := doJSON(t, h, http.MethodDelete, "/api/v1/blocks/"+otherID.String(), "", tok); rec.Code != http.StatusNoContent {
		t.Fatalf("unblock: status = %d, body = %s", rec.Code, rec.Body)
	}
	if count() != 0 {
		t.Error("block row not removed after unblock")
	}

	if rec := doJSON(t, h, http.MethodPost, "/api/v1/blocks", `{"blocked_id":"`+otherID.String()+`"}`, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no session: status = %d, want 401", rec.Code)
	}
}
