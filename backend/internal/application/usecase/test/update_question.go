package test

import (
	"context"
	"fmt"
	"github.com/shester1kov/testgen-backend/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/pkg/security"
)

// UpdateQuestionUseCase handles updating a question
type UpdateQuestionUseCase struct {
	testRepo     repository.TestRepository
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	userRepo     repository.UserRepository
}

// NewUpdateQuestionUseCase creates a new update question use case
func NewUpdateQuestionUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
) *UpdateQuestionUseCase {
	return &UpdateQuestionUseCase{
		testRepo:     testRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		userRepo:     userRepo,
	}
}

// UpdateQuestionParams contains question update parameters
type UpdateQuestionParams struct {
	TestID       uuid.UUID
	QuestionID   uuid.UUID
	UserID       uuid.UUID
	QuestionText string
	QuestionType string
	Difficulty   string
	Points       *float64
	Answers      []dto.UpdateAnswerRequest
}

// Execute executes the update question use case
func (uc *UpdateQuestionUseCase) Execute(ctx context.Context, params UpdateQuestionParams) (*dto.QuestionDTO, error) {
	// Verify test ownership
	test, err := uc.testRepo.FindByID(ctx, params.TestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTestNotFound, err)
	}

	user, err := uc.userRepo.FindByID(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if test.UserID != params.UserID && !user.IsAdmin() {
		return nil, domain.ErrAccessDenied
	}

	// Check if question belongs to test
	question, err := uc.questionRepo.FindByID(ctx, params.QuestionID)
	if err != nil || question.TestID != params.TestID {
		return nil, domain.ErrQuestionNotFound
	}

	// Sanitize and validate
	sanitizedQuestionText := security.SanitizeMultiline(params.QuestionText)
	if params.QuestionText != "" && len(sanitizedQuestionText) < 3 {
		return nil, fmt.Errorf("question text must be at least 3 characters")
	}

	// Update fields if provided
	if params.QuestionText != "" {
		question.QuestionText = sanitizedQuestionText
	}
	if params.QuestionType != "" {
		question.QuestionType = entity.QuestionType(params.QuestionType)
	}
	if params.Difficulty != "" {
		question.Difficulty = entity.Difficulty(params.Difficulty)
	}
	if params.Points != nil {
		question.Points = *params.Points
	}
	question.UpdatedAt = time.Now()

	if err := uc.questionRepo.Update(ctx, question); err != nil {
		return nil, fmt.Errorf("failed to update question: %w", err)
	}

	// Update answers if provided
	if len(params.Answers) > 0 {
		oldAnswers, _ := uc.answerRepo.FindByQuestionID(ctx, params.QuestionID)
		for _, oldAnswer := range oldAnswers {
			uc.answerRepo.Delete(ctx, oldAnswer.ID)
		}

		for _, answerReq := range params.Answers {
			sanitizedAnswerText := security.SanitizeInput(answerReq.AnswerText)
			answer := &entity.Answer{
				ID:         uuid.New(),
				QuestionID: params.QuestionID,
				AnswerText: sanitizedAnswerText,
				IsCorrect:  answerReq.IsCorrect,
				OrderNum:   answerReq.OrderNum,
				CreatedAt:  time.Now(),
			}
			if err := uc.answerRepo.Create(ctx, answer); err != nil {
				return nil, fmt.Errorf("failed to create answer: %w", err)
			}
		}
	}

	// Load updated answers
	answers, err := uc.answerRepo.FindByQuestionID(ctx, params.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load answers: %w", err)
	}

	answersDTO := make([]dto.AnswerDTO, len(answers))
	for i, a := range answers {
		answersDTO[i] = dto.AnswerDTO{
			ID:         a.ID.String(),
			AnswerText: a.AnswerText,
			IsCorrect:  a.IsCorrect,
			OrderNum:   a.OrderNum,
		}
	}

	return &dto.QuestionDTO{
		ID:           question.ID.String(),
		QuestionText: question.QuestionText,
		QuestionType: string(question.QuestionType),
		Difficulty:   string(question.Difficulty),
		Points:       question.Points,
		OrderNum:     question.OrderNum,
		Answers:      answersDTO,
	}, nil
}
