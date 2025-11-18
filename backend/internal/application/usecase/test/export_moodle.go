package test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
)

// ExportMoodleUseCase handles Moodle export
type ExportMoodleUseCase struct {
	testRepo     repository.TestRepository
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	xmlExporter  *moodle.MoodleXMLExporter
}

// NewExportMoodleUseCase creates a new export moodle use case
func NewExportMoodleUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	xmlExporter *moodle.MoodleXMLExporter,
) *ExportMoodleUseCase {
	return &ExportMoodleUseCase{
		testRepo:     testRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		xmlExporter:  xmlExporter,
	}
}

// Execute executes the export moodle use case
func (uc *ExportMoodleUseCase) Execute(ctx context.Context, testID uuid.UUID, userID uuid.UUID) (string, error) {
	// Get test
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		return "", fmt.Errorf("test not found: %w", err)
	}

	// Verify ownership
	if test.UserID != userID {
		return "", fmt.Errorf("unauthorized access to test")
	}

	// Get questions
	questions, err := uc.questionRepo.FindByTestID(ctx, testID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve questions: %w", err)
	}

	if len(questions) == 0 {
		return "", fmt.Errorf("test has no questions")
	}

	// Get answers for each question
	answersMap := make(map[string][]*entity.Answer)
	for _, q := range questions {
		answers, err := uc.answerRepo.FindByQuestionID(ctx, q.ID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve answers: %w", err)
		}
		answersMap[q.ID.String()] = answers
	}

	// Export to XML
	xmlContent, err := uc.xmlExporter.Export(test, questions, answersMap)
	if err != nil {
		return "", fmt.Errorf("failed to export XML: %w", err)
	}

	return xmlContent, nil
}
