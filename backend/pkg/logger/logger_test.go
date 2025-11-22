package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
)

func TestNewCreatesWritableLogger(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "app.log")
	cfg := Config{Level: "debug", OutputPath: logFile, Format: "json"}

	l, err := New(cfg)
	if err != nil {
		t.Fatalf("expected logger to initialize, got error: %v", err)
	}

	l.Info("hello")

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("expected log file to be created: %v", err)
	}

	if !strings.Contains(string(data), "\"message\":\"hello\"") {
		t.Fatalf("log output missing message, got %q", string(data))
	}
}

func TestNewReturnsErrorForBadPath(t *testing.T) {
	badPath := filepath.Join(t.TempDir(), "nested", "cannot_create.log")
	cfg := Config{Level: "info", OutputPath: badPath, Format: "json"}

	if _, err := New(cfg); err == nil {
		t.Fatalf("expected error when parent directories are missing")
	}
}

func TestWithHelpersPreserveConfigAndFields(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "fields.log")
	cfg := Config{Level: "debug", OutputPath: logFile, Format: "json"}

	base, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create base logger: %v", err)
	}

	derived := base.WithField("user", 42).WithFields(map[string]interface{}{"role": "admin"}).WithError(os.ErrClosed).Context()
	derived.InfoWithFields("with fields", map[string]interface{}{"ok": true})

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("expected log file to be created: %v", err)
	}

	content := string(data)
	for _, token := range []string{"with fields", "user", "role", "ok"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected %q to be present in log output", token)
		}
	}

	if derived.config != cfg {
		t.Fatalf("logger config should be preserved after helper calls")
	}
}

func TestFactoryHelpers(t *testing.T) {
	if got := NewDefault(); got == nil {
		t.Fatal("expected default logger instance")
	}

	if prod, err := NewProduction(); err != nil || prod == nil {
		t.Fatalf("expected production logger, got err=%v", err)
	}

	if dev, err := NewDevelopment(); err != nil || dev == nil {
		t.Fatalf("expected development logger, got err=%v", err)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(RequestIDMiddleware())
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.Header.Get("X-Request-ID") == "" {
		t.Fatalf("expected request id header to be set")
	}
}

func TestHTTPMiddlewareLoggingBranches(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "http.log")
	logger, err := New(Config{Level: "debug", OutputPath: logFile, Format: "json"})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	app := fiber.New()
	app.Use(RequestIDMiddleware())
	app.Use(HTTPMiddleware(logger))

	app.Get("/ok", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/client", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusBadRequest) })
	app.Get("/server", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusInternalServerError) })
	app.Get("/err", func(c *fiber.Ctx) error { return fiber.ErrInternalServerError })

	for _, path := range []string{"/ok", "/client", "/server", "/err"} {
		req := httptest.NewRequest(fiber.MethodGet, path, nil)
		if _, err := app.Test(req, -1); err != nil {
			t.Fatalf("request to %s failed: %v", path, err)
		}
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("expected log file to be created: %v", err)
	}

	content := string(data)
	for _, expected := range []string{"completed", "client error", "server error", "HTTP request error"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q to be logged", expected)
		}
	}
}
