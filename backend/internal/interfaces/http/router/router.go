package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/middleware"
	"github.com/shester1kov/testgen-backend/pkg/utils"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	documentHandler *handler.DocumentHandler,
	testHandler *handler.TestHandler,
	moodleHandler *handler.MoodleHandler,
	jwtManager *utils.JWTManager,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(jwtManager), authHandler.GetMe)

	// Document routes (protected)
	documents := api.Group("/documents", middleware.AuthMiddleware(jwtManager))
	documents.Post("/", documentHandler.Upload)
	documents.Get("/", documentHandler.List)
	documents.Get("/:id", documentHandler.GetByID)
	documents.Delete("/:id", documentHandler.Delete)
	documents.Post("/:id/parse", documentHandler.Parse)

	// Test routes (protected)
	tests := api.Group("/tests", middleware.AuthMiddleware(jwtManager))
	tests.Post("/", testHandler.Create)
	tests.Get("/", testHandler.List)
	tests.Get("/:id", testHandler.GetByID)
	tests.Delete("/:id", testHandler.Delete)
	tests.Post("/generate", testHandler.Generate)

	// Moodle integration routes (protected)
	moodle := api.Group("/moodle", middleware.AuthMiddleware(jwtManager))
	moodle.Get("/connection", moodleHandler.ValidateMoodleConnection)
	moodle.Get("/courses", moodleHandler.GetMoodleCourses)
	moodle.Get("/tests/:id/export", moodleHandler.ExportToXML)
	moodle.Post("/tests/:id/sync", moodleHandler.SyncToMoodle)
}
