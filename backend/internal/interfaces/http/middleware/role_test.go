package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

// POSITIVE TEST: RequireRole allows user with matching role
func TestRequireRole_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameTeacher))
		return c.Next()
	})

	app.Get("/test", RequireRole(entity.RoleNameTeacher), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// POSITIVE TEST: RequireRole allows user with one of multiple allowed roles
func TestRequireRole_MultipleRoles_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameTeacher))
		return c.Next()
	})

	app.Get("/test", RequireRole(entity.RoleNameTeacher, entity.RoleNameAdmin), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// NEGATIVE TEST: RequireRole denies user without matching role
func TestRequireRole_Forbidden(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameStudent))
		return c.Next()
	})

	app.Get("/test", RequireRole(entity.RoleNameTeacher), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// NEGATIVE TEST: RequireRole denies when userRole is missing
func TestRequireRole_MissingRole(t *testing.T) {
	app := fiber.New()

	app.Get("/test", RequireRole(entity.RoleNameTeacher), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

// NEGATIVE TEST: RequireRole handles invalid role format
func TestRequireRole_InvalidRoleFormat(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", 12345) // Invalid type (int instead of string)
		return c.Next()
	})

	app.Get("/test", RequireRole(entity.RoleNameTeacher), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

// POSITIVE TEST: RequireAdmin allows admin
func TestRequireAdmin_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameAdmin))
		return c.Next()
	})

	app.Get("/test", RequireAdmin(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// NEGATIVE TEST: RequireAdmin denies non-admin
func TestRequireAdmin_Forbidden(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameTeacher))
		return c.Next()
	})

	app.Get("/test", RequireAdmin(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// POSITIVE TEST: RequireTeacherOrAdmin allows teacher
func TestRequireTeacherOrAdmin_Teacher_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameTeacher))
		return c.Next()
	})

	app.Get("/test", RequireTeacherOrAdmin(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// POSITIVE TEST: RequireTeacherOrAdmin allows admin
func TestRequireTeacherOrAdmin_Admin_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameAdmin))
		return c.Next()
	})

	app.Get("/test", RequireTeacherOrAdmin(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// NEGATIVE TEST: RequireTeacherOrAdmin denies student
func TestRequireTeacherOrAdmin_Student_Forbidden(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameStudent))
		return c.Next()
	})

	app.Get("/test", RequireTeacherOrAdmin(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// POSITIVE TEST: RequireStudent allows student
func TestRequireStudent_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameStudent))
		return c.Next()
	})

	app.Get("/test", RequireStudent(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// NEGATIVE TEST: RequireStudent denies teacher
func TestRequireStudent_Teacher_Forbidden(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userRole", string(entity.RoleNameTeacher))
		return c.Next()
	})

	app.Get("/test", RequireStudent(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// POSITIVE TEST: RequireAnyAuthenticated allows any authenticated user
func TestRequireAnyAuthenticated_Success(t *testing.T) {
	testCases := []struct {
		name string
		role entity.RoleName
	}{
		{"admin", entity.RoleNameAdmin},
		{"teacher", entity.RoleNameTeacher},
		{"student", entity.RoleNameStudent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userRole", string(tc.role))
				return c.Next()
			})

			app.Get("/test", RequireAnyAuthenticated(), func(c *fiber.Ctx) error {
				return c.SendString("success")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

// NEGATIVE TEST: RequireAnyAuthenticated denies unauthenticated user
func TestRequireAnyAuthenticated_Unauthenticated(t *testing.T) {
	app := fiber.New()

	app.Get("/test", RequireAnyAuthenticated(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
