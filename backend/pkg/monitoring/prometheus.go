package monitoring

import (
	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// PrometheusConfig holds configuration for Prometheus metrics
type PrometheusConfig struct {
	ServiceName string
	Namespace   string
	Subsystem   string
}

// SetupPrometheus configures and registers Prometheus middleware
func SetupPrometheus(app *fiber.App, config PrometheusConfig) {
	// Create Prometheus middleware with service name
	prometheus := fiberprometheus.New(config.ServiceName)

	// Note: The fiberprometheus library automatically sets up metrics
	// with the service name. Namespace and Subsystem are included
	// in the config struct for future extensibility or custom metrics.

	// Register middleware to collect HTTP metrics
	// This should be registered early in the middleware chain
	prometheus.RegisterAt(app, "/metrics")

	// Use middleware to track all requests
	app.Use(prometheus.Middleware)
}
