package persistence

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock for UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockRoleRepository is a mock for RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

// POSITIVE TEST: Admin user created successfully
func TestSeeder_SeedAdminUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)

	cfg := &config.Config{
		Admin: config.AdminConfig{
			Email:    "admin@test.com",
			Password: "admin123",
			FullName: "Test Admin",
		},
	}

	log, _ := logger.New(logger.Config{Level: "error", OutputPath: "stdout", Format: "console"})
	seeder := NewSeeder(mockUserRepo, mockRoleRepo, cfg, log)

	adminRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameAdmin,
	}

	ctx := context.Background()

	// Expect FindByEmail to return not found (admin doesn't exist)
	mockUserRepo.On("FindByEmail", ctx, "admin@test.com").Return(nil, errors.New("not found"))

	// Expect FindByName to return admin role
	mockRoleRepo.On("FindByName", ctx, entity.RoleNameAdmin).Return(adminRole, nil)

	// Expect Create to be called and succeed
	mockUserRepo.On("Create", ctx, mock.MatchedBy(func(user *entity.User) bool {
		// Verify user fields
		return user.Email == "admin@test.com" &&
			user.FullName == "Test Admin" &&
			user.RoleID == adminRole.ID &&
			bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("admin123")) == nil
	})).Return(nil)

	err := seeder.SeedAdminUser(ctx)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// POSITIVE TEST: Admin user already exists, skip creation
func TestSeeder_SeedAdminUser_AlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)

	cfg := &config.Config{
		Admin: config.AdminConfig{
			Email:    "admin@test.com",
			Password: "admin123",
			FullName: "Test Admin",
		},
	}

	log, _ := logger.New(logger.Config{Level: "error", OutputPath: "stdout", Format: "console"})
	seeder := NewSeeder(mockUserRepo, mockRoleRepo, cfg, log)

	existingAdmin := &entity.User{
		ID:    uuid.New(),
		Email: "admin@test.com",
	}

	ctx := context.Background()

	// Expect FindByEmail to return existing admin
	mockUserRepo.On("FindByEmail", ctx, "admin@test.com").Return(existingAdmin, nil)

	// Should not call FindByName or Create
	err := seeder.SeedAdminUser(ctx)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertNotCalled(t, "FindByName")
	mockUserRepo.AssertNotCalled(t, "Create")
}

// NEGATIVE TEST: Role not found
func TestSeeder_SeedAdminUser_RoleNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)

	cfg := &config.Config{
		Admin: config.AdminConfig{
			Email:    "admin@test.com",
			Password: "admin123",
			FullName: "Test Admin",
		},
	}

	log, _ := logger.New(logger.Config{Level: "error", OutputPath: "stdout", Format: "console"})
	seeder := NewSeeder(mockUserRepo, mockRoleRepo, cfg, log)

	ctx := context.Background()

	// Admin doesn't exist
	mockUserRepo.On("FindByEmail", ctx, "admin@test.com").Return(nil, errors.New("not found"))

	// Role not found
	mockRoleRepo.On("FindByName", ctx, entity.RoleNameAdmin).Return(nil, errors.New("role not found"))

	err := seeder.SeedAdminUser(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role not found")
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// NEGATIVE TEST: Create user fails
func TestSeeder_SeedAdminUser_CreateFails(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)

	cfg := &config.Config{
		Admin: config.AdminConfig{
			Email:    "admin@test.com",
			Password: "admin123",
			FullName: "Test Admin",
		},
	}

	log, _ := logger.New(logger.Config{Level: "error", OutputPath: "stdout", Format: "console"})
	seeder := NewSeeder(mockUserRepo, mockRoleRepo, cfg, log)

	adminRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameAdmin,
	}

	ctx := context.Background()

	mockUserRepo.On("FindByEmail", ctx, "admin@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("FindByName", ctx, entity.RoleNameAdmin).Return(adminRole, nil)
	mockUserRepo.On("Create", ctx, mock.Anything).Return(errors.New("database error"))

	err := seeder.SeedAdminUser(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// POSITIVE TEST: Seed runs all seeders
func TestSeeder_Seed_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)

	cfg := &config.Config{
		Admin: config.AdminConfig{
			Email:    "admin@test.com",
			Password: "admin123",
			FullName: "Test Admin",
		},
	}

	log, _ := logger.New(logger.Config{Level: "error", OutputPath: "stdout", Format: "console"})
	seeder := NewSeeder(mockUserRepo, mockRoleRepo, cfg, log)

	adminRole := &entity.Role{
		ID:   uuid.New(),
		Name: entity.RoleNameAdmin,
	}

	ctx := context.Background()

	mockUserRepo.On("FindByEmail", ctx, "admin@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("FindByName", ctx, entity.RoleNameAdmin).Return(adminRole, nil)
	mockUserRepo.On("Create", ctx, mock.Anything).Return(nil)

	err := seeder.Seed(ctx)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}
