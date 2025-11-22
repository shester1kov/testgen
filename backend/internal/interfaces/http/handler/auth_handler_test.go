package handler

import (
	"bytes"
	"context"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
"github.com/shester1kov/testgen-backend/internal/domain/entity"
"github.com/shester1kov/testgen-backend/pkg/utils"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"github.com/stretchr/testify/require"
"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockRoleRepository is a mock implementation of RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func setupAuthHandler(t *testing.T) (*AuthHandler, *MockUserRepository, *MockRoleRepository, *utils.JWTManager) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	jwtManager, err := utils.NewJWTManager("test-secret-key-must-be-at-least-32-chars-long", "1h")
	assert.NoError(t, err)

	handler := NewAuthHandler(
		mockUserRepo,
		mockRoleRepo,
		jwtManager,
		"testgen_token",
		"",
		"/",
		"Lax",
		"1h",
		false,
		true,
	)

	return handler, mockUserRepo, mockRoleRepo, jwtManager
}

func TestLogin_Success_SetsCookie(t *testing.T) {
	handler, mockUserRepo, _, _ := setupAuthHandler(t)

	// Create test user with role
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}
	user := &entity.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}
	err := user.SetPassword("password123")
	assert.NoError(t, err)

	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/login", handler.Login)

	// Create request
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check cookie is set
	cookies := resp.Cookies()
	assert.NotEmpty(t, cookies, "Cookie should be set")

	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "testgen_token" {
			foundCookie = cookie
			break
		}
	}

	assert.NotNil(t, foundCookie, "testgen_token cookie should be present")
	assert.NotEmpty(t, foundCookie.Value, "Cookie value should not be empty")
	assert.True(t, foundCookie.HttpOnly, "Cookie should be HTTP-only")
	assert.Equal(t, "/", foundCookie.Path)
	assert.Greater(t, foundCookie.MaxAge, 0, "Cookie should have positive MaxAge")

	mockUserRepo.AssertExpectations(t)
}

func TestRegister_Success_SetsCookie(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo, _ := setupAuthHandler(t)

	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "newuser@example.com").Return(nil, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameStudent).Return(studentRole, nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/register", handler.Register)

	// Create request (note: Role field removed from RegisterRequest)
	registerReq := dto.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Check cookie is set
	cookies := resp.Cookies()
	assert.NotEmpty(t, cookies, "Cookie should be set")

	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "testgen_token" {
			foundCookie = cookie
			break
		}
	}

	assert.NotNil(t, foundCookie, "testgen_token cookie should be present")
	assert.NotEmpty(t, foundCookie.Value, "Cookie value should not be empty")
	assert.True(t, foundCookie.HttpOnly, "Cookie should be HTTP-only")

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestLogout_ClearsCookie(t *testing.T) {
	handler, _, _, _ := setupAuthHandler(t)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/logout", handler.Logout)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check cookie is cleared (MaxAge = -1)
	cookies := resp.Cookies()
	assert.NotEmpty(t, cookies, "Cookie should be set to clear")

	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "testgen_token" {
			foundCookie = cookie
			break
		}
	}

	assert.NotNil(t, foundCookie, "testgen_token cookie should be present")
	assert.Equal(t, "", foundCookie.Value, "Cookie value should be empty")
	assert.Equal(t, -1, foundCookie.MaxAge, "Cookie MaxAge should be -1 to clear it")
}

func TestLogin_ReturnsTokenInResponse(t *testing.T) {
	handler, mockUserRepo, _, _ := setupAuthHandler(t)

	// Create test user with role
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}
	user := &entity.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}
	err := user.SetPassword("password123")
	assert.NoError(t, err)

	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/login", handler.Login)

	// Create request
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	var authResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, authResp.Token, "Token should be present in response body")
	assert.Equal(t, user.Email, authResp.User.Email)
	assert.Equal(t, user.FullName, authResp.User.FullName)

	mockUserRepo.AssertExpectations(t)
}

func TestSetCookie_ParsesExpirationCorrectly(t *testing.T) {
	handler, _, _, _ := setupAuthHandler(t)

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		handler.setCookie(c, "test-token")
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "testgen_token" {
			foundCookie = cookie
			break
		}
	}

	assert.NotNil(t, foundCookie)
	// 1 hour = 3600 seconds
	assert.Equal(t, 3600, foundCookie.MaxAge)
}

func TestSetCookie_FallbackOnInvalidExpiration(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	jwtManager, _ := utils.NewJWTManager("test-secret-key-must-be-at-least-32-chars-long", "1h")

	// Create handler with invalid expiration format
	handler := NewAuthHandler(
		mockUserRepo,
		mockRoleRepo,
		jwtManager,
		"testgen_token",
		"",
		"/",
		"Lax",
		"invalid-duration",
		false,
		true,
	)

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		handler.setCookie(c, "test-token")
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "testgen_token" {
			foundCookie = cookie
			break
		}
	}

	assert.NotNil(t, foundCookie)
	// Should fallback to 24 hours = 86400 seconds
	expectedMaxAge := int((24 * time.Hour).Seconds())
	assert.Equal(t, expectedMaxAge, foundCookie.MaxAge)
}

func TestRegister_AssignsStudentRoleByDefault(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo, _ := setupAuthHandler(t)

	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "newuser@example.com").Return(nil, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameStudent).Return(studentRole, nil)

	// Capture the user being created to verify role assignment
	var capturedUser *entity.User
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*entity.User)
		}).
		Return(nil)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/register", handler.Register)

	// Create request without role field
	registerReq := dto.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Verify that the created user has student role assigned
	assert.NotNil(t, capturedUser, "User should be created")
	assert.Equal(t, studentRole.ID, capturedUser.RoleID, "User should have student role ID")

	// Parse response and verify role in response
	var authResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)
	assert.Equal(t, "student", authResp.User.Role, "Response should contain student role")

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRegister_FailsWhenStudentRoleNotFound(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo, _ := setupAuthHandler(t)

	mockUserRepo.On("FindByEmail", mock.Anything, "newuser@example.com").Return(nil, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameStudent).Return(nil, assert.AnError)

	// Setup Fiber app
	app := fiber.New()
	app.Post("/register", handler.Register)

	// Create request
	registerReq := dto.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestLogin_InvalidBodyReturnsBadRequest(t *testing.T) {
    handler, _, _, _ := setupAuthHandler(t)
    app := fiber.New()
    app.Post("/login", handler.Login)

    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid"))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
    handler, mockUserRepo, _, _ := setupAuthHandler(t)
    mockUserRepo.On("FindByEmail", mock.Anything, "missing@example.com").Return(nil, gorm.ErrRecordNotFound)

    app := fiber.New()
    app.Post("/login", handler.Login)

    body, _ := json.Marshal(dto.LoginRequest{Email: "missing@example.com", Password: "pass"})
    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestLogin_InvalidPassword(t *testing.T) {
    handler, mockUserRepo, _, _ := setupAuthHandler(t)

    role := &entity.Role{ID: uuid.New(), Name: entity.RoleNameStudent}
    user := &entity.User{ID: uuid.New(), Email: "user@example.com", FullName: "User", RoleID: role.ID, Role: role}
    require.NoError(t, user.SetPassword("correct"))

    mockUserRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil)

    app := fiber.New()
    app.Post("/login", handler.Login)

    body, _ := json.Marshal(dto.LoginRequest{Email: "user@example.com", Password: "wrong"})
    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
    handler, mockUserRepo, mockRoleRepo, _ := setupAuthHandler(t)

    existing := &entity.User{ID: uuid.New(), Email: "dup@example.com"}
    mockUserRepo.On("FindByEmail", mock.Anything, "dup@example.com").Return(existing, nil)
    mockRoleRepo.AssertNotCalled(t, "FindByName")

    app := fiber.New()
    app.Post("/register", handler.Register)

    body, _ := json.Marshal(dto.RegisterRequest{Email: "dup@example.com", Password: "pass", FullName: "Dup"})
    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestRegister_InvalidBody(t *testing.T) {
    handler, _, _, _ := setupAuthHandler(t)
    app := fiber.New()
    app.Post("/register", handler.Register)

    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{invalid"))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
