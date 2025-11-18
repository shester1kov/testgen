package auth

import (
	"context"
	"fmt"

	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(userRepo repository.UserRepository, jwtManager *utils.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, fmt.Errorf("invalid email or password")
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
