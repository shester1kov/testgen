package stats

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// DashboardUseCase handles retrieving dashboard statistics
type DashboardUseCase struct {
	testRepo     repository.TestRepository
	documentRepo repository.DocumentRepository
	questionRepo repository.QuestionRepository
	userRepo     repository.UserRepository
}

// NewDashboardUseCase creates a new dashboard use case
func NewDashboardUseCase(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	userRepo repository.UserRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		testRepo:     testRepo,
		documentRepo: documentRepo,
		questionRepo: questionRepo,
		userRepo:     userRepo,
	}
}

// Execute executes the dashboard use case
func (uc *DashboardUseCase) Execute(ctx context.Context, userID uuid.UUID) (*dto.DashboardStatsResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	var documentsCount, testsCount, questionsCount int64

	if user.IsAdmin() {
		documentsCount, _ = uc.documentRepo.CountAll(ctx)
		testsCount, _ = uc.testRepo.CountAll(ctx)
		questionsCount, _ = uc.questionRepo.CountAll(ctx)
	} else {
		documentsCount, _ = uc.documentRepo.CountByUserID(ctx, userID)
		testsCount, _ = uc.testRepo.CountByUserID(ctx, userID)
		questionsCount, _ = uc.questionRepo.CountByUserID(ctx, userID)
	}

	return &dto.DashboardStatsResponse{
		DocumentsCount: documentsCount,
		TestsCount:     testsCount,
		QuestionsCount: questionsCount,
	}, nil
}
