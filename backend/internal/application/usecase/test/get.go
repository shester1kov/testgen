package test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// GetUseCase handles getting a test by ID with questions and answers
type GetUseCase struct {
	testRepo     repository.TestRepository
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	userRepo     repository.UserRepository
}

// NewGetUseCase creates a new get use case
func NewGetUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
) *GetUseCase {
	return &GetUseCase{
		testRepo:     testRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		userRepo:     userRepo,
	}
}

// Execute executes the get test use case
func (uc *GetUseCase) Execute(ctx context.Context, testID uuid.UUID, userID uuid.UUID) (*dto.TestResponse, error) {
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("test not found: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if test.UserID != userID && !user.IsAdmin() {
		return nil, fmt.Errorf("test not found")
	}

	questions, err := uc.questionRepo.FindByTestID(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}

	questionsDTO := make([]dto.QuestionDTO, len(questions))
	for i, q := range questions {
		answers, err := uc.answerRepo.FindByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load answers: %w", err)
		}

		answersDTO := make([]dto.AnswerDTO, len(answers))
		for j, a := range answers {
			answersDTO[j] = dto.AnswerDTO{
				ID:         a.ID.String(),
				AnswerText: a.AnswerText,
				IsCorrect:  a.IsCorrect,
				OrderNum:   a.OrderNum,
			}
		}

		questionsDTO[i] = dto.QuestionDTO{
			ID:           q.ID.String(),
			QuestionText: q.QuestionText,
			QuestionType: string(q.QuestionType),
			Difficulty:   string(q.Difficulty),
			Points:       q.Points,
			OrderNum:     q.OrderNum,
			Answers:      answersDTO,
		}
	}

	return &dto.TestResponse{
		ID:             test.ID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   test.MoodleSynced,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
		Questions:      questionsDTO,
	}, nil
}
