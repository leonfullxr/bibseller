package auth

import (
	"context"
	"net/http"

	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
)

type ctxKey int

const userKey ctxKey = iota

// ResolveUser is middleware that resolves the session cookie to its user once
// per request and stashes it in the context, so handlers gate via
// UserFromContext instead of each re-running the lookup. Best-effort: a
// missing, unknown, or expired token simply leaves no user in context (the
// handler answers 401). Anonymous requests carry no cookie and skip the DB
// entirely (see Authenticate -> requestToken).
func ResolveUser(q *sqlcgen.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if row, ok := Authenticate(r.Context(), q, r); ok {
				r = r.WithContext(context.WithValue(r.Context(), userKey, row))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// UserFromContext returns the signed-in user resolved by ResolveUser, or
// ok=false if the request carried no valid session.
func UserFromContext(ctx context.Context) (sqlcgen.GetSessionWithUserRow, bool) {
	row, ok := ctx.Value(userKey).(sqlcgen.GetSessionWithUserRow)
	return row, ok
}
