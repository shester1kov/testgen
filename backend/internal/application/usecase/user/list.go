package user

import (
	"context"
	"fmt"

	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// ListUseCase handles listing users with pagination
type ListUseCase struct {
	userRepo repository.UserRepository
}

// NewListUseCase creates a new list use case
func NewListUseCase(userRepo repository.UserRepository) *ListUseCase {
	return &ListUseCase{
		userRepo: userRepo,
	}
}

// Execute executes the list users use case
func (uc *ListUseCase) Execute(ctx context.Context, limit, offset int) (*dto.UserListResponse, error) {
	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = dto.UserDTO{
			ID:       u.ID.String(),
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.GetRoleName(),
		}
	}

	return &dto.UserListResponse{
		Users:  userDTOs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
