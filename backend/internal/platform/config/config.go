// Package config loads runtime configuration from the environment.
// Plain env vars with dev-friendly defaults, no config framework -
// see docs/ARCHITECTURE.md.
package config

import "os"

type Config struct {
	Env         string // "development" or "production"
	Port        string
	DatabaseURL string

	// Email (M3 verification). Dev defaults target Mailpit (docker-compose).
	SMTPAddr  string // host:port of the SMTP server
	EmailFrom string // From header on transactional mail
	AppURL    string // frontend base URL, for building verification links

	// Object storage (M5 in-chat images). Dev defaults target MinIO (docker-compose).
	S3Endpoint  string // host[:port], optionally scheme-prefixed (http:// -> insecure)
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func Load() Config {
	return Config{
		Env:         getenv("ENV", "development"),
		Port:        getenv("PORT", "8080"),
		// Dev Postgres publishes on 54320, off the default port, so dev tooling
		// can never collide with the prod stack's 127.0.0.1:5432 (#159).
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:dev@localhost:54320/bibseller?sslmode=disable"),
		SMTPAddr:    getenv("SMTP_ADDR", "localhost:1025"),
		EmailFrom:   getenv("EMAIL_FROM", "Bibseller <noreply@bibseller.dev>"),
		AppURL:      getenv("APP_URL", "http://localhost:5173"),
		S3Endpoint:  getenv("S3_ENDPOINT", "http://localhost:9000"),
		S3AccessKey: getenv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getenv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    getenv("S3_BUCKET", "bibseller-dev"),
	}
}

func (c Config) IsDev() bool { return c.Env == "development" }

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
