package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Errorf("Failed to load configuration: %v", err)
	}

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

	cfg, err := Load()
	if err != nil {
		t.Errorf("Failed to load configuration: %v", err)
	}

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

func TestValidate(t *testing.T) {
	// Helper function to set all required env vars with valid values
	setValidEnv := func() {
		t.Setenv("ENV", "development") // Explicitly set to development
		t.Setenv("JWT_SECRET", "some-secure-secret-key-at-least-32-chars-long")
		t.Setenv("DB_HOST", "localhost")
		t.Setenv("DB_NAME", "testgen_db")
		t.Setenv("MINIO_ENDPOINT", "minio.example.com:9000") // Non-default value
		t.Setenv("MINIO_ACCESS_KEY", "testkey")
		t.Setenv("MINIO_SECRET_KEY", "testsecret")
		t.Setenv("MINIO_BUCKET_NAME", "test-bucket")
	}

	tests := []struct {
		name        string
		setupEnv    func()
		expectError bool
	}{
		{
			name: "valid development configuration",
			setupEnv: func() {
				setValidEnv()
			},
			expectError: false,
		},
		{
			name: "forbidden JWT secret",
			setupEnv: func() {
				setValidEnv()
				t.Setenv("JWT_SECRET", "your-secret-key") // Forbidden default
			},
			expectError: true,
		},
		{
			name: "production with default JWT secret",
			setupEnv: func() {
				setValidEnv()
				t.Setenv("ENV", "production")
				t.Setenv("JWT_SECRET", "your-secret-key-change-in-production") // Default forbidden in prod
			},
			expectError: true,
		},
		{
			name: "production with minioadmin credentials",
			setupEnv: func() {
				setValidEnv()
				t.Setenv("ENV", "production")
				t.Setenv("MINIO_ACCESS_KEY", "minioadmin") // Forbidden in production
			},
			expectError: true,
		},
		{
			name: "development with default credentials is OK",
			setupEnv: func() {
				setValidEnv()
				t.Setenv("ENV", "development")
				t.Setenv("JWT_SECRET", "your-secret-key-change-in-production") // OK in dev
				t.Setenv("MINIO_ACCESS_KEY", "minioadmin")                     // OK in dev
				t.Setenv("MINIO_SECRET_KEY", "minioadmin")                     // OK in dev
			},
			expectError: false,
		},
		{
			name: "missing database configuration",
			setupEnv: func() {
				// Don't call setValidEnv, start from scratch
				t.Setenv("ENV", "development")
				t.Setenv("JWT_SECRET", "valid-secret-key-here")
				t.Setenv("DB_HOST", "")   // Will use default "localhost"
				t.Setenv("DB_NAME", "")   // Will use default "testgen_db"
				t.Setenv("MINIO_ENDPOINT", "localhost:9000")
				t.Setenv("MINIO_ACCESS_KEY", "key")
				t.Setenv("MINIO_SECRET_KEY", "secret")
				t.Setenv("MINIO_BUCKET_NAME", "bucket")
			},
			expectError: false, // Defaults are valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			cfg, err := Load()

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && cfg == nil {
				t.Errorf("expected valid config but got nil")
			}
		})
	}
}
