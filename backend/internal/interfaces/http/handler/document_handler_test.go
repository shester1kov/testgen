package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDocumentRepository struct {
	mock.Mock
}

func (m *mockDocumentRepository) Create(ctx context.Context, document *entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

type mockDocUserRepository struct {
	mock.Mock
}

func (m *mockDocUserRepository) Create(ctx context.Context, user *entity.User) error { return nil }
func (m *mockDocUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockDocUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockDocUserRepository) Update(ctx context.Context, user *entity.User) error { return nil }
func (m *mockDocUserRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockDocUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}
func (m *mockDocUserRepository) Count(ctx context.Context) (int64, error) { return 0, nil }

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

func (m *mockDocumentRepository) Update(ctx context.Context, document *entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *mockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockDocumentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockDocumentRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *mockDocumentRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type stubParser struct {
	result string
	err    error
}

func (p stubParser) Parse(reader io.Reader) (string, error) {
	return p.result, p.err
}

func (p stubParser) SupportedType() string { return "txt" }

func createUploadRequest(t *testing.T, path, fieldName, fileName string) (*http.Request, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte("content"))
	assert.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, path, &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestDocumentUpload_SucceedsWithSupportedType(t *testing.T) {
	repo := new(mockDocumentRepository)
	factory := parser.NewDocumentParserFactory()
	uploadDir := t.TempDir()
	userID := uuid.New()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Document")).Run(func(args mock.Arguments) {
		doc := args.Get(1).(*entity.Document)
		doc.ID = uuid.New()
	}).Return(nil)

	handler := NewDocumentHandler(repo, new(mockDocUserRepository), factory, uploadDir, 1024)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/route", handler.Upload)

	req, err := createUploadRequest(t, "/route", "file", "sample.txt")
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var result dto.DocumentUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "sample.txt", result.FileName)

	repo.AssertExpectations(t)
}

func TestDocumentUpload_RejectsUnsupportedOrLargeFiles(t *testing.T) {
	repo := new(mockDocumentRepository)
	factory := parser.NewDocumentParserFactory()
	uploadDir := t.TempDir()
	userID := uuid.New()

	handler := NewDocumentHandler(repo, new(mockDocUserRepository), factory, uploadDir, 4)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/route", handler.Upload)

	req, err := createUploadRequest(t, "/route", "file", "sample.exe")
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Too large
	largeReq, err := createUploadRequest(t, "/route", "file", "big.txt")
	require.NoError(t, err)
	resp, err = app.Test(largeReq)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	repo.AssertNotCalled(t, "Create")
}

func TestDocumentList_HandlesErrors(t *testing.T) {
	repo := new(mockDocumentRepository)
	userRepo := new(mockDocUserRepository)
	factory := parser.NewDocumentParserFactory()
	userID := uuid.New()

	// Mock user with teacher role
	teacherRole := &entity.Role{ID: uuid.New(), Name: "teacher"}
	teacherUser := &entity.User{ID: userID, Email: "teacher@test.com", RoleID: teacherRole.ID, Role: teacherRole}
	userRepo.On("FindByID", mock.Anything, userID).Return(teacherUser, nil)

	repo.On("FindByUserID", mock.Anything, userID, 20, 0).Return(nil, assert.AnError)
	handler := NewDocumentHandler(repo, userRepo, factory, t.TempDir(), 1024)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Get("/route", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/route", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	repo.AssertExpectations(t)
}

func TestDocumentGetByID_ForbiddenAndNotFound(t *testing.T) {
	repo := new(mockDocumentRepository)
	factory := parser.NewDocumentParserFactory()
	userID := uuid.New()
	otherUser := uuid.New()

	doc := &entity.Document{ID: uuid.New(), UserID: otherUser}
	repo.On("FindByID", mock.Anything, doc.ID).Return(doc, nil)

	handler := NewDocumentHandler(repo, new(mockDocUserRepository), factory, t.TempDir(), 1024)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Get("/route/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/route/"+doc.ID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

	missingID := uuid.New()
	repo.On("FindByID", mock.Anything, missingID).Return(nil, assert.AnError)
	req = httptest.NewRequest(http.MethodGet, "/route/"+missingID.String(), nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDocumentDelete_RemovesFileAndHandlesRepoFailure(t *testing.T) {
	repo := new(mockDocumentRepository)
	factory := parser.NewDocumentParserFactory()
	userID := uuid.New()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("data"), 0644))

	doc := &entity.Document{ID: uuid.New(), UserID: userID, FilePath: filePath}
	repo.On("FindByID", mock.Anything, doc.ID).Return(doc, nil)
	repo.On("Delete", mock.Anything, doc.ID).Return(assert.AnError)

	handler := NewDocumentHandler(repo, new(mockDocUserRepository), factory, tempDir, 1024)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Delete("/route/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/route/"+doc.ID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestDocumentParse_SuccessAndParserError(t *testing.T) {
	repo := new(mockDocumentRepository)
	factory := parser.NewDocumentParserFactory()
	parserStub := stubParser{result: "parsed text"}
	factory.Register(parserStub)

	userID := uuid.New()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "doc.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("hello"), 0644))

	doc := &entity.Document{ID: uuid.New(), UserID: userID, FilePath: filePath, FileType: entity.FileTypeTXT}
	repo.On("FindByID", mock.Anything, doc.ID).Return(doc, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil)

	handler := NewDocumentHandler(repo, new(mockDocUserRepository), factory, tempDir, 1024)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/route/:id/parse", handler.Parse)

	req := httptest.NewRequest(http.MethodPost, "/route/"+doc.ID.String()+"/parse", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	repo.AssertCalled(t, "Update", mock.Anything, mock.AnythingOfType("*entity.Document"))

	// parser error branch
	factoryErr := parser.NewDocumentParserFactory()
	failingParser := stubParser{err: assert.AnError}
	factoryErr.Register(failingParser)
	doc.Status = entity.StatusUploaded
	repo.On("FindByID", mock.Anything, doc.ID).Return(doc, nil)

	handler = NewDocumentHandler(repo, new(mockDocUserRepository), factoryErr, tempDir, 1024)
	app = fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Post("/route/:id/parse", handler.Parse)

	req = httptest.NewRequest(http.MethodPost, "/route/"+doc.ID.String()+"/parse", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func getBodyBytes(t *testing.T, resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return body
}
