package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// AnswerRepository defines the interface for answer data operations
type AnswerRepository interface {
	// Create creates a new answer
	Create(ctx context.Context, answer *entity.Answer) error

	// FindByID retrieves an answer by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Answer, error)

	// FindByQuestionID retrieves all answers for a specific question
	FindByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error)

	// Update updates an existing answer
	Update(ctx context.Context, answer *entity.Answer) error

	// Delete deletes an answer by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByQuestionID deletes all answers for a specific question
	DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error
}
