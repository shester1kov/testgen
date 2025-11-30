package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/security"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	jwtManager *utils.JWTManager
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtManager *utils.JWTManager) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtManager: jwtManager,
	}
}

// Execute executes the register use case
func (uc *RegisterUseCase) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Sanitize user input to prevent XSS attacks
	sanitizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	sanitizedFullName := security.SanitizeInput(req.FullName)

	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, sanitizedEmail)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", sanitizedEmail)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Get default student role
	studentRole, err := uc.roleRepo.FindByName(ctx, entity.RoleNameStudent)
	if err != nil {
		return nil, fmt.Errorf("failed to find student role: %w", err)
	}

	// Create new user entity
	user := &entity.User{
		ID:       uuid.New(),
		Email:    sanitizedEmail,
		FullName: sanitizedFullName,
		RoleID:   studentRole.ID,
		Role:     studentRole,
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
	token, err := uc.jwtManager.GenerateToken(user.ID, user.Email, user.GetRoleName())
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
			Role:     user.GetRoleName(),
		},
	}, nil
}
