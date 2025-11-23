package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserHandler(t *testing.T) (*UserHandler, *MockUserRepository, *MockRoleRepository) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	handler := NewUserHandler(mockUserRepo, mockRoleRepo)
	return handler, mockUserRepo, mockRoleRepo
}

func TestListUsers_Success(t *testing.T) {
	handler, mockUserRepo, _ := setupUserHandler(t)

	adminRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameAdmin,
		Description: "Admin role",
	}

	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	users := []*entity.User{
		{
			ID:       uuid.New(),
			Email:    "user1@example.com",
			FullName: "User One",
			RoleID:   adminRole.ID,
			Role:     adminRole,
		},
		{
			ID:       uuid.New(),
			Email:    "user2@example.com",
			FullName: "User Two",
			RoleID:   studentRole.ID,
			Role:     studentRole,
		},
	}

	mockUserRepo.On("List", mock.Anything, 10, 0).Return(users, nil)
	mockUserRepo.On("Count", mock.Anything).Return(int64(2), nil)

	// Setup Fiber app
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users?limit=10&offset=0", nil)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	var listResp dto.UserListResponse
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)
	assert.Len(t, listResp.Users, 2)
	assert.Equal(t, int64(2), listResp.Total)
	assert.Equal(t, 10, listResp.Limit)
	assert.Equal(t, 0, listResp.Offset)
	assert.Equal(t, "admin", listResp.Users[0].Role)
	assert.Equal(t, "student", listResp.Users[1].Role)

	mockUserRepo.AssertExpectations(t)
}

func TestListUsers_DefaultPagination(t *testing.T) {
	handler, mockUserRepo, _ := setupUserHandler(t)

	// Default limit is 10, default offset is 0
	mockUserRepo.On("List", mock.Anything, 10, 0).Return([]*entity.User{}, nil)
	mockUserRepo.On("Count", mock.Anything).Return(int64(0), nil)

	// Setup Fiber app
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	// Create request without pagination params
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	var listResp dto.UserListResponse
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)
	assert.Equal(t, 10, listResp.Limit) // Default limit
	assert.Equal(t, 0, listResp.Offset)  // Default offset

	mockUserRepo.AssertExpectations(t)
}

func TestListUsers_InvalidLimit(t *testing.T) {
	handler, mockUserRepo, _ := setupUserHandler(t)

	// Setup Fiber app
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	tests := []struct {
		name  string
		query string
	}{
		{"negative limit", "/users?limit=-1"},
		{"limit too large", "/users?limit=101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}

	mockUserRepo.AssertNotCalled(t, "List")
	mockUserRepo.AssertNotCalled(t, "Count")
}

func TestListUsers_InvalidOffset(t *testing.T) {
	handler, mockUserRepo, _ := setupUserHandler(t)

	// Setup Fiber app
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users?offset=-1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUserRepo.AssertNotCalled(t, "List")
	mockUserRepo.AssertNotCalled(t, "Count")
}

func TestUpdateUserRole_Success(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	teacherRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameTeacher,
		Description: "Teacher role",
	}

	user := &entity.User{
		ID:       userID,
		Email:    "user@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}

	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameTeacher).Return(teacherRole, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	// Create request
	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	var userDTO dto.UserDTO
	err = json.NewDecoder(resp.Body).Decode(&userDTO)
	assert.NoError(t, err)
	assert.Equal(t, "teacher", userDTO.Role)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestUpdateUserRole_InvalidUUID(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/invalid-uuid/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUserRepo.AssertNotCalled(t, "FindByID")
	mockRoleRepo.AssertNotCalled(t, "FindByName")
	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUserRole_UserNotFound(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()
	mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, assert.AnError)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertNotCalled(t, "FindByName")
	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUserRole_InvalidRoleName(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "invalid_role",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUserRepo.AssertNotCalled(t, "FindByID")
	mockRoleRepo.AssertNotCalled(t, "FindByName")
	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUserRole_RoleNotFound(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	user := &entity.User{
		ID:       userID,
		Email:    "user@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}

	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameTeacher).Return(nil, assert.AnError)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUserRole_UpdateFails(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	teacherRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameTeacher,
		Description: "Teacher role",
	}

	user := &entity.User{
		ID:       userID,
		Email:    "user@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}

	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameTeacher).Return(teacherRole, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(assert.AnError)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// TestUpdateUserRole_UpdatesTimestamp verifies that UpdatedAt timestamp is set when updating role
func TestUpdateUserRole_UpdatesTimestamp(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupUserHandler(t)

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
	}

	teacherRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameTeacher,
		Description: "Teacher role",
	}

	user := &entity.User{
		ID:       userID,
		Email:    "user@example.com",
		FullName: "Test User",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}

	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	mockRoleRepo.On("FindByName", mock.Anything, entity.RoleNameTeacher).Return(teacherRole, nil)

	// Capture the user object passed to Update to verify UpdatedAt was set
	var capturedUser *entity.User
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*entity.User)
		}).
		Return(nil)

	// Setup Fiber app
	app := fiber.New()
	app.Put("/users/:id/role", handler.UpdateUserRole)

	// Create request
	updateReq := dto.UpdateUserRoleRequest{
		RoleName: "teacher",
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify that UpdatedAt was set and role_id was updated
	assert.NotNil(t, capturedUser)
	assert.Equal(t, teacherRole.ID, capturedUser.RoleID, "RoleID should be updated to teacher role")
	assert.False(t, capturedUser.UpdatedAt.IsZero(), "UpdatedAt should be set")

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}
