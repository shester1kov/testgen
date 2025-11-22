package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockRoleRepository struct {
	repository.RoleRepository
	findByNameFunc func(ctx context.Context, name entity.RoleName) (*entity.Role, error)
}

func (m *mockRoleRepository) FindByName(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(ctx, name)
	}
	return nil, gorm.ErrRecordNotFound
}

type mockUserRepo struct {
	repository.UserRepository
	findByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
	createFunc      func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *entity.User) error { return nil }
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Count(ctx context.Context) (int64, error) { return 0, nil }

func TestRegisterUseCase_Execute(t *testing.T) {
	jwtManager, err := utils.NewJWTManager("secret", "1h")
	require.NoError(t, err)

	studentRole := &entity.Role{ID: uuid.New(), Name: entity.RoleNameStudent}

	tests := []struct {
		name        string
		userRepo    repository.UserRepository
		roleRepo    repository.RoleRepository
		request     dto.RegisterRequest
		expectedErr string
	}{
		{
			name: "successfully registers new user",
			userRepo: &mockUserRepo{
				findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) { return nil, gorm.ErrRecordNotFound },
				createFunc:      func(ctx context.Context, user *entity.User) error { return nil },
			},
			roleRepo: &mockRoleRepository{findByNameFunc: func(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
				return studentRole, nil
			}},
			request: dto.RegisterRequest{Email: "new@example.com", Password: "password123", FullName: "New User"},
		},
		{
			name: "returns error when user already exists",
			userRepo: &mockUserRepo{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return &entity.User{Email: email}, nil
			}},
			roleRepo:    &mockRoleRepository{findByNameFunc: func(ctx context.Context, name entity.RoleName) (*entity.Role, error) { return studentRole, nil }},
			request:     dto.RegisterRequest{Email: "existing@example.com", Password: "password123", FullName: "Existing User"},
			expectedErr: "already exists",
		},
		{
			name: "returns error when role lookup fails",
			userRepo: &mockUserRepo{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return nil, gorm.ErrRecordNotFound
			}},
			roleRepo: &mockRoleRepository{findByNameFunc: func(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
				return nil, errors.New("db error")
			}},
			request:     dto.RegisterRequest{Email: "user@example.com", Password: "password123", FullName: "User"},
			expectedErr: "failed to find student role",
		},
		{
			name: "returns error when creating user fails",
			userRepo: &mockUserRepo{
				findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) { return nil, gorm.ErrRecordNotFound },
				createFunc:      func(ctx context.Context, user *entity.User) error { return errors.New("insert failed") },
			},
			roleRepo:    &mockRoleRepository{findByNameFunc: func(ctx context.Context, name entity.RoleName) (*entity.Role, error) { return studentRole, nil }},
			request:     dto.RegisterRequest{Email: "user@example.com", Password: "password123", FullName: "User"},
			expectedErr: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewRegisterUseCase(tt.userRepo, tt.roleRepo, jwtManager)
			resp, err := uc.Execute(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.request.Email, resp.User.Email)
			require.Equal(t, string(entity.RoleNameStudent), resp.User.Role)
			require.NotEmpty(t, resp.Token)
		})
	}
}
