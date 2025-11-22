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

type mockUserRepository struct {
	repository.UserRepository
	findByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func TestLoginUseCase_Execute(t *testing.T) {
	jwtManager, err := utils.NewJWTManager("secret", "1h")
	require.NoError(t, err)

	validUser := &entity.User{
		ID:       uuid.New(),
		Email:    "user@example.com",
		FullName: "Test User",
		RoleID:   uuid.New(),
	}
	require.NoError(t, validUser.SetPassword("password123"))

	tests := []struct {
		name        string
		repo        repository.UserRepository
		request     dto.LoginRequest
		expectedErr string
	}{
		{
			name: "successfully logs in with valid credentials",
			repo: &mockUserRepository{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return validUser, nil
			}},
			request: dto.LoginRequest{Email: "user@example.com", Password: "password123"},
		},
		{
			name: "returns error when user not found",
			repo: &mockUserRepository{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return nil, gorm.ErrRecordNotFound
			}},
			request:     dto.LoginRequest{Email: "missing@example.com", Password: "password123"},
			expectedErr: "invalid email or password",
		},
		{
			name: "returns error when password is invalid",
			repo: &mockUserRepository{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return validUser, nil
			}},
			request:     dto.LoginRequest{Email: "user@example.com", Password: "wrong"},
			expectedErr: "invalid email or password",
		},
		{
			name: "propagates repository errors",
			repo: &mockUserRepository{findByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
				return nil, errors.New("database unavailable")
			}},
			request:     dto.LoginRequest{Email: "user@example.com", Password: "password123"},
			expectedErr: "failed to find user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewLoginUseCase(tt.repo, jwtManager)
			resp, err := uc.Execute(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, validUser.Email, resp.User.Email)
			require.NotEmpty(t, resp.Token)
		})
	}
}
