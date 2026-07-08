// Package config loads runtime configuration from the environment.
// Plain env vars with dev-friendly defaults, no config framework -
// see docs/ARCHITECTURE.md.
package config

import (
	"cmp"
	"os"
)

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

	// TrustProxyHeader: trust CF-Connecting-IP as the client address (#182).
	// Set "1" only where every request provably traverses the Cloudflare edge
	// (the prod compose stack); default off keys rate limits and session audit
	// on the unforgeable RemoteAddr.
	TrustProxyHeader bool
}

func Load() Config {
	return Config{
		Env:  cmp.Or(os.Getenv("ENV"), "development"),
		Port: cmp.Or(os.Getenv("PORT"), "8080"),
		// Dev Postgres publishes on 54320, off the default port, so dev tooling
		// can never collide with the prod stack's 127.0.0.1:5432 (#159).
		DatabaseURL: cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres:dev@localhost:54320/bibseller?sslmode=disable"),
		SMTPAddr:    cmp.Or(os.Getenv("SMTP_ADDR"), "localhost:1025"),
		EmailFrom:   cmp.Or(os.Getenv("EMAIL_FROM"), "Bibseller <noreply@bibseller.dev>"),
		AppURL:      cmp.Or(os.Getenv("APP_URL"), "http://localhost:5173"),
		S3Endpoint:  cmp.Or(os.Getenv("S3_ENDPOINT"), "http://localhost:9000"),
		S3AccessKey: cmp.Or(os.Getenv("S3_ACCESS_KEY"), "minioadmin"),
		S3SecretKey: cmp.Or(os.Getenv("S3_SECRET_KEY"), "minioadmin"),
		S3Bucket:    cmp.Or(os.Getenv("S3_BUCKET"), "bibseller-dev"),

		TrustProxyHeader: os.Getenv("TRUST_PROXY_HEADER") == "1",
	}
}

func (c Config) IsDev() bool { return c.Env == "development" }
