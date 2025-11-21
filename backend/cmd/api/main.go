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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "github.com/shester1kov/testgen-backend/docs"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/router"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	documentRepo := postgres.NewDocumentRepository(db)
	testRepo := postgres.NewTestRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	answerRepo := postgres.NewAnswerRepository(db)

	// Initialize JWT manager
	jwtManager, err := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		log.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	// Initialize document parser factory (Factory Pattern)
	parserFactory := parser.NewDocumentParserFactory()

	// Initialize LLM factory (Factory Pattern + Strategy Pattern)
	llmFactory := llm.NewLLMFactory(
		cfg.LLM.PerplexityAPIKey,
		cfg.LLM.OpenAIAPIKey,
		cfg.LLM.YandexAPIKey,
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
		parserFactory,
		cfg.File.UploadDir,
		cfg.File.MaxFileSize,
	)
	testHandler := handler.NewTestHandler(testRepo, documentRepo, llmFactory)
	moodleHandler := handler.NewMoodleHandler(
		testRepo,
		questionRepo,
		answerRepo,
		xmlExporter,
		moodleClient,
	)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Test Generation System API is running",
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup routes
	router.SetupRoutes(app, authHandler, userHandler, documentHandler, testHandler, moodleHandler, jwtManager, cfg.Cookie.Name)

	// Root endpoint
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Test Generation System API v1",
			"endpoints": fiber.Map{
				"auth":      "/api/v1/auth",
				"documents": "/api/v1/documents",
				"tests":     "/api/v1/tests",
				"moodle":    "/api/v1/moodle",
			},
		})
	})

	// Start server
	port := cfg.Server.Port

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
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
