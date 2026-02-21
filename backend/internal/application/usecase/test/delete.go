package test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// DeleteUseCase handles deleting a test
type DeleteUseCase struct {
	testRepo repository.TestRepository
}

// NewDeleteUseCase creates a new delete use case
func NewDeleteUseCase(testRepo repository.TestRepository) *DeleteUseCase {
	return &DeleteUseCase{
		testRepo: testRepo,
	}
}

// Execute executes the delete test use case
func (uc *DeleteUseCase) Execute(ctx context.Context, testID uuid.UUID, userID uuid.UUID) error {
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		return fmt.Errorf("test not found: %w", err)
	}

	if test.UserID != userID {
		return fmt.Errorf("test not found")
	}

	if err := uc.testRepo.Delete(ctx, testID); err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	return nil
}
