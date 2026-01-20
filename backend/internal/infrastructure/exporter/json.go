package exporter

import (
	"encoding/json"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// JSONExporter exports tests to JSON format
type JSONExporter struct{}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

// Export exports a test to JSON format
func (e *JSONExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	// Create export structure
	exportData := map[string]interface{}{
		"test": map[string]interface{}{
			"id":              test.ID.String(),
			"title":           test.Title,
			"description":     test.Description,
			"total_questions": test.TotalQuestions,
			"status":          test.Status,
			"moodle_synced":   test.MoodleSynced,
			"moodle_test_id":  test.MoodleTestID,
			"created_at":      test.CreatedAt,
			"updated_at":      test.UpdatedAt,
		},
		"questions": e.formatQuestions(questions, answers),
	}

	// Marshal to JSON with pretty print
	jsonBytes, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, err
	}

	return &ExportResult{
		Content:     string(jsonBytes),
		ContentType: "application/json",
		FileExt:     "json",
		Format:      FormatJSON,
	}, nil
}

// formatQuestions formats questions with their answers
func (e *JSONExporter) formatQuestions(questions []*entity.Question, answers map[string][]*entity.Answer) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(questions))

	for _, q := range questions {
		questionData := map[string]interface{}{
			"id":            q.ID.String(),
			"question_text": q.QuestionText,
			"question_type": q.QuestionType,
			"difficulty":    q.Difficulty,
			"points":        q.Points,
			"order_num":     q.OrderNum,
		}

		// Add answers
		questionAnswers := answers[q.ID.String()]
		answersData := make([]map[string]interface{}, 0, len(questionAnswers))
		for _, ans := range questionAnswers {
			answersData = append(answersData, map[string]interface{}{
				"id":          ans.ID.String(),
				"answer_text": ans.AnswerText,
				"is_correct":  ans.IsCorrect,
				"order_num":   ans.OrderNum,
			})
		}
		questionData["answers"] = answersData

		result = append(result, questionData)
	}

	return result
}

// Format returns the export format
func (e *JSONExporter) Format() ExportFormat {
	return FormatJSON
}

// Name returns the human-readable name of the format
func (e *JSONExporter) Name() string {
	return "JSON"
}

// Description returns the description of the format
func (e *JSONExporter) Description() string {
	return "JSON format for data interchange and backup"
}

// SupportedQuestionTypes returns the question types supported by this exporter
func (e *JSONExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
		entity.QuestionTypeTrueFalse,
		entity.QuestionTypeShortAnswer,
	}
}
