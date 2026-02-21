package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// UpdateRoleUseCase handles updating a user's role
type UpdateRoleUseCase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUpdateRoleUseCase creates a new update role use case
func NewUpdateRoleUseCase(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Execute executes the update role use case
func (uc *UpdateRoleUseCase) Execute(ctx context.Context, userID uuid.UUID, roleName string) (*dto.UserDTO, error) {
	// Validate role name
	rn := entity.RoleName(roleName)
	if rn != entity.RoleNameAdmin && rn != entity.RoleNameTeacher && rn != entity.RoleNameStudent {
		return nil, fmt.Errorf("invalid role name: %s", roleName)
	}

	// Find user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Find role by name
	role, err := uc.roleRepo.FindByName(ctx, rn)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Update user role
	user.RoleID = role.ID
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	user.Role = role

	return &dto.UserDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.GetRoleName(),
	}, nil
}
