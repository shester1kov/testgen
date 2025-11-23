package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// RequireRole creates a middleware that checks if user has one of the required roles
func RequireRole(allowedRoles ...entity.RoleName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by JWT middleware)
		userRoleRaw := c.Locals("userRole")
		if userRoleRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
			)
		}

		// The role is stored as string in JWT claims
		roleStr, ok := userRoleRaw.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeInternalError, "Invalid role format"),
			)
		}

		userRole := entity.RoleName(roleStr)

		// Check if user's role is in the allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "Access denied: insufficient permissions"),
		)
	}
}

// RequireAdmin middleware that only allows admin role
func RequireAdmin() fiber.Handler {
	return RequireRole(entity.RoleNameAdmin)
}

// RequireTeacherOrAdmin middleware that allows teacher or admin roles
func RequireTeacherOrAdmin() fiber.Handler {
	return RequireRole(entity.RoleNameTeacher, entity.RoleNameAdmin)
}

// RequireStudent middleware that only allows student role
func RequireStudent() fiber.Handler {
	return RequireRole(entity.RoleNameStudent)
}

// RequireAnyAuthenticated middleware that allows any authenticated user
func RequireAnyAuthenticated() fiber.Handler {
	return RequireRole(entity.RoleNameAdmin, entity.RoleNameTeacher, entity.RoleNameStudent)
}
