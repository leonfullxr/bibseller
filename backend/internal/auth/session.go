package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
)

// CookieName is the browser-side session cookie. The __Host- prefix is
// enforced BY THE BROWSER, not by us: a browser refuses to store it unless
// the cookie is Secure, has Path=/, and has no Domain attribute. That makes
// it impossible for a subdomain (or a non-HTTPS man-in-the-middle) to plant
// or overwrite our session cookie.
//
// Division of labor (docs/ARCHITECTURE.md -> Auth & sessions): the Go API
// MINTS tokens and READS this cookie, but never sets it - Set-Cookie on a
// server-to-server response would never reach the browser. The SvelteKit
// form action receives the raw token in the JSON response (server-to-server
// only) and sets the cookie with these exact attributes.
const CookieName = "__Host-session"

// tokenBytes = 256 bits from crypto/rand. The token is a pure random
// capability - there is nothing to "crack"; the only attack is guessing,
// and 2^256 is not guessable. Encoded as unpadded base64url (43 chars,
// cookie-safe alphabet).
const tokenBytes = 32

// touchInterval throttles the sliding-expiry write: the 30-day idle window
// (the TTL itself lives in db/queries/auth.sql) only needs refreshing once
// per hour, not on every request that validates the session.
const touchInterval = time.Hour

// newToken mints a session: the raw token goes to the client, only the
// SHA-256 of it goes to the database.
func newToken() (token string, tokenHash []byte, err error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", nil, err
	}
	token = base64.RawURLEncoding.EncodeToString(buf)
	return token, hashToken(token), nil
}

// hashToken is what makes a database leak harmless: the sessions table holds
// SHA-256(token), and SHA-256 cannot be reversed to the token the browser
// presents. A *fast* hash is correct here - unlike passwords, tokens are
// 256-bit random values with no dictionary to attack, so argon2 would add
// cost without adding security. (This also means we never store anything a
// log line or backup could turn into a working credential.)
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

// issueSession mints a token and persists its hash, with client metadata
// (ip, user agent) for a future "active sessions" view.
func issueSession(ctx context.Context, q *sqlcgen.Queries, userID uuid.UUID, r *http.Request) (token string, expiresAt time.Time, err error) {
	token, tokenHash, err := newToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt, err = q.CreateSession(ctx, sqlcgen.CreateSessionParams{
		TokenHash: tokenHash,
		UserID:    userID,
		Ip:        clientAddr(r),
		UserAgent: clientUA(r),
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

// requestToken extracts the session token presented by the caller - the
// __Host-session cookie, forwarded verbatim by the SvelteKit server or sent
// automatically by the browser on direct same-origin API calls.
func requestToken(r *http.Request) (string, bool) {
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return "", false
	}
	return c.Value, true
}

// authenticate resolves a presented token to its user, enforcing expiry and
// sliding the 30-day idle window (throttled by touchInterval). Returns
// ok=false for missing, unknown, or expired tokens - the caller answers 401.
// Internal to the package: ResolveUser runs it once per request as v1
// middleware, and handlers read the result via UserFromContext.
func authenticate(ctx context.Context, q *sqlcgen.Queries, r *http.Request) (sqlcgen.GetSessionWithUserRow, bool) {
	token, ok := requestToken(r)
	if !ok {
		return sqlcgen.GetSessionWithUserRow{}, false
	}
	row, err := q.GetSessionWithUser(ctx, hashToken(token))
	if err != nil {
		// pgx.ErrNoRows = unknown or expired token; anything else is a DB
		// problem, and failing closed (401) is the safe behavior for auth.
		return sqlcgen.GetSessionWithUserRow{}, false
	}
	if time.Since(row.LastSeenAt) > touchInterval {
		_ = q.TouchSession(ctx, hashToken(token)) // best-effort; expiry slides next time
	}
	return row, true
}

func clientAddr(r *http.Request) *netip.Addr {
	// httpx.ClientIP returns a bare address (no port) in both of its shapes.
	addr, err := netip.ParseAddr(httpx.ClientIP(r))
	if err != nil {
		return nil
	}
	// Unmap 4-in-6 so the stored inet matches how ClientIPKey buckets the
	// same client (::ffff:a.b.c.d and a.b.c.d are different inet values).
	addr = addr.Unmap()
	return &addr
}

func clientUA(r *http.Request) *string {
	if ua := r.UserAgent(); ua != "" {
		return &ua
	}
	return nil
}
