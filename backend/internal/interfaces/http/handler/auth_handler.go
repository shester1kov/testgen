package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/auth"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	registerUseCase *auth.RegisterUseCase
	loginUseCase    *auth.LoginUseCase
	getMeUseCase    *auth.GetMeUseCase
	cookieName      string
	cookieDomain    string
	cookiePath      string
	cookieSecure    bool
	cookieHTTPOnly  bool
	cookieSameSite  string
	jwtExpiration   string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	registerUseCase *auth.RegisterUseCase,
	loginUseCase *auth.LoginUseCase,
	getMeUseCase *auth.GetMeUseCase,
	cookieName, cookieDomain, cookiePath, cookieSameSite, jwtExpiration string,
	cookieSecure, cookieHTTPOnly bool,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		getMeUseCase:    getMeUseCase,
		cookieName:      cookieName,
		cookieDomain:    cookieDomain,
		cookiePath:      cookiePath,
		cookieSecure:    cookieSecure,
		cookieHTTPOnly:  cookieHTTPOnly,
		cookieSameSite:  cookieSameSite,
		jwtExpiration:   jwtExpiration,
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

	// Normalize email to lowercase for case-insensitive matching
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)

	result, err := h.registerUseCase.Execute(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(
				dto.NewErrorResponse(dto.ErrCodeUserExists, "User with this email already exists"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "Failed to register user"),
		)
	}

	// Set token in HTTP-only cookie
	h.setCookie(c, result.Token)

	return c.Status(fiber.StatusCreated).JSON(result)
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

	// Normalize email to lowercase for case-insensitive matching
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	result, err := h.loginUseCase.Execute(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid email or password") {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidCredentials, "Invalid email or password"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "Failed to login"),
		)
	}

	// Set token in HTTP-only cookie
	h.setCookie(c, result.Token)

	return c.JSON(result)
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
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	result, err := h.getMeUseCase.Execute(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeUserNotFound, "User not found"),
		)
	}

	return c.JSON(result)
}
