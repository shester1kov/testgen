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
)

func TestGetUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	uc := NewGetUseCase(mockRepo)

	docID := uuid.New()
	userID := uuid.New()
	doc := &entity.Document{ID: docID, UserID: userID}

	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	result, err := uc.Execute(context.Background(), docID, userID)
	assert.NoError(t, err)
	assert.Equal(t, doc, result)

	mockRepo.On("FindByID", mock.Anything, docID).Return(nil, errors.New("missing")).Once()
	_, err = uc.Execute(context.Background(), docID, userID)
	assert.Error(t, err)

	otherUser := uuid.New()
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	_, err = uc.Execute(context.Background(), docID, otherUser)
	assert.Error(t, err)
}

func TestListUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	uc := NewListUseCase(mockRepo)

	userID := uuid.New()
	documents := []*entity.Document{{ID: uuid.New()}, {ID: uuid.New()}}

	mockRepo.On("FindByUserID", mock.Anything, userID, 10, 0).Return(documents, nil).Once()
	mockRepo.On("CountByUserID", mock.Anything, userID).Return(int64(len(documents)), nil).Once()

	result, err := uc.Execute(context.Background(), ListParams{UserID: userID, Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, documents, result.Documents)
	assert.EqualValues(t, len(documents), result.Total)
}

func TestDeleteUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	uc := NewDeleteUseCase(mockRepo)

	docID := uuid.New()
	userID := uuid.New()
	doc := &entity.Document{ID: docID, UserID: userID, FilePath: filepath.Join(t.TempDir(), "file.txt")}
	os.WriteFile(doc.FilePath, []byte("content"), 0644)

	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	mockRepo.On("Delete", mock.Anything, docID).Return(nil).Once()

	err := uc.Execute(context.Background(), docID, userID)
	assert.NoError(t, err)
	_, statErr := os.Stat(doc.FilePath)
	assert.True(t, os.IsNotExist(statErr))

	mockRepo.On("FindByID", mock.Anything, docID).Return(nil, errors.New("missing")).Once()
	assert.Error(t, uc.Execute(context.Background(), docID, userID))

	otherUser := uuid.New()
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	err = uc.Execute(context.Background(), docID, otherUser)
	assert.Error(t, err)
}

func TestParseUseCase(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	parserFactory := parser.NewDocumentParserFactory()
	uc := NewParseUseCase(mockRepo, parserFactory)

	docID := uuid.New()
	userID := uuid.New()
	tempFile := filepath.Join(t.TempDir(), "doc.txt")
	os.WriteFile(tempFile, []byte("parsed"), 0644)

	doc := &entity.Document{ID: docID, UserID: userID, FilePath: tempFile, FileType: entity.FileTypeTXT}

	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil).Twice()

	err := uc.Execute(context.Background(), docID, userID)
	assert.NoError(t, err)

	// Unauthorized
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
	err = uc.Execute(context.Background(), docID, uuid.New())
	assert.Error(t, err)

	// Already parsed
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	parsedDoc := &entity.Document{ID: docID, UserID: userID, FilePath: tempFile, FileType: entity.FileTypeTXT, Status: entity.StatusParsed}
	mockRepo.On("FindByID", mock.Anything, docID).Return(parsedDoc, nil).Once()
	err = uc.Execute(context.Background(), docID, userID)
	assert.Error(t, err)

	// File open failure path
	mockRepo.ExpectedCalls = nil
	mockRepo.Calls = nil
	missing := &entity.Document{ID: docID, UserID: userID, FilePath: filepath.Join(t.TempDir(), "missing.txt"), FileType: entity.FileTypeTXT}
	mockRepo.On("FindByID", mock.Anything, docID).Return(missing, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(nil).Twice()
	err = uc.Execute(context.Background(), docID, userID)
	assert.Error(t, err)
}
