package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()

	if cfg.Server.Port == "" || cfg.JWT.Secret == "" {
		t.Fatalf("defaults not loaded correctly")
	}
	if cfg.File.MaxFileSize <= 0 {
		t.Fatalf("expected positive max file size")
	}
}

func TestEnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9999")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("MAX_FILE_SIZE", "123")
	t.Setenv("ENABLE_METRICS", "false")

	cfg := Load()

	if cfg.Server.Port != "9999" {
		t.Fatalf("expected port override, got %s", cfg.Server.Port)
	}
	if !cfg.Cookie.Secure {
		t.Fatalf("expected secure cookie override")
	}
	if cfg.File.MaxFileSize != 123 {
		t.Fatalf("expected max file size override, got %d", cfg.File.MaxFileSize)
	}
	if cfg.Server.EnableMetrics {
		t.Fatalf("expected metrics disabled")
	}
}

func TestGetEnvHelpers(t *testing.T) {
	if val := getEnv("MISSING", "fallback"); val != "fallback" {
		t.Fatalf("expected fallback value")
	}

	t.Setenv("BOOL_VALUE", "notabool")
	if val := getEnvBool("BOOL_VALUE", true); !val {
		t.Fatalf("invalid bool should return default")
	}

	t.Setenv("INT_VALUE", "notanint")
	if val := getEnvInt64("INT_VALUE", 42); val != 42 {
		t.Fatalf("invalid int should return default")
	}
}
