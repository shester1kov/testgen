package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	apperrors "github.com/shester1kov/testgen-backend/pkg/errors"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	jwtManager     *utils.JWTManager
	cookieName     string
	cookieDomain   string
	cookiePath     string
	cookieSecure   bool
	cookieHTTPOnly bool
	cookieSameSite string
	jwtExpiration  string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtManager *utils.JWTManager, cookieName, cookieDomain, cookiePath, cookieSameSite, jwtExpiration string, cookieSecure, cookieHTTPOnly bool) *AuthHandler {
	return &AuthHandler{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		jwtManager:     jwtManager,
		cookieName:     cookieName,
		cookieDomain:   cookieDomain,
		cookiePath:     cookiePath,
		cookieSecure:   cookieSecure,
		cookieHTTPOnly: cookieHTTPOnly,
		cookieSameSite: cookieSameSite,
		jwtExpiration:  jwtExpiration,
	}
}

// setCookie sets the JWT token in an HTTP-only cookie
func (h *AuthHandler) setCookie(c *fiber.Ctx, token string) {
	// Parse JWT expiration duration
	expiration, err := time.ParseDuration(h.jwtExpiration)
	if err != nil {
		expiration = 24 * time.Hour // fallback to 24 hours
	}

	cookie := &fiber.Cookie{
		Name:     h.cookieName,
		Value:    token,
		Path:     h.cookiePath,
		Domain:   h.cookieDomain,
		MaxAge:   int(expiration.Seconds()),
		Secure:   h.cookieSecure,
		HTTPOnly: h.cookieHTTPOnly,
		SameSite: h.cookieSameSite,
	}

	c.Cookie(cookie)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with default student role
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "Invalid request body"),
		)
	}

	// Check if user already exists
	existingUser, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(
			dto.NewErrorResponse(dto.ErrCodeUserExists, "User with this email already exists"),
		)
	}

	// Get default student role
	studentRole, err := h.roleRepo.FindByName(c.Context(), entity.RoleNameStudent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeRoleNotFound, "Failed to get default student role"),
		)
	}

	// Create new user with student role
	user := &entity.User{
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   studentRole.ID,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "Failed to hash password"),
		)
	}

	if err := h.userRepo.Create(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "Failed to create user"),
		)
	}

	// Load role for response
	user.Role = studentRole

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, user.GetRoleName())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "Failed to generate token"),
		)
	}

	// Set token in HTTP-only cookie
	h.setCookie(c, token)

	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.GetRoleName(),
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "Invalid request body"),
		)
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidCredentials, "Invalid email or password"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "Failed to find user"),
		)
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidCredentials, "Invalid email or password"),
		)
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, user.GetRoleName())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "Failed to generate token"),
		)
	}

	// Set token in HTTP-only cookie
	h.setCookie(c, token)

	return c.JSON(dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.GetRoleName(),
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Clear authentication cookie
// @Tags auth
// @Success 200 {object} dto.MessageResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear cookie by setting it with expired time
	c.Cookie(&fiber.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     h.cookiePath,
		Domain:   h.cookieDomain,
		MaxAge:   -1,
		Secure:   h.cookieSecure,
		HTTPOnly: h.cookieHTTPOnly,
		SameSite: h.cookieSameSite,
	})

	return c.JSON(dto.NewMessageResponse("Logged out successfully"))
}

// GetMe godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserDTO
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	rawUserID := c.Locals("userID")
	if rawUserID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	var userID uuid.UUID
	switch v := rawUserID.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsedID, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
			)
		}
		userID = parsedID
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		appErr := apperrors.NotFound("user not found")
		return c.Status(appErr.Code).JSON(
			dto.NewErrorResponse(dto.ErrCodeUserNotFound, appErr.Message),
		)
	}

	return c.JSON(dto.UserDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.GetRoleName(),
	})
}
