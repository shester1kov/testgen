package handler

import (
	"bytes"
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

type mockTestRepository struct{ mock.Mock }

func (m *mockTestRepository) Create(ctx context.Context, test *entity.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *mockTestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Test), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockTestRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error) {
	args := m.Called(ctx, userID, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Test), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockTestRepository) Update(ctx context.Context, test *entity.Test) error { return nil }

func (m *mockTestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTestRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

type mockTestDocRepository struct{ mock.Mock }

func (m *mockTestDocRepository) Create(ctx context.Context, doc *entity.Document) error { return nil }
func (m *mockTestDocRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Document), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestDocRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockTestDocRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, userID, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Document), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestDocRepository) Update(ctx context.Context, doc *entity.Document) error { return nil }
func (m *mockTestDocRepository) Delete(ctx context.Context, id uuid.UUID) error         { return nil }
func (m *mockTestDocRepository) MarkAsParsed(ctx context.Context, id uuid.UUID, parsedText string) error {
	return nil
}
func (m *mockTestDocRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateTest_Success(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	docRepo := new(mockTestDocRepository)

	testRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Test")).Run(func(args mock.Arguments) {
		test := args.Get(1).(*entity.Test)
		test.ID = uuid.New()
	}).Return(nil)

	handler := NewTestHandler(testRepo, docRepo, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests", handler.Create)

	body, _ := json.Marshal(dto.CreateTestRequest{Title: "Sample", Description: "Desc"})
	req := httptest.NewRequest(http.MethodPost, "/tests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	testRepo.AssertExpectations(t)
}

func TestCreateTest_InvalidBody(t *testing.T) {
	handler := NewTestHandler(new(mockTestRepository), new(mockTestDocRepository), nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", uuid.New()); return c.Next() })
	app.Post("/tests", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/tests", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGenerate_DocumentNotFound(t *testing.T) {
	userID := uuid.New()
	docRepo := new(mockTestDocRepository)
	docRepo.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)

	handler := NewTestHandler(new(mockTestRepository), docRepo, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests/generate", handler.Generate)

	body, _ := json.Marshal(dto.GenerateTestRequest{DocumentID: uuid.New().String(), NumQuestions: 1})
	req := httptest.NewRequest(http.MethodPost, "/tests/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestListTests_Success(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	docRepo := new(mockTestDocRepository)

	testRepo.On("FindByUserID", mock.Anything, userID, 20, 0).Return([]*entity.Test{{ID: uuid.New(), Title: "T1", UserID: userID}}, nil)
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(1), nil)

	handler := NewTestHandler(testRepo, docRepo, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/tests", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetByID_NotFound(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	testRepo.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/tests/"+uuid.New().String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteTest_Success(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	testRepo := new(mockTestRepository)

	testRepo.On("FindByID", mock.Anything, testID).Return(&entity.Test{ID: testID, UserID: userID}, nil)
	testRepo.On("Delete", mock.Anything, testID).Return(nil)

	handler := NewTestHandler(testRepo, new(mockDocumentRepository), nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Delete("/tests/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/tests/"+testID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	testRepo.AssertExpectations(t)
}
