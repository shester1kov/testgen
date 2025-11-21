package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/pkg/utils"
)

// AuthMiddleware creates authentication middleware that supports both cookie and Authorization header
func AuthMiddleware(jwtManager *utils.JWTManager, cookieName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		// Try to get token from cookie first
		token = c.Cookies(cookieName)

		// If no cookie, try Authorization header
		if token == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				// Extract token from "Bearer <token>"
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		// If still no token found
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authentication token",
			})
		}

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware creates role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole").(string)

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "insufficient permissions",
		})
	}
}
