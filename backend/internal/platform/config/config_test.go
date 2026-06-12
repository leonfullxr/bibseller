package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()
	if cfg.Env != "development" {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if !cfg.IsDev() {
		t.Error("IsDev() = false, want true")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("PORT", "9999")
	t.Setenv("DATABASE_URL", "postgres://example/app")

	cfg := Load()
	if cfg.Env != "production" {
		t.Errorf("Env = %q, want production", cfg.Env)
	}
	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want 9999", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://example/app" {
		t.Errorf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if cfg.IsDev() {
		t.Error("IsDev() = true, want false")
	}
}
