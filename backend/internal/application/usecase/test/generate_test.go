package test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/stretchr/testify/require"
)

type mockDocumentRepository struct {
	repository.DocumentRepository
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*entity.Document, error)
}

func (m *mockDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDocumentRepository) Create(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockDocumentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockDocumentRepository) Update(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockDocumentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}

func TestGenerateUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	documentID := uuid.New()

	parsedDocument := &entity.Document{
		ID:         documentID,
		UserID:     userID,
		Status:     entity.StatusParsed,
		ParsedText: "parsed content",
	}

	factory := llm.NewLLMFactory("perplexity-key", "", "", "", "")

	tests := []struct {
		name        string
		repo        repository.DocumentRepository
		params      GenerateParams
		expectedErr string
	}{
		{
			name: "generates questions when document is parsed and owned by user",
			repo: &mockDocumentRepository{findByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
				return parsedDocument, nil
			}},
			params: GenerateParams{UserID: userID, DocumentID: documentID, NumQuestions: 2, Difficulty: "easy"},
		},
		{
			name: "fails when document is missing",
			repo: &mockDocumentRepository{findByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
				return nil, errors.New("not found")
			}},
			params:      GenerateParams{UserID: userID, DocumentID: documentID, NumQuestions: 1},
			expectedErr: "document not found",
		},
		{
			name: "rejects access when user IDs do not match",
			repo: &mockDocumentRepository{findByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
				wrongUserDoc := *parsedDocument
				wrongUserDoc.UserID = uuid.New()
				return &wrongUserDoc, nil
			}},
			params:      GenerateParams{UserID: userID, DocumentID: documentID, NumQuestions: 1},
			expectedErr: "unauthorized access",
		},
		{
			name: "errors when document is not parsed yet",
			repo: &mockDocumentRepository{findByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
				notParsed := *parsedDocument
				notParsed.Status = entity.StatusUploaded
				return &notParsed, nil
			}},
			params:      GenerateParams{UserID: userID, DocumentID: documentID, NumQuestions: 1},
			expectedErr: "document not parsed",
		},
		{
			name: "errors on unknown LLM provider",
			repo: &mockDocumentRepository{findByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
				return parsedDocument, nil
			}},
			params:      GenerateParams{UserID: userID, DocumentID: documentID, NumQuestions: 1, LLMProvider: "unknown"},
			expectedErr: "unknown LLM provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewGenerateUseCase(tt.repo, factory)
			questions, err := uc.Execute(context.Background(), tt.params)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				require.Nil(t, questions)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, questions)
			require.NotEmpty(t, questions)
		})
	}
}
