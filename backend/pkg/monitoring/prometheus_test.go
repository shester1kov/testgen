package monitoring

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSetupPrometheus(t *testing.T) {
	t.Run("registers /metrics endpoint", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "test_service",
			Namespace:   "test",
			Subsystem:   "http",
		}

		SetupPrometheus(app, config)

		// Add a test endpoint
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Make a request first to generate metrics
		testReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		_, err := app.Test(testReq)
		if err != nil {
			t.Fatalf("failed to make test request: %v", err)
		}

		// Now check metrics endpoint
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		bodyStr := string(body)
		// Check that we get some metrics output (not empty)
		if len(bodyStr) == 0 {
			t.Error("expected non-empty metrics output")
		}
	})

	t.Run("collects http request metrics", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "test_service",
			Namespace:   "test",
			Subsystem:   "http",
		}

		SetupPrometheus(app, config)

		// Add a test endpoint
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Make a request to the test endpoint
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		_, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}

		// Check metrics
		metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		metricsResp, err := app.Test(metricsReq)
		if err != nil {
			t.Fatalf("failed to get metrics: %v", err)
		}
		defer metricsResp.Body.Close()

		body, err := io.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		// Check for http_requests_total metric
		if !strings.Contains(bodyStr, "http_requests_total") {
			t.Error("expected http_requests_total metric")
		}

		// Check for http_request_duration_seconds metric
		if !strings.Contains(bodyStr, "http_request_duration_seconds") {
			t.Error("expected http_request_duration_seconds metric")
		}

		// Check that our test endpoint was tracked
		if !strings.Contains(bodyStr, `path="/test"`) {
			t.Error("expected metrics for /test endpoint")
		}

		// Check that method was tracked
		if !strings.Contains(bodyStr, `method="GET"`) {
			t.Error("expected metrics for GET method")
		}

		// Check that status code was tracked
		if !strings.Contains(bodyStr, `status_code="200"`) {
			t.Error("expected metrics for 200 status code")
		}

		// Check that service name is in metrics
		if !strings.Contains(bodyStr, `service="test_service"`) {
			t.Error("expected service label in metrics")
		}
	})

	t.Run("tracks different http methods", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "test_service",
			Namespace:   "test",
			Subsystem:   "http",
		}

		SetupPrometheus(app, config)

		app.Post("/api/data", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "created"})
		})

		// Make POST request
		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		req.Header.Set("Content-Type", "application/json")
		_, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}

		// Check metrics
		metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		metricsResp, err := app.Test(metricsReq)
		if err != nil {
			t.Fatalf("failed to get metrics: %v", err)
		}
		defer metricsResp.Body.Close()

		body, err := io.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		if !strings.Contains(bodyStr, `method="POST"`) {
			t.Error("expected metrics for POST method")
		}

		if !strings.Contains(bodyStr, `path="/api/data"`) {
			t.Error("expected metrics for /api/data endpoint")
		}
	})

	t.Run("tracks different status codes", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "test_service",
			Namespace:   "test",
			Subsystem:   "http",
		}

		SetupPrometheus(app, config)

		app.Get("/not-found", func(c *fiber.Ctx) error {
			return c.Status(404).SendString("Not Found")
		})

		app.Get("/error", func(c *fiber.Ctx) error {
			return c.Status(500).SendString("Internal Error")
		})

		// Make 404 request
		req404 := httptest.NewRequest(http.MethodGet, "/not-found", nil)
		_, err := app.Test(req404)
		if err != nil {
			t.Fatalf("failed to make 404 request: %v", err)
		}

		// Make 500 request
		req500 := httptest.NewRequest(http.MethodGet, "/error", nil)
		_, err = app.Test(req500)
		if err != nil {
			t.Fatalf("failed to make 500 request: %v", err)
		}

		// Check metrics
		metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		metricsResp, err := app.Test(metricsReq)
		if err != nil {
			t.Fatalf("failed to get metrics: %v", err)
		}
		defer metricsResp.Body.Close()

		body, err := io.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		if !strings.Contains(bodyStr, `status_code="404"`) {
			t.Error("expected metrics for 404 status code")
		}

		if !strings.Contains(bodyStr, `status_code="500"`) {
			t.Error("expected metrics for 500 status code")
		}
	})

	t.Run("PrometheusConfig allows custom service name", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "custom_service_name",
			Namespace:   "custom",
			Subsystem:   "api",
		}

		SetupPrometheus(app, config)

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		_, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}

		// Check metrics
		metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		metricsResp, err := app.Test(metricsReq)
		if err != nil {
			t.Fatalf("failed to get metrics: %v", err)
		}
		defer metricsResp.Body.Close()

		body, err := io.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		if !strings.Contains(bodyStr, `service="custom_service_name"`) {
			t.Error("expected custom service name in metrics")
		}
	})
}

func TestPrometheusMetricsFormat(t *testing.T) {
	t.Run("metrics follow Prometheus text format", func(t *testing.T) {
		app := fiber.New()

		config := PrometheusConfig{
			ServiceName: "test",
			Namespace:   "test",
			Subsystem:   "http",
		}

		SetupPrometheus(app, config)

		// Add a test endpoint
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Make a request first to generate metrics
		testReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		_, err := app.Test(testReq)
		if err != nil {
			t.Fatalf("failed to make test request: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check Content-Type header
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/plain") && !strings.Contains(contentType, "text") {
			t.Errorf("expected text content type, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response: %v", err)
		}

		bodyStr := string(body)

		// Check that we have actual metrics output (not empty)
		if len(bodyStr) == 0 {
			t.Error("expected non-empty metrics output")
		}

		// Check that output is in text format with newlines
		lines := strings.Split(bodyStr, "\n")
		if len(lines) < 1 {
			t.Error("expected at least one line of metrics output")
		}
	})
}
