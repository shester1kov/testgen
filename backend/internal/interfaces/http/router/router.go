package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/middleware"
	"github.com/shester1kov/testgen-backend/pkg/utils"
)

// action constants for activity logging
const (
	actionAuthRegister    = "auth.register"
	actionAuthLogin       = "auth.login"
	actionAuthLogout      = "auth.logout"
	actionUserUpdateRole  = "user.update_role"
	actionDocumentUpload  = "document.upload"
	actionDocumentDelete  = "document.delete"
	actionDocumentParse   = "document.parse"
	actionTestCreate      = "test.create"
	actionTestUpdate      = "test.update"
	actionTestDelete      = "test.delete"
	actionTestGenerate    = "test.generate"
	actionTestExport      = "test.export"
	actionQuestionUpdate  = "question.update"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	documentHandler *handler.DocumentHandler,
	testHandler *handler.TestHandler,
	statsHandler *handler.StatsHandler,
	jwtManager *utils.JWTManager,
	cookieName string,
	activityLogger *middleware.ActivityLogger,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", activityLogger.Log(actionAuthRegister, "user"), authHandler.Register)
	auth.Post("/login", activityLogger.Log(actionAuthLogin, "user"), authHandler.Login)
	auth.Post("/logout", activityLogger.Log(actionAuthLogout, "user"), authHandler.Logout)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(jwtManager, cookieName), authHandler.GetMe)

	// User management routes (protected)
	users := api.Group("/users", middleware.AuthMiddleware(jwtManager, cookieName))
	users.Get("/", middleware.RequireTeacherOrAdmin(), userHandler.ListUsers)                                                        // Teachers can view users
	users.Put("/:id/role", middleware.RequireAdmin(), activityLogger.Log(actionUserUpdateRole, "user"), userHandler.UpdateUserRole) // Only admin can change roles

	// Document routes (protected - teacher and admin only for upload)
	documents := api.Group("/documents", middleware.AuthMiddleware(jwtManager, cookieName))
	documents.Post("/", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionDocumentUpload, "document"), documentHandler.Upload)         // Only teachers/admin can upload
	documents.Get("/", documentHandler.List)                                                                                                       // All can list
	documents.Get("/:id", documentHandler.GetByID)                                                                                                 // All can view
	documents.Delete("/:id", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionDocumentDelete, "document"), documentHandler.Delete)    // Only teachers/admin can delete
	documents.Post("/:id/parse", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionDocumentParse, "document"), documentHandler.Parse) // Only teachers/admin can parse

	// Test routes (protected - teacher and admin only for creation/editing)
	tests := api.Group("/tests", middleware.AuthMiddleware(jwtManager, cookieName))
	tests.Post("/", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionTestCreate, "test"), testHandler.Create)                                                         // Only teachers/admin can create
	tests.Get("/", testHandler.List)                                                                                                                                               // All can list (students see assigned tests)
	tests.Get("/:id", testHandler.GetByID)                                                                                                                                         // All can view
	tests.Put("/:id", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionTestUpdate, "test"), testHandler.Update)                                                       // Only teachers/admin can update
	tests.Delete("/:id", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionTestDelete, "test"), testHandler.Delete)                                                    // Only teachers/admin can delete
	tests.Post("/generate", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionTestGenerate, "test"), testHandler.Generate)                                             // Only teachers/admin can generate
	tests.Put("/:testId/questions/:questionId", middleware.RequireTeacherOrAdmin(), activityLogger.Log(actionQuestionUpdate, "question"), testHandler.UpdateQuestion)             // Only teachers/admin can update questions
	tests.Get("/:id/export/json", testHandler.ExportToJSON)                                                                                                                        // Export test to JSON
	tests.Get("/:id/export/xml", testHandler.ExportToXML)                                                                                                                          // Export test to Moodle XML (legacy)
	tests.Get("/:id/export", activityLogger.Log(actionTestExport, "test"), testHandler.Export)                                                                                    // Export test to various formats
	tests.Get("/export/formats", testHandler.GetExportFormats)                                                                                                                     // Get available export formats

	// Stats routes (protected - all authenticated users)
	stats := api.Group("/stats", middleware.AuthMiddleware(jwtManager, cookieName))
	stats.Get("/dashboard", statsHandler.GetDashboardStats)
}
