package document

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserRepository for testing
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

func TestGetUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	logger := zap.NewNop()
	uc := NewGetUseCase(mockRepo, logger)

	docID := uuid.New()
	userID := uuid.New()
	doc := &entity.Document{
		ID:         docID,
		UserID:     userID,
		Title:      "Test Doc",
		FileName:   "test.pdf",
		FileType:   entity.FileTypePDF,
		FileSize:   1024,
		ParsedText: "Sample parsed text",
		Status:     entity.StatusParsed,
	}

	// Success case - without parsed text
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	result, err := uc.Execute(context.Background(), docID, userID, false)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, docID.String(), result.ID)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, "Test Doc", result.Title)
	assert.Nil(t, result.ParsedText) // Should be nil when includeParsedText=false

	// Success case - with parsed text
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	result, err = uc.Execute(context.Background(), docID, userID, true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ParsedText)
	assert.Equal(t, "Sample parsed text", *result.ParsedText)

	// Document not found
	mockRepo.On("FindByID", mock.Anything, docID).Return(nil, errors.New("missing")).Once()
	_, err = uc.Execute(context.Background(), docID, userID, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document not found")

	// Access denied - different user
	otherUser := uuid.New()
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	_, err = uc.Execute(context.Background(), docID, otherUser, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")

	mockRepo.AssertExpectations(t)
}

func TestListUseCase(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	mockUserRepo := new(MockUserRepository)
	logger := zap.NewNop()
	uc := NewListUseCase(mockDocRepo, mockUserRepo, logger)

	userID := uuid.New()
	studentRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameStudent,
	}
	user := &entity.User{
		ID:       userID,
		Email:    "student@test.com",
		FullName: "Test Student",
		RoleID:   studentRole.ID,
		Role:     studentRole,
	}

	doc1 := &entity.Document{
		ID:       uuid.New(),
		UserID:   userID,
		Title:    "Doc 1",
		FileName: "doc1.pdf",
		FileType: entity.FileTypePDF,
		FileSize: 1024,
		Status:   entity.StatusUploaded,
	}
	doc2 := &entity.Document{
		ID:       uuid.New(),
		UserID:   userID,
		Title:    "Doc 2",
		FileName: "doc2.pdf",
		FileType: entity.FileTypePDF,
		FileSize: 2048,
		Status:   entity.StatusParsed,
	}
	documents := []*entity.Document{doc1, doc2}

	// Student user - sees only their documents
	mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil).Once()
	mockDocRepo.On("FindByUserID", mock.Anything, userID, 10, 0).Return(documents, nil).Once()
	mockDocRepo.On("CountByUserID", mock.Anything, userID).Return(int64(2), nil).Once()

	result, err := uc.Execute(context.Background(), userID, 1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Documents, 2)
	assert.EqualValues(t, 2, result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, doc1.ID.String(), result.Documents[0].ID)
	assert.Equal(t, doc2.ID.String(), result.Documents[1].ID)

	mockDocRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestListUseCase_Admin(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	mockUserRepo := new(MockUserRepository)
	logger := zap.NewNop()
	uc := NewListUseCase(mockDocRepo, mockUserRepo, logger)

	adminID := uuid.New()
	adminRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameAdmin,
	}
	admin := &entity.User{
		ID:       adminID,
		Email:    "admin@test.com",
		FullName: "Test Admin",
		RoleID:   adminRole.ID,
		Role:     adminRole,
	}

	otherUserID := uuid.New()
	otherUser := &entity.User{
		ID:       otherUserID,
		Email:    "other@test.com",
		FullName: "Other User",
	}

	doc1 := &entity.Document{
		ID:       uuid.New(),
		UserID:   otherUserID,
		Title:    "Doc 1",
		FileName: "doc1.pdf",
		FileType: entity.FileTypePDF,
		FileSize: 1024,
		Status:   entity.StatusUploaded,
		User:     *otherUser,
	}
	documents := []*entity.Document{doc1}

	// Admin user - sees all documents
	mockUserRepo.On("FindByID", mock.Anything, adminID).Return(admin, nil).Once()
	mockDocRepo.On("FindAll", mock.Anything, 20, 0).Return(documents, nil).Once()
	mockDocRepo.On("CountAll", mock.Anything).Return(int64(1), nil).Once()

	result, err := uc.Execute(context.Background(), adminID, 1, 20)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Documents, 1)
	assert.EqualValues(t, 1, result.Total)
	// Admin should see user info
	assert.NotNil(t, result.Documents[0].UserName)
	assert.NotNil(t, result.Documents[0].UserEmail)
	assert.Equal(t, "Other User", *result.Documents[0].UserName)
	assert.Equal(t, "other@test.com", *result.Documents[0].UserEmail)

	mockDocRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewDeleteUseCase(mockRepo, mockStorage, logger)

	docID := uuid.New()
	userID := uuid.New()
	objectKey := "documents/" + userID.String() + "/test.txt"
	doc := &entity.Document{
		ID:        docID,
		UserID:    userID,
		ObjectKey: objectKey,
		FileName:  "test.txt",
		FileType:  entity.FileTypeTXT,
	}

	// Success case
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	mockStorage.On("Delete", mock.Anything, objectKey).Return(nil).Once()
	mockRepo.On("Delete", mock.Anything, docID).Return(nil).Once()

	err := uc.Execute(context.Background(), docID, userID)
	assert.NoError(t, err)
	mockStorage.AssertCalled(t, "Delete", mock.Anything, objectKey)

	// Document not found
	mockRepo.On("FindByID", mock.Anything, docID).Return(nil, errors.New("missing")).Once()
	assert.Error(t, uc.Execute(context.Background(), docID, userID))

	// Access denied - different user
	otherUser := uuid.New()
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	err = uc.Execute(context.Background(), docID, otherUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestParseUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	parserFactory := parser.NewDocumentParserFactory()
	logger := zap.NewNop()

	uc := NewParseUseCase(mockRepo, mockStorage, parserFactory, logger)

	docID := uuid.New()
	userID := uuid.New()
	objectKey := "documents/" + userID.String() + "/doc.txt"

	// Create a temporary file for reading
	tempFile := filepath.Join(t.TempDir(), "doc.txt")
	os.WriteFile(tempFile, []byte("parsed content"), 0644)

	doc := &entity.Document{
		ID:        docID,
		UserID:    userID,
		ObjectKey: objectKey,
		FileName:  "doc.txt",
		FileType:  entity.FileTypeTXT,
		Status:    entity.StatusUploaded,
	}

	// Success case
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil).Twice()

	// Mock storage download - return the temp file
	file, _ := os.Open(tempFile)
	mockStorage.On("Download", mock.Anything, objectKey).Return(file, nil).Once()

	err := uc.Execute(context.Background(), docID, userID)
	assert.NoError(t, err)

	// Unauthorized
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	mockStorage.ExpectedCalls = nil
	mockStorage.Calls = nil

	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	err = uc.Execute(context.Background(), docID, uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Already parsed
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	mockStorage.ExpectedCalls = nil
	mockStorage.Calls = nil

	parsedDoc := &entity.Document{
		ID:        docID,
		UserID:    userID,
		ObjectKey: objectKey,
		FileName:  "doc.txt",
		FileType:  entity.FileTypeTXT,
		Status:    entity.StatusParsed,
	}
	mockRepo.On("FindByID", mock.Anything, docID).Return(parsedDoc, nil).Once()
	err = uc.Execute(context.Background(), docID, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already parsed")

	// File download failure
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	mockStorage.ExpectedCalls = nil
	mockStorage.Calls = nil

	missingKey := "documents/" + userID.String() + "/missing.txt"
	missing := &entity.Document{
		ID:        docID,
		UserID:    userID,
		ObjectKey: missingKey,
		FileName:  "missing.txt",
		FileType:  entity.FileTypeTXT,
		Status:    entity.StatusUploaded,
	}

	mockRepo.On("FindByID", mock.Anything, docID).Return(missing, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil).Twice()
	mockStorage.On("Download", mock.Anything, missingKey).Return(nil, errors.New("file not found")).Once()

	err = uc.Execute(context.Background(), docID, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download file")
}
