package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(userRepo repository.UserRepository, jwtManager *utils.JWTManager) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Execute executes the register use case
func (uc *RegisterUseCase) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Create new user entity
	user := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		FullName: req.FullName,
		Role:     entity.UserRole(req.Role),
	}

	// Set password hash
	if err := user.SetPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Save user to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := uc.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return auth response
	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	}, nil
}
