package test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/security"
)

// CreateUseCase handles test creation
type CreateUseCase struct {
	testRepo repository.TestRepository
}

// NewCreateUseCase creates a new create use case
func NewCreateUseCase(testRepo repository.TestRepository) *CreateUseCase {
	return &CreateUseCase{
		testRepo: testRepo,
	}
}

// CreateParams contains test creation parameters
type CreateParams struct {
	UserID      uuid.UUID
	Title       string
	Description string
	DocumentID  *string
}

// Execute executes the create test use case
func (uc *CreateUseCase) Execute(ctx context.Context, params CreateParams) (*dto.TestResponse, error) {
	sanitizedTitle := security.SanitizeInput(params.Title)
	sanitizedDescription := security.SanitizeMultiline(params.Description)

	test := &entity.Test{
		UserID:      params.UserID,
		Title:       sanitizedTitle,
		Description: sanitizedDescription,
	}

	if params.DocumentID != nil {
		docID, _ := uuid.Parse(*params.DocumentID)
		test.DocumentID = &docID
	}

	if err := uc.testRepo.Create(ctx, test); err != nil {
		return nil, fmt.Errorf("failed to create test: %w", err)
	}

	return &dto.TestResponse{
		ID:             test.ID.String(),
		UserID:         test.UserID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: 0,
		Status:         string(test.Status),
		CreatedAt:      test.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
