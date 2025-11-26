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
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "github.com/shester1kov/testgen-backend/docs"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/router"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/logger"
	"github.com/shester1kov/testgen-backend/pkg/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Silently continue if .env file not found
	}

	// Load configuration
	cfg := config.Load()

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
	})
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	appLogger.Info("Database connection established")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	documentRepo := postgres.NewDocumentRepository(db)
	testRepo := postgres.NewTestRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	answerRepo := postgres.NewAnswerRepository(db)

	// Run database seeders
	seeder := persistence.NewSeeder(userRepo, roleRepo, cfg, appLogger)
	if err := seeder.Seed(context.Background()); err != nil {
		appLogger.Error("Failed to run database seeders", zap.Error(err))
	}

	// Initialize JWT manager
	jwtManager, err := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		appLogger.Fatal("Failed to initialize JWT manager", zap.Error(err))
	}

	// Initialize document parser factory (Factory Pattern)
	parserFactory := parser.NewDocumentParserFactory()

	// Initialize LLM factory (Factory Pattern + Strategy Pattern)
	llmFactory := llm.NewLLMFactory(
		cfg.LLM.PerplexityAPIKey,
		cfg.LLM.OpenAIAPIKey,
		cfg.LLM.YandexAPIKey,
		cfg.LLM.YandexFolderID,
		cfg.LLM.YandexModel,
	)

	// Initialize Moodle components
	xmlExporter := moodle.NewMoodleXMLExporter()
	var moodleClient *moodle.Client
	if cfg.Moodle.URL != "" && cfg.Moodle.Token != "" {
		moodleClient = moodle.NewClient(cfg.Moodle.URL, cfg.Moodle.Token)
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(
		userRepo,
		roleRepo,
		jwtManager,
		cfg.Cookie.Name,
		cfg.Cookie.Domain,
		cfg.Cookie.Path,
		cfg.Cookie.SameSite,
		cfg.JWT.Expiration,
		cfg.Cookie.Secure,
		cfg.Cookie.HTTPOnly,
	)
	userHandler := handler.NewUserHandler(userRepo, roleRepo)
	documentHandler := handler.NewDocumentHandler(
		documentRepo,
		userRepo,
		parserFactory,
		cfg.File.UploadDir,
		cfg.File.MaxFileSize,
	)
	testHandler := handler.NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, llmFactory, xmlExporter)
	moodleHandler := handler.NewMoodleHandler(
		testRepo,
		questionRepo,
		answerRepo,
		xmlExporter,
		moodleClient,
	)
	statsHandler := handler.NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.RequestIDMiddleware())
	app.Use(logger.HTTPMiddleware(appLogger))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
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
	router.SetupRoutes(app, authHandler, userHandler, documentHandler, testHandler, moodleHandler, statsHandler, jwtManager, cfg.Cookie.Name)

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
				"moodle":    "/api/v1/moodle",
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
	if err := app.Shutdown(); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
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
