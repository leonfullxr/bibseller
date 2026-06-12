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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(logger *slog.Logger, pool *pgxpool.Pool, apiV1 ...func(chi.Router)) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		// Liveness: the process is up.
		r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})
		// Readiness: dependencies are reachable.
		r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
			defer cancel()
			if err := pool.Ping(ctx); err != nil {
				JSON(w, http.StatusServiceUnavailable,
					map[string]string{"status": "unavailable", "reason": "database unreachable"})
				return
			}
			JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		r.Route("/v1", func(v1 chi.Router) {
			for _, mount := range apiV1 {
				mount(v1)
			}
		})
	})

	return r
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
// → API conventions).
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}
