package exporter

import (
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// GIFTExporter exports tests to GIFT format
// GIFT (General Import Format Technology) is used by Moodle and other LMS
// Format spec: https://docs.moodle.org/en/GIFT_format
type GIFTExporter struct{}

// NewGIFTExporter creates a new GIFT exporter
func NewGIFTExporter() *GIFTExporter {
	return &GIFTExporter{}
}

// Format returns the export format identifier
func (e *GIFTExporter) Format() ExportFormat {
	return FormatGIFT
}

// Name returns human-readable format name
func (e *GIFTExporter) Name() string {
	return "GIFT"
}

// Description returns format description
func (e *GIFTExporter) Description() string {
	return "General Import Format Technology - text-based format for Moodle and other LMS"
}

// SupportedQuestionTypes returns list of supported question types
func (e *GIFTExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
		entity.QuestionTypeTrueFalse,
		entity.QuestionTypeShortAnswer,
	}
}

// Export converts a test with questions and answers to GIFT format
func (e *GIFTExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	var builder strings.Builder

	// Add category if test has a title
	if test != nil && test.Title != "" {
		builder.WriteString(fmt.Sprintf("$CATEGORY: %s\n\n", e.escapeGIFT(test.Title)))
	}

	for i, q := range questions {
		qAnswers := answers[q.ID.String()]

		line, err := e.formatQuestion(q, qAnswers, i+1)
		if err != nil {
			continue // Skip unsupported questions
		}

		builder.WriteString(line)
		builder.WriteString("\n\n")
	}

	content := strings.TrimSpace(builder.String())

	if content == "" {
		return nil, fmt.Errorf("no compatible questions found for GIFT export")
	}

	return &ExportResult{
		Content:     content,
		ContentType: "text/plain",
		FileExt:     "txt",
		Format:      FormatGIFT,
	}, nil
}

// formatQuestion formats a single question in GIFT format
func (e *GIFTExporter) formatQuestion(q *entity.Question, answers []*entity.Answer, num int) (string, error) {
	var builder strings.Builder

	// Question title (optional)
	title := fmt.Sprintf("Q%d", num)
	builder.WriteString(fmt.Sprintf("::%s::", e.escapeGIFT(title)))

	// Question text
	builder.WriteString(e.escapeGIFT(q.QuestionText))

	switch q.QuestionType {
	case entity.QuestionTypeSingleChoice:
		builder.WriteString(" {\n")
		for _, ans := range answers {
			if ans.IsCorrect {
				builder.WriteString(fmt.Sprintf("  =%s\n", e.escapeGIFT(ans.AnswerText)))
			} else {
				builder.WriteString(fmt.Sprintf("  ~%s\n", e.escapeGIFT(ans.AnswerText)))
			}
		}
		builder.WriteString("}")

	case entity.QuestionTypeMultipleChoice:
		builder.WriteString(" {\n")
		// For multiple choice, calculate weights
		correctCount := 0
		for _, ans := range answers {
			if ans.IsCorrect {
				correctCount++
			}
		}

		if correctCount == 0 {
			return "", fmt.Errorf("no correct answers")
		}

		correctWeight := 100.0 / float64(correctCount)
		incorrectWeight := -100.0 / float64(len(answers)-correctCount)

		for _, ans := range answers {
			if ans.IsCorrect {
				builder.WriteString(fmt.Sprintf("  ~%%%.0f%%%s\n", correctWeight, e.escapeGIFT(ans.AnswerText)))
			} else {
				builder.WriteString(fmt.Sprintf("  ~%%%.0f%%%s\n", incorrectWeight, e.escapeGIFT(ans.AnswerText)))
			}
		}
		builder.WriteString("}")

	case entity.QuestionTypeTrueFalse:
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
			builder.WriteString(" {TRUE}")
		} else {
			builder.WriteString(" {FALSE}")
		}

	case entity.QuestionTypeShortAnswer:
		builder.WriteString(" {\n")
		for _, ans := range answers {
			if ans.IsCorrect {
				builder.WriteString(fmt.Sprintf("  =%s\n", e.escapeGIFT(ans.AnswerText)))
			}
		}
		builder.WriteString("}")

	default:
		return "", fmt.Errorf("unsupported question type: %s", q.QuestionType)
	}

	return builder.String(), nil
}

// escapeGIFT escapes special characters for GIFT format
func (e *GIFTExporter) escapeGIFT(s string) string {
	s = strings.TrimSpace(s)

	// GIFT special characters that need escaping: ~ = # { } :
	replacer := strings.NewReplacer(
		"~", "\\~",
		"=", "\\=",
		"#", "\\#",
		"{", "\\{",
		"}", "\\}",
		":", "\\:",
	)

	return replacer.Replace(s)
}
