package exporter

import (
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// StepikCSVExporter exports tests to Stepik CSV format
// Format: Question,Answer1,y/n,Answer2,y/n,...
type StepikCSVExporter struct{}

// NewStepikCSVExporter creates a new Stepik CSV exporter
func NewStepikCSVExporter() *StepikCSVExporter {
	return &StepikCSVExporter{}
}

// Format returns the export format identifier
func (e *StepikCSVExporter) Format() ExportFormat {
	return FormatStepikCSV
}

// Name returns human-readable format name
func (e *StepikCSVExporter) Name() string {
	return "Stepik CSV"
}

// Description returns format description
func (e *StepikCSVExporter) Description() string {
	return "Stepik platform test import format (CSV with y/n markers)"
}

// SupportedQuestionTypes returns list of supported question types
func (e *StepikCSVExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
	}
}

// Export converts a test with questions and answers to Stepik CSV format
func (e *StepikCSVExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	var builder strings.Builder

	for _, q := range questions {
		// Skip unsupported question types
		if q.QuestionType != entity.QuestionTypeSingleChoice &&
			q.QuestionType != entity.QuestionTypeMultipleChoice {
			continue
		}

		qAnswers := answers[q.ID.String()]
		if len(qAnswers) == 0 {
			continue
		}

		e.formatQuestion(&builder, q, qAnswers)

	}

	content := strings.TrimSuffix(builder.String(), "\n")

	if content == "" {
		return nil, fmt.Errorf("no compatible questions found for Stepik export")
	}

	return &ExportResult{
		Content:     content,
		ContentType: "text/csv",
		FileExt:     "csv",
		Format:      FormatStepikCSV,
	}, nil
}

// formatQuestion formats a single question with answers for Stepik CSV
func (e *StepikCSVExporter) formatQuestion(builder *strings.Builder, q *entity.Question, answers []*entity.Answer) {
	builder.WriteString("text,")
	builder.WriteString(e.escapeCSV(q.QuestionText))
	builder.WriteString(",_\n")

	// Add each answer with y/n marker
	for _, ans := range answers {
		builder.WriteString("option,")
		builder.WriteString(e.escapeCSV(ans.AnswerText))

		if ans.IsCorrect {
			builder.WriteString(",y\n")
		} else {
			builder.WriteString(",n\n")
		}
	}

}

// escapeCSV escapes a string for CSV format
func (e *StepikCSVExporter) escapeCSV(s string) string {
	s = strings.TrimSpace(s)
	// Replace newlines with spaces
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// If contains comma, quote, or leading/trailing space, wrap in quotes
	needsQuoting := strings.ContainsAny(s, ",\"") || s != strings.TrimSpace(s)

	if needsQuoting {
		// Escape existing quotes by doubling them
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}

	return s
}
