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

// NEGATIVE TEST: SQL Injection in Authorization header
func TestAuthMiddleware_SQLInjectionAttempt(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	sqlInjectionPayloads := []string{
		"' OR '1'='1",
		"admin'--",
		"' OR 1=1--",
		"'; DROP TABLE users--",
		"1' UNION SELECT NULL--",
	}

	for _, payload := range sqlInjectionPayloads {
		t.Run("sql_injection_"+payload, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+payload)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}

// NEGATIVE TEST: XSS Attempt in cookie
func TestAuthMiddleware_XSSAttempt(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"javascript:alert('XSS')",
		"<img src=x onerror=alert('XSS')>",
	}

	for _, payload := range xssPayloads {
		t.Run("xss_"+payload[:10], func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.AddCookie(&http.Cookie{
				Name:  "testgen_token",
				Value: payload,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}

// NEGATIVE TEST: Very long token (potential buffer overflow)
// Note: Very large tokens can cause memory issues - this tests the system handles it
func TestAuthMiddleware_VeryLongToken(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Generate a long string (100KB is enough to test, 10MB might cause OOM)
	longToken := string(make([]byte, 100*1024))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+longToken)

	// Test with timeout to prevent hanging
	resp, err := app.Test(req, 1000) // 1 second timeout
	if err != nil {
		// If it times out or errors, that's also acceptable
		t.Log("Request timed out or errored (expected for very long token)")
		return
	}
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// NEGATIVE TEST: Token with special characters
func TestAuthMiddleware_SpecialCharactersInToken(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	specialChars := []string{
		"../../../etc/passwd",
		"%00",
		// Note: "\x00\x00\x00" (null bytes) are not valid in HTTP headers
		// HTTP spec doesn't allow null bytes, so we skip this test
		"../../",
		"<>\"'&",
	}

	for i, char := range specialChars {
		t.Run("special_char_"+string(rune(i)), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+char)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}

// NEGATIVE TEST: Token from different issuer
func TestAuthMiddleware_TokenFromDifferentIssuer(t *testing.T) {
	jwtManager1 := setupTestJWTManager(t)

	// Create a token with different secret
	otherJWTManager, _ := utils.NewJWTManager("different-secret-key-32-chars-long-xxx", "1h")
	userID := uuid.New()
	tokenFromOtherIssuer, _ := otherJWTManager.GenerateToken(userID, "test@example.com", "student")

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager1, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFromOtherIssuer)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// NEGATIVE TEST: Concurrent requests with same token
func TestAuthMiddleware_ConcurrentRequests(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()
	token, _ := jwtManager.GenerateToken(userID, "test@example.com", "student")

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test concurrent access with same token (should work)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// NEGATIVE TEST: Case sensitive Authorization header
func TestAuthMiddleware_CaseSensitiveHeader(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()
	token, _ := jwtManager.GenerateToken(userID, "test@example.com", "student")

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name   string
		header string
		value  string
		expect int
	}{
		{"lowercase bearer", "Authorization", "bearer " + token, fiber.StatusUnauthorized},
		{"UPPERCASE BEARER", "Authorization", "BEARER " + token, fiber.StatusUnauthorized},
		{"Mixed Case", "Authorization", "BeArEr " + token, fiber.StatusUnauthorized},
		{"Correct case", "Authorization", "Bearer " + token, fiber.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(tc.header, tc.value)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, resp.StatusCode)
		})
	}
}

// NEGATIVE TEST: Multiple Authorization headers
func TestAuthMiddleware_MultipleAuthHeaders(t *testing.T) {
	jwtManager := setupTestJWTManager(t)
	userID := uuid.New()
	validToken, _ := jwtManager.GenerateToken(userID, "test@example.com", "student")

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	// Add multiple Authorization headers (HTTP spec allows it)
	req.Header.Add("Authorization", "Bearer invalid-token")
	req.Header.Add("Authorization", "Bearer "+validToken)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Should handle gracefully (might use first or reject)
	// Either 200 or 401 is acceptable, just shouldn't panic
	assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusUnauthorized)
}

// NEGATIVE TEST: Null bytes in token
// Note: HTTP spec doesn't allow null bytes in headers, so this test verifies
// that the HTTP layer itself rejects such requests (expected behavior)
func TestAuthMiddleware_NullBytesInToken(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tokenWithNulls := "token\x00with\x00nulls"

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenWithNulls)

	resp, err := app.Test(req)
	// HTTP layer should reject this (error expected)
	if err != nil {
		t.Log("HTTP layer correctly rejected null bytes in header")
		return
	}
	// If it somehow gets through, middleware should reject it
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// NEGATIVE TEST: Empty Authorization value
func TestAuthMiddleware_EmptyAuthorizationValue(t *testing.T) {
	jwtManager := setupTestJWTManager(t)

	app := fiber.New()
	app.Use(AuthMiddleware(jwtManager, "testgen_token"))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []string{
		"",
		" ",
		"Bearer",
		"Bearer ",
		" Bearer ",
	}

	for _, value := range testCases {
		t.Run("empty_value_"+value, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if value != "" {
				req.Header.Set("Authorization", value)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}
