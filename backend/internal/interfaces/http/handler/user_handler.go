package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/user"
)

// UserHandler handles user management requests
type UserHandler struct {
	listUseCase       *user.ListUseCase
	updateRoleUseCase *user.UpdateRoleUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(listUseCase *user.ListUseCase, updateRoleUseCase *user.UpdateRoleUseCase) *UserHandler {
	return &UserHandler{
		listUseCase:       listUseCase,
		updateRoleUseCase: updateRoleUseCase,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Get paginated list of all users (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.UserListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "limit must be between 1 and 100"),
		)
	}
	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "offset must be non-negative"),
		)
	}

	result, err := h.listUseCase.Execute(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch users"),
		)
	}

	return c.JSON(result)
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update a user's role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRoleRequest true "Role update request"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	// Parse user ID
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid user ID"),
		)
	}

	// Parse request body
	var req dto.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request body"),
		)
	}

	result, err := h.updateRoleUseCase.Execute(c.Context(), userID, req.RoleName)
	if err != nil {
		if strings.Contains(err.Error(), "invalid role name") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidRole, "invalid role name"),
			)
		}
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeUserNotFound, "user not found"),
			)
		}
		if strings.Contains(err.Error(), "role not found") {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeRoleNotFound, "role not found in database"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to update user role"),
		)
	}

	return c.JSON(result)
}
