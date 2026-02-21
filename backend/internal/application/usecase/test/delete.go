package test

import (
	"context"
	"fmt"
	"github.com/shester1kov/testgen-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// DeleteUseCase handles deleting a test
type DeleteUseCase struct {
	testRepo repository.TestRepository
	userRepo repository.UserRepository
}

// NewDeleteUseCase creates a new delete use case
func NewDeleteUseCase(testRepo repository.TestRepository, userRepo repository.UserRepository) *DeleteUseCase {
	return &DeleteUseCase{
		testRepo: testRepo,
		userRepo: userRepo,
	}
}

// Execute executes the delete test use case
func (uc *DeleteUseCase) Execute(ctx context.Context, testID uuid.UUID, userID uuid.UUID) error {
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrTestNotFound, err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if test.UserID != userID && !user.IsAdmin() {
		return domain.ErrTestNotFound
	}

	if err := uc.testRepo.Delete(ctx, testID); err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	return nil
}
