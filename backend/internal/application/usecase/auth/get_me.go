package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// GetMeUseCase handles retrieving the current user
type GetMeUseCase struct {
	userRepo repository.UserRepository
}

// NewGetMeUseCase creates a new get me use case
func NewGetMeUseCase(userRepo repository.UserRepository) *GetMeUseCase {
	return &GetMeUseCase{
		userRepo: userRepo,
	}
}

// Execute executes the get me use case
func (uc *GetMeUseCase) Execute(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &dto.UserDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.GetRoleName(),
	}, nil
}
