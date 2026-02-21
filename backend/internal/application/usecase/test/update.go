package test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/security"
)

// UpdateUseCase handles updating a test
type UpdateUseCase struct {
	testRepo repository.TestRepository
}

// NewUpdateUseCase creates a new update use case
func NewUpdateUseCase(testRepo repository.TestRepository) *UpdateUseCase {
	return &UpdateUseCase{
		testRepo: testRepo,
	}
}

// UpdateParams contains test update parameters
type UpdateParams struct {
	TestID      uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
}

// Execute executes the update test use case
func (uc *UpdateUseCase) Execute(ctx context.Context, params UpdateParams) (*dto.TestResponse, error) {
	test, err := uc.testRepo.FindByID(ctx, params.TestID)
	if err != nil {
		return nil, fmt.Errorf("test not found: %w", err)
	}

	if test.UserID != params.UserID {
		return nil, fmt.Errorf("access denied")
	}

	sanitizedTitle := security.SanitizeInput(params.Title)
	sanitizedDescription := security.SanitizeMultiline(params.Description)

	if params.Title != "" && len(sanitizedTitle) < 3 {
		return nil, fmt.Errorf("title must be at least 3 characters")
	}

	if params.Title != "" {
		test.Title = sanitizedTitle
	}
	test.Description = sanitizedDescription
	test.UpdatedAt = time.Now()

	if err := uc.testRepo.Update(ctx, test); err != nil {
		return nil, fmt.Errorf("failed to update test: %w", err)
	}

	return &dto.TestResponse{
		ID:             test.ID.String(),
		UserID:         test.UserID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   test.MoodleSynced,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
	}, nil
}
