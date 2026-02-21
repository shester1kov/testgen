package test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// ListUseCase handles listing tests with pagination
type ListUseCase struct {
	testRepo repository.TestRepository
	userRepo repository.UserRepository
}

// NewListUseCase creates a new list use case
func NewListUseCase(testRepo repository.TestRepository, userRepo repository.UserRepository) *ListUseCase {
	return &ListUseCase{
		testRepo: testRepo,
		userRepo: userRepo,
	}
}

// Execute executes the list tests use case
func (uc *ListUseCase) Execute(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TestListResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var tests []*entity.Test
	var total int64

	if user.IsAdmin() {
		tests, err = uc.testRepo.FindAll(ctx, pageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tests: %w", err)
		}
		total, err = uc.testRepo.CountAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count tests: %w", err)
		}
	} else {
		tests, err = uc.testRepo.FindByUserID(ctx, userID, pageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tests: %w", err)
		}
		total, err = uc.testRepo.CountByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count tests: %w", err)
		}
	}

	result := make([]dto.TestResponse, len(tests))
	for i, t := range tests {
		var userName *string
		var userEmail *string
		if user.IsAdmin() && t.User.ID != uuid.Nil {
			userName = &t.User.FullName
			userEmail = &t.User.Email
		}

		result[i] = dto.TestResponse{
			ID:             t.ID.String(),
			UserID:         t.UserID.String(),
			UserName:       userName,
			UserEmail:      userEmail,
			Title:          t.Title,
			Description:    t.Description,
			TotalQuestions: t.TotalQuestions,
			Status:         string(t.Status),
			MoodleSynced:   t.MoodleSynced,
			CreatedAt:      t.CreatedAt.Format(time.RFC3339),
		}
	}

	return &dto.TestListResponse{Tests: result, Total: total, Page: page, PageSize: pageSize}, nil
}
