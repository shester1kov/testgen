package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Seeder handles database seeding
type Seeder struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	config   *config.Config
	logger   *logger.Logger
}

// NewSeeder creates a new seeder instance
func NewSeeder(userRepo repository.UserRepository, roleRepo repository.RoleRepository, cfg *config.Config, log *logger.Logger) *Seeder {
	return &Seeder{
		userRepo: userRepo,
		roleRepo: roleRepo,
		config:   cfg,
		logger:   log,
	}
}

// SeedAdminUser creates the default admin user if it doesn't exist
func (s *Seeder) SeedAdminUser(ctx context.Context) error {
	// Check if admin already exists
	existingAdmin, err := s.userRepo.FindByEmail(ctx, s.config.Admin.Email)
	if err == nil && existingAdmin != nil {
		s.logger.Info("Admin user already exists, skipping seed", zap.String("email", s.config.Admin.Email))
		return nil
	}

	// Get admin role
	adminRole, err := s.roleRepo.FindByName(ctx, entity.RoleNameAdmin)
	if err != nil {
		s.logger.Error("Failed to find admin role", zap.Error(err))
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.config.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash admin password", zap.Error(err))
		return err
	}

	// Create admin user
	admin := &entity.User{
		ID:           uuid.New(),
		Email:        s.config.Admin.Email,
		PasswordHash: string(hashedPassword),
		FullName:     s.config.Admin.FullName,
		RoleID:       adminRole.ID,
	}

	if err := s.userRepo.Create(ctx, admin); err != nil {
		s.logger.Error("Failed to create admin user", zap.Error(err))
		return err
	}

	s.logger.Info("Admin user created successfully",
		zap.String("email", admin.Email),
		zap.String("role_id", admin.RoleID.String()),
	)

	return nil
}

// Seed runs all seeders
func (s *Seeder) Seed(ctx context.Context) error {
	s.logger.Info("Running database seeders...")

	if err := s.SeedAdminUser(ctx); err != nil {
		return err
	}

	s.logger.Info("Database seeding completed successfully")
	return nil
}
