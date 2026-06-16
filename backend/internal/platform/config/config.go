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
}

func Load() Config {
	return Config{
		Env:         getenv("ENV", "development"),
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:dev@localhost:5432/bibseller?sslmode=disable"),
		SMTPAddr:    getenv("SMTP_ADDR", "localhost:1025"),
		EmailFrom:   getenv("EMAIL_FROM", "Bibseller <noreply@bibseller.dev>"),
		AppURL:      getenv("APP_URL", "http://localhost:5173"),
	}
}

func (c Config) IsDev() bool { return c.Env == "development" }

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
