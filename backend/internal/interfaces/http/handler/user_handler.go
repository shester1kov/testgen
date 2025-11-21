package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// UserHandler handles user management requests
type UserHandler struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		roleRepo: roleRepo,
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

	// Get users
	users, err := h.userRepo.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch users"),
		)
	}

	// Get total count
	total, err := h.userRepo.Count(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to count users"),
		)
	}

	// Convert to DTOs
	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.GetRoleName(),
		}
	}

	return c.JSON(dto.UserListResponse{
		Users:  userDTOs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
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

	// Validate role name
	roleName := entity.RoleName(req.RoleName)
	if roleName != entity.RoleNameAdmin && roleName != entity.RoleNameTeacher && roleName != entity.RoleNameStudent {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidRole, "invalid role name"),
		)
	}

	// Find user
	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeUserNotFound, "user not found"),
		)
	}

	// Find role by name
	role, err := h.roleRepo.FindByName(c.Context(), roleName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeRoleNotFound, "role not found in database"),
		)
	}

	// Update user role
	user.RoleID = role.ID
	if err := h.userRepo.Update(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to update user role"),
		)
	}

	// Load updated role for response
	user.Role = role

	return c.JSON(dto.UserDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.GetRoleName(),
	})
}
