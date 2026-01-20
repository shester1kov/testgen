package exporter

import (
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// AikenExporter exports tests to Aiken format
// Aiken format is a simple text format for multiple choice questions
// Used by Moodle and other LMS
// Format:
// Question text
// A. First answer
// B. Second answer
// C. Third answer
// ANSWER: B
type AikenExporter struct{}

// NewAikenExporter creates a new Aiken exporter
func NewAikenExporter() *AikenExporter {
	return &AikenExporter{}
}

// Format returns the export format identifier
func (e *AikenExporter) Format() ExportFormat {
	return FormatAiken
}

// Name returns human-readable format name
func (e *AikenExporter) Name() string {
	return "Aiken"
}

// Description returns format description
func (e *AikenExporter) Description() string {
	return "Simple text format for multiple choice questions (Moodle compatible)"
}

// SupportedQuestionTypes returns list of supported question types
func (e *AikenExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
	}
}

// Export converts a test with questions and answers to Aiken format
func (e *AikenExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	var builder strings.Builder

	exportedCount := 0
	for _, q := range questions {
		// Aiken only supports single choice
		if q.QuestionType != entity.QuestionTypeSingleChoice {
			continue
		}

		qAnswers := answers[q.ID.String()]
		if len(qAnswers) == 0 || len(qAnswers) > 26 { // Max 26 answers (A-Z)
			continue
		}

		line, err := e.formatQuestion(q, qAnswers)
		if err != nil {
			continue
		}

		if exportedCount > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(line)
		exportedCount++
	}

	content := strings.TrimSpace(builder.String())

	if content == "" {
		return nil, fmt.Errorf("no compatible questions found for Aiken export (only single choice supported)")
	}

	return &ExportResult{
		Content:     content,
		ContentType: "text/plain",
		FileExt:     "txt",
		Format:      FormatAiken,
	}, nil
}

// formatQuestion formats a single question in Aiken format
func (e *AikenExporter) formatQuestion(q *entity.Question, answers []*entity.Answer) (string, error) {
	var builder strings.Builder

	// Question text (must be on single line)
	questionText := e.sanitizeText(q.QuestionText)
	builder.WriteString(questionText)
	builder.WriteString("\n")

	// Answer options
	correctLetter := ""
	for i, ans := range answers {
		letter := string(rune('A' + i))
		answerText := e.sanitizeText(ans.AnswerText)
		builder.WriteString(fmt.Sprintf("%s. %s\n", letter, answerText))

		if ans.IsCorrect {
			correctLetter = letter
		}
	}

	if correctLetter == "" {
		return "", fmt.Errorf("no correct answer found")
	}

	// Correct answer line
	builder.WriteString(fmt.Sprintf("ANSWER: %s", correctLetter))

	return builder.String(), nil
}

// sanitizeText cleans text for Aiken format (no line breaks allowed)
func (e *AikenExporter) sanitizeText(s string) string {
	s = strings.TrimSpace(s)
	// Replace newlines with spaces
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	// Replace multiple spaces with single space
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}
