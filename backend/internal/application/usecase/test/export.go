package test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/exporter"
)

// ExportUseCase handles exporting a test to various formats
type ExportUseCase struct {
	testRepo        repository.TestRepository
	questionRepo    repository.QuestionRepository
	answerRepo      repository.AnswerRepository
	exporterFactory *exporter.ExporterFactory
}

// NewExportUseCase creates a new export use case
func NewExportUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	exporterFactory *exporter.ExporterFactory,
) *ExportUseCase {
	return &ExportUseCase{
		testRepo:        testRepo,
		questionRepo:    questionRepo,
		answerRepo:      answerRepo,
		exporterFactory: exporterFactory,
	}
}

// ExportResult contains export output
type ExportResult struct {
	Content     string
	ContentType string
	FileExt     string
	Title       string
}

// Execute executes the export use case
func (uc *ExportUseCase) Execute(ctx context.Context, testID uuid.UUID, userID uuid.UUID, format string) (*ExportResult, error) {
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("test not found: %w", err)
	}

	if test.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	questions, err := uc.questionRepo.FindByTestID(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve questions: %w", err)
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("test has no questions")
	}

	answersMap := make(map[string][]*entity.Answer)
	for _, q := range questions {
		answers, err := uc.answerRepo.FindByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve answers: %w", err)
		}
		answersMap[q.ID.String()] = answers
	}

	exp, err := uc.exporterFactory.GetExporter(exporter.ExportFormat(format))
	if err != nil {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}

	result, err := exp.Export(test, questions, answersMap)
	if err != nil {
		return nil, fmt.Errorf("failed to export: %w", err)
	}

	return &ExportResult{
		Content:     result.Content,
		ContentType: result.ContentType,
		FileExt:     result.FileExt,
		Title:       test.Title,
	}, nil
}

// GetAvailableFormats returns all available export formats
func (uc *ExportUseCase) GetAvailableFormats() []exporter.FormatInfo {
	return uc.exporterFactory.GetAvailableFormats()
}
