package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/document"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock Repositories
type mockDocumentRepository struct {
	mock.Mock
}

func (m *mockDocumentRepository) Create(ctx context.Context, doc *entity.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *mockDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *mockDocumentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *mockDocumentRepository) Update(ctx context.Context, doc *entity.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *mockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockDocumentRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *mockDocumentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockDocumentRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Mock Storage
type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *mockStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockStorage) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *mockStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

// Helper to create multipart upload request
func createUploadRequest(t *testing.T, path, fieldName, fileName, title string) (*http.Request, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add file
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte("test file content"))
	require.NoError(t, err)

	// Add title if provided
	if title != "" {
		err = writer.WriteField("title", title)
		require.NoError(t, err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, path, &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// Tests

func TestDocumentHandler_Upload_Success(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()

	// Setup mocks
	mockStor.On("Upload", mock.Anything, mock.Anything, mock.Anything, int64(17), "application/pdf").Return(nil)
	mockDocRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil)

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/documents", handler.Upload)

	req, err := createUploadRequest(t, "/documents", "file", "test.pdf", "")
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	mockStor.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentHandler_Upload_FileTooLarge(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()

	// Create use cases with small maxFileSize
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 10, logger) // 10 bytes max
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/documents", handler.Upload)

	req, err := createUploadRequest(t, "/documents", "file", "large.pdf", "")
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var errResp dto.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrCodeFileTooLarge, errResp.Error.Code)
}

func TestDocumentHandler_List_Success(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameStudent,
	}
	user := &entity.User{
		ID:     userID,
		Email:  "test@example.com",
		RoleID: studentRole.ID,
		Role:   studentRole,
	}

	documents := []*entity.Document{
		{
			ID:       uuid.New(),
			UserID:   userID,
			Title:    "Doc 1",
			FileName: "doc1.pdf",
			FileType: entity.FileTypePDF,
			FileSize: 1024,
			Status:   entity.StatusUploaded,
		},
	}

	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	mockDocRepo.On("FindByUserID", mock.Anything, userID, 20, 0).Return(documents, nil)
	mockDocRepo.On("CountByUserID", mock.Anything, userID).Return(int64(1), nil)

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Get("/documents", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/documents?page=1&page_size=20", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.DocumentListResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result.Documents, 1)
	assert.Equal(t, int64(1), result.Total)

	mockUserRepo.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentHandler_GetByID_Success(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()
	docID := uuid.New()
	doc := &entity.Document{
		ID:       docID,
		UserID:   userID,
		Title:    "Test Doc",
		FileName: "test.pdf",
		FileType: entity.FileTypePDF,
		FileSize: 1024,
		Status:   entity.StatusUploaded,
	}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Get("/documents/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/documents/"+docID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.DocumentUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, docID.String(), result.ID)

	mockDocRepo.AssertExpectations(t)
}

func TestDocumentHandler_GetByID_Forbidden(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()
	otherUserID := uuid.New()
	docID := uuid.New()
	doc := &entity.Document{
		ID:       docID,
		UserID:   otherUserID,
		Title:    "Test Doc",
		FileName: "test.pdf",
		FileType: entity.FileTypePDF,
		FileSize: 1024,
		Status:   entity.StatusUploaded,
	}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Get("/documents/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/documents/"+docID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

	mockDocRepo.AssertExpectations(t)
}

func TestDocumentHandler_Delete_Success(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()
	docID := uuid.New()
	objectKey := "documents/" + userID.String() + "/test.pdf"
	doc := &entity.Document{
		ID:        docID,
		UserID:    userID,
		ObjectKey: objectKey,
		FileName:  "test.pdf",
		FileType:  entity.FileTypePDF,
	}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)
	mockStor.On("Delete", mock.Anything, objectKey).Return(nil)
	mockDocRepo.On("Delete", mock.Anything, docID).Return(nil)

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Delete("/documents/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/documents/"+docID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "document deleted successfully", result.Message)

	mockDocRepo.AssertExpectations(t)
	mockStor.AssertExpectations(t)
}

func TestDocumentHandler_Delete_NotFound(t *testing.T) {
	mockDocRepo := new(mockDocumentRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	userID := uuid.New()
	docID := uuid.New()

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(nil, errors.New("not found"))

	// Create use cases
	uploadUC := document.NewUploadUseCase(mockDocRepo, mockStor, 50*1024*1024, logger)
	deleteUC := document.NewDeleteUseCase(mockDocRepo, mockStor, logger)
	parseUC := document.NewParseUseCase(mockDocRepo, mockStor, parser.NewDocumentParserFactory(), logger)
	listUC := document.NewListUseCase(mockDocRepo, mockUserRepo, logger)
	getUC := document.NewGetUseCase(mockDocRepo, logger)

	handler := NewDocumentHandler(uploadUC, deleteUC, parseUC, listUC, getUC)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Delete("/documents/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/documents/"+docID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockDocRepo.AssertExpectations(t)
}
