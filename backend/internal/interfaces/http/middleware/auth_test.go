package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func setupTestJWTManager(t *testing.T) *utils.JWTManager {
	jwtManager, err := utils.NewJWTManager("test-secret-key-must-be-at-least-32-chars-long", "1h")
	assert.NoError(t, err)
	return jwtManager
}

func TestAuthMiddleware_WithValidCookie(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()

	// Generate a valid token
	token, err := jwtManager.GenerateToken(userID, "test@example.com", "student")
	assert.NoError(t, err)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userID":    c.Locals("userID"),
			"userEmail": c.Locals("userEmail"),
			"userRole":  c.Locals("userRole"),
		})
	})

	// Create request with cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "testgen_token",
		Value: token,
	})

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_WithValidAuthorizationHeader(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()

	// Generate a valid token
	token, err := jwtManager.GenerateToken(userID, "test@example.com", "student")
	assert.NoError(t, err)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userID":    c.Locals("userID"),
			"userEmail": c.Locals("userEmail"),
			"userRole":  c.Locals("userRole"),
		})
	})

	// Create request with Authorization header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_CookiePriorityOverHeader(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID1 := uuid.New()
	userID2 := uuid.New()

	// Generate two different valid tokens
	tokenCookie, err := jwtManager.GenerateToken(userID1, "cookie@example.com", "student")
	assert.NoError(t, err)

	tokenHeader, err := jwtManager.GenerateToken(userID2, "header@example.com", "teacher")
	assert.NoError(t, err)

	var capturedEmail interface{}

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		capturedEmail = c.Locals("userEmail")
		return c.JSON(fiber.Map{
			"userEmail": c.Locals("userEmail"),
		})
	})

	// Create request with both cookie and header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "testgen_token",
		Value: tokenCookie,
	})
	req.Header.Set("Authorization", "Bearer "+tokenHeader)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify that cookie took priority (email from cookie token)
	assert.Equal(t, "cookie@example.com", capturedEmail)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Create request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "testgen_token",
		Value: "invalid-token",
	})

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Create JWT manager with very short expiration
	jwtManager, err := utils.NewJWTManager("test-secret-key-must-be-at-least-32-chars-long", "1ns")
	assert.NoError(t, err)

	userID := uuid.New()
	token, err := jwtManager.GenerateToken(userID, "test@example.com", "student")
	assert.NoError(t, err)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Create request with expired token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "testgen_token",
		Value: token,
	})

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test cases for invalid header formats
	testCases := []string{
		"InvalidFormat",
		"Bearer",
		"Token abc123",
		"Bearer token with spaces",
	}

	for _, headerValue := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", headerValue)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Header: "+headerValue)
	}
}

func TestAuthMiddleware_StoresUserInfoInContext(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()
	userEmail := "test@example.com"
	userRole := "student"

	// Generate a valid token
	token, err := jwtManager.GenerateToken(userID, userEmail, userRole)
	assert.NoError(t, err)

	// Setup Fiber app with middleware
	var capturedUserID interface{}
	var capturedUserEmail interface{}
	var capturedUserRole interface{}

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		capturedUserID = c.Locals("userID")
		capturedUserEmail = c.Locals("userEmail")
		capturedUserRole = c.Locals("userRole")
		return c.SendStatus(fiber.StatusOK)
	})

	// Create request with cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "testgen_token",
		Value: token,
	})

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify context values
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, userEmail, capturedUserEmail)
	assert.Equal(t, userRole, capturedUserRole)
}

func TestRoleMiddleware_AllowedRole(t *testing.T) {
	app := fiber.New()

	// Simulate auth middleware setting user role
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", "admin")
		return c.Next()
	})

	app.Use(RoleMiddleware("admin", "teacher"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRoleMiddleware_ForbiddenRole(t *testing.T) {
	app := fiber.New()

	// Simulate auth middleware setting user role
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", "student")
		return c.Next()
	})

	app.Use(RoleMiddleware("admin", "teacher"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}
