package handler

import (
	"context"
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
	"github.com/stretchr/testify/require"
)

// Mock repositories for stats handler
type mockStatsTestRepository struct {
	mock.Mock
}

func (m *mockStatsTestRepository) Create(ctx context.Context, test *entity.Test) error {
	return nil
}
func (m *mockStatsTestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	return nil, nil
}
func (m *mockStatsTestRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error) {
	return nil, nil
}
func (m *mockStatsTestRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Test, error) {
	return nil, nil
}
func (m *mockStatsTestRepository) Update(ctx context.Context, test *entity.Test) error { return nil }
func (m *mockStatsTestRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockStatsTestRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockStatsTestRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockStatsDocumentRepository struct {
	mock.Mock
}

func (m *mockStatsDocumentRepository) Create(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockStatsDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	return nil, nil
}
func (m *mockStatsDocumentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockStatsDocumentRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockStatsDocumentRepository) Update(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockStatsDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockStatsDocumentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockStatsDocumentRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockStatsQuestionRepository struct {
	mock.Mock
}

func (m *mockStatsQuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	return nil
}
func (m *mockStatsQuestionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	return nil, nil
}
func (m *mockStatsQuestionRepository) FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
	return nil, nil
}
func (m *mockStatsQuestionRepository) Update(ctx context.Context, question *entity.Question) error {
	return nil
}
func (m *mockStatsQuestionRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockStatsQuestionRepository) CountByTestID(ctx context.Context, testID uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockStatsQuestionRepository) ReorderQuestions(ctx context.Context, testID uuid.UUID, questionIDs []uuid.UUID) error {
	return nil
}
func (m *mockStatsQuestionRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockStatsQuestionRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockStatsUserRepository struct {
	mock.Mock
}

func (m *mockStatsUserRepository) Create(ctx context.Context, user *entity.User) error { return nil }
func (m *mockStatsUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockStatsUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockStatsUserRepository) Update(ctx context.Context, user *entity.User) error { return nil }
func (m *mockStatsUserRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockStatsUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}
func (m *mockStatsUserRepository) Count(ctx context.Context) (int64, error) { return 0, nil }

// TestGetDashboardStats_AdminUser tests that admin sees all stats
func TestGetDashboardStats_AdminUser(t *testing.T) {
	userID := uuid.New()
	adminRoleID := uuid.New()

	testRepo := new(mockStatsTestRepository)
	documentRepo := new(mockStatsDocumentRepository)
	questionRepo := new(mockStatsQuestionRepository)
	userRepo := new(mockStatsUserRepository)

	// Mock admin user
	adminRole := &entity.Role{
		ID:   adminRoleID,
		Name: "admin",
	}
	adminUser := &entity.User{
		ID:     userID,
		Email:  "admin@test.com",
		RoleID: adminRoleID,
		Role:   adminRole,
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(adminUser, nil)

	// Mock stats for all users (admin sees all)
	testRepo.On("CountAll", mock.Anything).Return(int64(25), nil)
	documentRepo.On("CountAll", mock.Anything).Return(int64(15), nil)
	questionRepo.On("CountAll", mock.Anything).Return(int64(120), nil)

	handler := NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.String()) // Pass as string
		return c.Next()
	})
	app.Get("/stats/dashboard", handler.GetDashboardStats)

	req := httptest.NewRequest(http.MethodGet, "/stats/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.DashboardStatsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Admin should see all stats
	assert.Equal(t, int64(15), response.DocumentsCount)
	assert.Equal(t, int64(25), response.TestsCount)
	assert.Equal(t, int64(120), response.QuestionsCount)

	testRepo.AssertExpectations(t)
	documentRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// TestGetDashboardStats_TeacherUser tests that teacher sees only their own stats
func TestGetDashboardStats_TeacherUser(t *testing.T) {
	userID := uuid.New()
	teacherRoleID := uuid.New()

	testRepo := new(mockStatsTestRepository)
	documentRepo := new(mockStatsDocumentRepository)
	questionRepo := new(mockStatsQuestionRepository)
	userRepo := new(mockStatsUserRepository)

	// Mock teacher user
	teacherRole := &entity.Role{
		ID:   teacherRoleID,
		Name: "teacher",
	}
	teacherUser := &entity.User{
		ID:     userID,
		Email:  "teacher@test.com",
		RoleID: teacherRoleID,
		Role:   teacherRole,
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(teacherUser, nil)

	// Mock stats for this user only
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(8), nil)
	documentRepo.On("CountByUserID", mock.Anything, userID).Return(int64(5), nil)
	questionRepo.On("CountByUserID", mock.Anything, userID).Return(int64(40), nil)

	handler := NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.String()) // Pass as string
		return c.Next()
	})
	app.Get("/stats/dashboard", handler.GetDashboardStats)

	req := httptest.NewRequest(http.MethodGet, "/stats/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.DashboardStatsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Teacher should see only their own stats
	assert.Equal(t, int64(5), response.DocumentsCount)
	assert.Equal(t, int64(8), response.TestsCount)
	assert.Equal(t, int64(40), response.QuestionsCount)

	testRepo.AssertExpectations(t)
	documentRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// TestGetDashboardStats_StudentUser tests that student sees only their own stats
func TestGetDashboardStats_StudentUser(t *testing.T) {
	userID := uuid.New()
	studentRoleID := uuid.New()

	testRepo := new(mockStatsTestRepository)
	documentRepo := new(mockStatsDocumentRepository)
	questionRepo := new(mockStatsQuestionRepository)
	userRepo := new(mockStatsUserRepository)

	// Mock student user
	studentRole := &entity.Role{
		ID:   studentRoleID,
		Name: "student",
	}
	studentUser := &entity.User{
		ID:     userID,
		Email:  "student@test.com",
		RoleID: studentRoleID,
		Role:   studentRole,
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(studentUser, nil)

	// Mock stats for this user only (students typically have 0 documents/tests)
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(0), nil)
	documentRepo.On("CountByUserID", mock.Anything, userID).Return(int64(0), nil)
	questionRepo.On("CountByUserID", mock.Anything, userID).Return(int64(0), nil)

	handler := NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.String()) // Pass as string
		return c.Next()
	})
	app.Get("/stats/dashboard", handler.GetDashboardStats)

	req := httptest.NewRequest(http.MethodGet, "/stats/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.DashboardStatsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Student should see only their own stats (typically 0)
	assert.Equal(t, int64(0), response.DocumentsCount)
	assert.Equal(t, int64(0), response.TestsCount)
	assert.Equal(t, int64(0), response.QuestionsCount)

	testRepo.AssertExpectations(t)
	documentRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// TestGetDashboardStats_Unauthorized tests unauthorized access
func TestGetDashboardStats_Unauthorized(t *testing.T) {
	testRepo := new(mockStatsTestRepository)
	documentRepo := new(mockStatsDocumentRepository)
	questionRepo := new(mockStatsQuestionRepository)
	userRepo := new(mockStatsUserRepository)

	handler := NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)
	app := fiber.New()
	// No userID in context - simulates unauthorized request
	app.Get("/stats/dashboard", handler.GetDashboardStats)

	req := httptest.NewRequest(http.MethodGet, "/stats/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response dto.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, dto.ErrCodeUnauthorized, response.Error.Code)
}

// TestGetDashboardStats_UserNotFound tests when user doesn't exist
func TestGetDashboardStats_UserNotFound(t *testing.T) {
	userID := uuid.New()

	testRepo := new(mockStatsTestRepository)
	documentRepo := new(mockStatsDocumentRepository)
	questionRepo := new(mockStatsQuestionRepository)
	userRepo := new(mockStatsUserRepository)

	// User not found
	userRepo.On("FindByID", mock.Anything, userID).Return(nil, assert.AnError)

	handler := NewStatsHandler(testRepo, documentRepo, questionRepo, userRepo)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.String()) // Pass as string
		return c.Next()
	})
	app.Get("/stats/dashboard", handler.GetDashboardStats)

	req := httptest.NewRequest(http.MethodGet, "/stats/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response dto.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, dto.ErrCodeDatabaseError, response.Error.Code)
	assert.Contains(t, response.Error.Message, "failed to fetch user")

	userRepo.AssertExpectations(t)
}
