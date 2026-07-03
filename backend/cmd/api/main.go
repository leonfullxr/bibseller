package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leonfullxr/bibseller/backend/internal/auth"
	"github.com/leonfullxr/bibseller/backend/internal/chat"
	"github.com/leonfullxr/bibseller/backend/internal/listing"
	"github.com/leonfullxr/bibseller/backend/internal/moderation"
	"github.com/leonfullxr/bibseller/backend/internal/platform/config"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db"
	"github.com/leonfullxr/bibseller/backend/internal/platform/db/sqlcgen"
	"github.com/leonfullxr/bibseller/backend/internal/platform/email"
	"github.com/leonfullxr/bibseller/backend/internal/platform/httpx"
	"github.com/leonfullxr/bibseller/backend/internal/platform/storage"
	"github.com/leonfullxr/bibseller/backend/internal/race"
	"github.com/leonfullxr/bibseller/backend/internal/user"
)

func main() {
	// Container healthcheck mode (distroless image has no shell/curl): probe
	// the serving process's readiness endpoint and exit 0/1.
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		os.Exit(healthcheck("http://localhost:" + config.Load().Port))
	}
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()
	logger := newLogger(cfg)
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Lazy pool: boots even while Postgres is still starting; /api/readyz
	// reports actual readiness.
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := sqlcgen.New(pool)
	mailer := email.SMTPMailer{Addr: cfg.SMTPAddr, From: cfg.EmailFrom}

	store, err := storage.New(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	if err != nil {
		return err
	}
	// Ensure the private image bucket exists. Best-effort so the API still boots
	// if the object store is briefly unavailable; the image routes error until it
	// returns.
	bucketCtx, cancelBucket := context.WithTimeout(ctx, 5*time.Second)
	if err := store.EnsureBucket(bucketCtx); err != nil {
		logger.Warn("object storage not ready", "err", err)
	}
	cancelBucket()

	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: httpx.NewRouter(logger, pool,
			[]httpx.Middleware{auth.RateLimit(), auth.ResolveUser(queries)},
			race.Routes(queries), listing.Routes(queries), user.Routes(queries),
			auth.Routes(pool, mailer, cfg.AppURL), chat.Routes(pool, mailer, store, cfg.AppURL),
			moderation.Routes(queries)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	// Background jobs run on every instance and coordinate via Postgres advisory
	// locks; they stop when ctx is cancelled on shutdown.
	go listing.StartExpiryJob(ctx, pool, logger, time.Hour)
	go chat.StartRetentionJob(ctx, pool, store, logger, 24*time.Hour)

	errc := make(chan error, 1)
	go func() {
		logger.Info("api listening", "addr", srv.Addr, "env", cfg.Env)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func newLogger(cfg config.Config) *slog.Logger {
	if cfg.IsDev() {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
