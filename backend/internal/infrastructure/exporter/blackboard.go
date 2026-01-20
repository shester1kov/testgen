package exporter

import (
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// BlackboardExporter exports tests to Blackboard tab-delimited format
// Format: Each question on a line, tab-separated fields
// MC (Multiple Choice): MC<tab>Question<tab>Answer1<tab>correct/incorrect<tab>Answer2<tab>correct/incorrect...
// TF (True/False): TF<tab>Question<tab>true/false
// SA (Short Answer): SA<tab>Question<tab>Answer1<tab>Answer2...
type BlackboardExporter struct{}

// NewBlackboardExporter creates a new Blackboard exporter
func NewBlackboardExporter() *BlackboardExporter {
	return &BlackboardExporter{}
}

// Format returns the export format identifier
func (e *BlackboardExporter) Format() ExportFormat {
	return FormatBlackboard
}

// Name returns human-readable format name
func (e *BlackboardExporter) Name() string {
	return "Blackboard"
}

// Description returns format description
func (e *BlackboardExporter) Description() string {
	return "Blackboard Learn tab-delimited import format"
}

// SupportedQuestionTypes returns list of supported question types
func (e *BlackboardExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
		entity.QuestionTypeTrueFalse,
		entity.QuestionTypeShortAnswer,
	}
}

// Export converts a test with questions and answers to Blackboard format
func (e *BlackboardExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	var builder strings.Builder

	exportedCount := 0
	for _, q := range questions {
		qAnswers := answers[q.ID.String()]

		line, err := e.formatQuestion(q, qAnswers)
		if err != nil {
			continue
		}

		builder.WriteString(line)
		builder.WriteString("\n")
		exportedCount++
	}

	content := strings.TrimSuffix(builder.String(), "\n")

	if content == "" {
		return nil, fmt.Errorf("no compatible questions found for Blackboard export")
	}

	return &ExportResult{
		Content:     content,
		ContentType: "text/tab-separated-values",
		FileExt:     "txt",
		Format:      FormatBlackboard,
	}, nil
}

// formatQuestion formats a single question in Blackboard format
func (e *BlackboardExporter) formatQuestion(q *entity.Question, answers []*entity.Answer) (string, error) {
	var parts []string

	questionText := e.sanitizeText(q.QuestionText)

	switch q.QuestionType {
	case entity.QuestionTypeSingleChoice, entity.QuestionTypeMultipleChoice:
		parts = append(parts, "MC")
		parts = append(parts, questionText)

		for _, ans := range answers {
			parts = append(parts, e.sanitizeText(ans.AnswerText))
			if ans.IsCorrect {
				parts = append(parts, "correct")
			} else {
				parts = append(parts, "incorrect")
			}
		}

	case entity.QuestionTypeTrueFalse:
		parts = append(parts, "TF")
		parts = append(parts, questionText)

		// Find the correct answer
		isTrue := false
		for _, ans := range answers {
			if ans.IsCorrect {
				if strings.ToLower(strings.TrimSpace(ans.AnswerText)) == "true" {
					isTrue = true
				}
				break
			}
		}
		if isTrue {
			parts = append(parts, "true")
		} else {
			parts = append(parts, "false")
		}

	case entity.QuestionTypeShortAnswer:
		parts = append(parts, "SA")
		parts = append(parts, questionText)

		// Add all correct answers
		for _, ans := range answers {
			if ans.IsCorrect {
				parts = append(parts, e.sanitizeText(ans.AnswerText))
			}
		}

	default:
		return "", fmt.Errorf("unsupported question type: %s", q.QuestionType)
	}

	return strings.Join(parts, "\t"), nil
}

// sanitizeText cleans text for Blackboard format (no tabs or newlines)
func (e *BlackboardExporter) sanitizeText(s string) string {
	s = strings.TrimSpace(s)
	// Replace tabs and newlines with spaces
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	// Replace multiple spaces with single space
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}
