package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type ctxKey int

const requestIDKey ctxKey = iota

// RequestID returns the id attached by the requestID middleware, or "".
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// requestID tags every request with a short random id for log correlation,
// honoring an inbound X-Request-Id (e.g. from a reverse proxy) when present.
//
// SECURITY: this trusts X-Request-Id, which is correct only behind our own
// proxy. If this port ever becomes directly client-facing, a client can forge
// correlation ids — sanitize (length/charset cap) or stop honoring the header
// at that point (tracked: near-term considerations, docs/TECH_NOTES.md).
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			var b [8]byte
			_, _ = rand.Read(b[:])
			id = hex.EncodeToString(b[:])
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

// statusWriter records the status and bytes written for the request log.
// It deliberately drops optional interfaces (Flusher, Hijacker): the API is
// plain JSON over HTTP. Revisit if chat upgrades to SSE (docs/CONTEXT.md D13).
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

// requestLogger emits one structured line per request.
func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(sw, r)
			status := sw.status
			if status == 0 {
				status = http.StatusOK
			}
			logger.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"bytes", sw.bytes,
				"dur_ms", time.Since(start).Milliseconds(),
				"request_id", RequestID(r.Context()),
			)
		})
	}
}

// csrfGuard rejects browser-issued cross-site mutations using fetch metadata
// (docs/ARCHITECTURE.md → Auth & sessions). Browsers attach Sec-Fetch-Site to
// every request they make: "cross-site" means another origin's page triggered
// it, which is never legitimate for a state change here. Server-to-server
// callers (the SvelteKit actions) and curl send no such header and pass —
// CSRF is only about riding the *browser's* ambient cookie; a client that
// sets its own headers holds the credential anyway. This is defense in depth
// on top of SameSite=Lax, which already keeps the session cookie off
// cross-site POSTs.
func csrfGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r) // safe methods must not mutate; nothing to guard
			return
		}
		if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			Error(w, http.StatusForbidden, "cross_site_request", "cross-site requests are not allowed")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// recoverer converts handler panics into 500s instead of dropped connections.
func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// net/http's sentinel for cleanly aborting a response.
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					logger.Error("panic recovered",
						"err", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"request_id", RequestID(r.Context()),
						"stack", string(debug.Stack()),
					)
					Error(w, http.StatusInternalServerError, "internal", "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
