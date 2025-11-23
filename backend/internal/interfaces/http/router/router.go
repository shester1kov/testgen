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
	userHandler *handler.UserHandler,
	documentHandler *handler.DocumentHandler,
	testHandler *handler.TestHandler,
	moodleHandler *handler.MoodleHandler,
	jwtManager *utils.JWTManager,
	cookieName string,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(jwtManager, cookieName), authHandler.GetMe)

	// User management routes (protected)
	users := api.Group("/users", middleware.AuthMiddleware(jwtManager, cookieName))
	users.Get("/", middleware.RequireTeacherOrAdmin(), userHandler.ListUsers)           // Teachers can view users
	users.Put("/:id/role", middleware.RequireAdmin(), userHandler.UpdateUserRole)       // Only admin can change roles

	// Document routes (protected - teacher and admin only for upload)
	documents := api.Group("/documents", middleware.AuthMiddleware(jwtManager, cookieName))
	documents.Post("/", middleware.RequireTeacherOrAdmin(), documentHandler.Upload)     // Only teachers/admin can upload
	documents.Get("/", documentHandler.List)                                             // All can list
	documents.Get("/:id", documentHandler.GetByID)                                       // All can view
	documents.Delete("/:id", middleware.RequireTeacherOrAdmin(), documentHandler.Delete) // Only teachers/admin can delete
	documents.Post("/:id/parse", middleware.RequireTeacherOrAdmin(), documentHandler.Parse) // Only teachers/admin can parse

	// Test routes (protected - teacher and admin only for creation/editing)
	tests := api.Group("/tests", middleware.AuthMiddleware(jwtManager, cookieName))
	tests.Post("/", middleware.RequireTeacherOrAdmin(), testHandler.Create)             // Only teachers/admin can create
	tests.Get("/", testHandler.List)                                                     // All can list (students see assigned tests)
	tests.Get("/:id", testHandler.GetByID)                                               // All can view
	tests.Delete("/:id", middleware.RequireTeacherOrAdmin(), testHandler.Delete)         // Only teachers/admin can delete
	tests.Post("/generate", middleware.RequireTeacherOrAdmin(), testHandler.Generate)    // Only teachers/admin can generate

	// Moodle integration routes (protected - teacher and admin only)
	moodle := api.Group("/moodle", middleware.AuthMiddleware(jwtManager, cookieName), middleware.RequireTeacherOrAdmin())
	moodle.Get("/connection", moodleHandler.ValidateMoodleConnection)
	moodle.Get("/courses", moodleHandler.GetMoodleCourses)
	moodle.Get("/tests/:id/export", moodleHandler.ExportToXML)
	moodle.Post("/tests/:id/sync", moodleHandler.SyncToMoodle)
}
