// Package httpx holds the HTTP router, shared middleware, response helpers,
// and cross-domain endpoints (health checks). Domain packages register their
// routes under /api/v1 via the mount functions passed to NewRouter.
package httpx

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Middleware wraps a handler. v1Middleware (e.g. session resolution, rate
// limiting) is applied only to the /api/v1 sub-mux, never to the health checks.
type Middleware = func(http.Handler) http.Handler

func NewRouter(logger *slog.Logger, pool *pgxpool.Pool, v1Middleware []Middleware, apiV1 ...func(*http.ServeMux)) http.Handler {
	mux := http.NewServeMux()

	// Liveness: the process is up.
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, _ *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	// Readiness: dependencies are reachable.
	mux.HandleFunc("GET /api/readyz", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			JSON(w, http.StatusServiceUnavailable,
				map[string]string{"status": "unavailable", "reason": "database unreachable"})
			return
		}
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Domain routes register method+pattern routes relative to /api/v1;
	// StripPrefix rebases incoming paths onto that sub-mux. v1 middleware sees
	// the rebased path ("/auth/login") and wraps the mux with v1Middleware[0]
	// outermost.
	v1 := http.NewServeMux()
	for _, mount := range apiV1 {
		mount(v1)
	}
	var v1h http.Handler = v1
	for i := len(v1Middleware) - 1; i >= 0; i-- {
		v1h = v1Middleware[i](v1h)
	}
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1h))

	// Outermost first: tag the request, log it, convert panics to 500s,
	// block cross-site mutations.
	var h http.Handler = mux
	h = csrfGuard(h)
	h = recoverer(logger)(h)
	h = requestLogger(logger)(h)
	h = requestID(h)
	return h
}

// JSON writes v as a JSON response with the given status.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error writes the API's standard error envelope: stable code strings the
// frontend can switch on, human message for debugging (docs/ARCHITECTURE.md
// -> API conventions).
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}
