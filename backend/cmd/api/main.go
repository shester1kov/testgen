// @title Test Generation System API
// @version 1.0
// @description API for test generation system with LLM integration and Moodle export
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@testgen.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/shester1kov/testgen-backend/docs"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/router"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/logger"
	"github.com/shester1kov/testgen-backend/pkg/monitoring"
)

func main() {

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	appLogger, err := logger.New(logger.Config{
		Level:      cfg.Logger.Level,
		OutputPath: "stdout",
		Format:     cfg.Logger.Format,
	})
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer appLogger.Sync()

	appLogger.Info("Starting Test Generation System API",
		zap.String("version", "1.0"),
		zap.String("environment", cfg.Server.Environment),
	)

	// Initialize database
	db, err := postgres.NewDatabase(&postgres.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		Logger:   appLogger, // Pass logger for structured GORM logs
	})
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	appLogger.Info("Database connection established")

	container, err := InitializeApplication(cfg, db, appLogger)
	if err != nil {
		appLogger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Initialize JWT manager
	jwtManager := container.JWTManager

	// Initialize handlers
	authHandler := container.AuthHandler
	userHandler := container.UserHandler
	testHandler := container.TestHandler
	statsHandler := container.StatsHandler
	documentHandler := container.DocumentHandler
	redisClient := container.RedisClient

	// Run database seeders
	seeder := persistence.NewSeeder(container.UserRepo, container.RoleRepo, cfg, appLogger)
	if err := seeder.Seed(context.Background()); err != nil {
		appLogger.Error("Failed to run database seeders", zap.Error(err))
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Setup Prometheus metrics (should be registered early)
	monitoring.SetupPrometheus(app, monitoring.PrometheusConfig{
		ServiceName: "testgen_api",
		Namespace:   "testgen",
		Subsystem:   "http",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.RequestIDMiddleware())
	app.Use(logger.HTTPMiddleware(appLogger))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5173,http://localhost,http://109.73.195.85,http://petproj.ru.net,https://petproj.ru.net",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		ExposeHeaders:    "Content-Disposition",
		AllowCredentials: true,
	}))

	// Health check endpoint
	// @Summary Health check
	// @Description Check if the API is running and healthy
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Test Generation System API is running",
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup routes
	router.SetupRoutes(app, authHandler, userHandler, documentHandler, testHandler, statsHandler, jwtManager, cfg.Cookie.Name, container.ActivityLogger)

	// Root endpoint
	// @Summary API version information
	// @Description Get API version and available endpoints
	// @Tags info
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /api/v1 [get]
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Test Generation System API v1",
			"endpoints": fiber.Map{
				"auth":      "/api/v1/auth",
				"documents": "/api/v1/documents",
				"tests":     "/api/v1/tests",
				"stats":     "/api/v1/stats",
			},
		})
	})

	// Start server
	port := cfg.Server.Port

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	appLogger.Info("Server started successfully",
		zap.String("port", port),
		zap.String("environment", cfg.Server.Environment),
	)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		appLogger.Error("Failed to close Redis connection", zap.Error(err))
	}

	appLogger.Info("Server exited successfully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}
