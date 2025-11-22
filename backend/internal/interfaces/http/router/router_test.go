package router

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes_RegistersExpectedEndpoints(t *testing.T) {
	jwtManager, err := utils.NewJWTManager("test-secret-key-must-be-at-least-32-chars-long", "1h")
	assert.NoError(t, err)

	app := fiber.New()

	SetupRoutes(
		app,
		&handler.AuthHandler{},
		&handler.UserHandler{},
		&handler.DocumentHandler{},
		&handler.TestHandler{},
		&handler.MoodleHandler{},
		jwtManager,
		"token",
	)

	routes := app.GetRoutes()

	expected := map[string]bool{
		"POST /api/v1/auth/register":          true,
		"POST /api/v1/auth/login":             true,
		"POST /api/v1/auth/logout":            true,
		"GET /api/v1/auth/me":                 true,
		"GET /api/v1/users/":                  true,
		"PUT /api/v1/users/:id/role":          true,
		"POST /api/v1/documents/":             true,
		"GET /api/v1/documents/":              true,
		"GET /api/v1/documents/:id":           true,
		"DELETE /api/v1/documents/:id":        true,
		"POST /api/v1/documents/:id/parse":    true,
		"POST /api/v1/tests/":                 true,
		"GET /api/v1/tests/":                  true,
		"GET /api/v1/tests/:id":               true,
		"DELETE /api/v1/tests/:id":            true,
		"POST /api/v1/tests/generate":         true,
		"GET /api/v1/moodle/connection":       true,
		"GET /api/v1/moodle/courses":          true,
		"GET /api/v1/moodle/tests/:id/export": true,
		"POST /api/v1/moodle/tests/:id/sync":  true,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		delete(expected, key)
	}

	assert.Empty(t, expected, "all expected routes should be registered")
}
