package exporter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONExporter_Export(t *testing.T) {
	exporter := NewJSONExporter()

	// Create test data
	testID := uuid.New()
	test := &entity.Test{
		ID:             testID,
		Title:          "Test JSON Export",
		Description:    "Testing JSON export functionality",
		TotalQuestions: 2,
		Status:         entity.TestStatusDraft,
		MoodleSynced:   false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	q1ID := uuid.New()
	q2ID := uuid.New()
	questions := []*entity.Question{
		{
			ID:           q1ID,
			QuestionText: "What is Go?",
			QuestionType: entity.QuestionTypeSingleChoice,
			Difficulty:   entity.DifficultyMedium,
			Points:       1.0,
			OrderNum:     1,
		},
		{
			ID:           q2ID,
			QuestionText: "What is Vue?",
			QuestionType: entity.QuestionTypeSingleChoice,
			Difficulty:   entity.DifficultyEasy,
			Points:       1.0,
			OrderNum:     2,
		},
	}

	answersMap := map[string][]*entity.Answer{
		q1ID.String(): {
			{
				ID:         uuid.New(),
				AnswerText: "A programming language",
				IsCorrect:  true,
				OrderNum:   1,
			},
			{
				ID:         uuid.New(),
				AnswerText: "A database",
				IsCorrect:  false,
				OrderNum:   2,
			},
		},
		q2ID.String(): {
			{
				ID:         uuid.New(),
				AnswerText: "A frontend framework",
				IsCorrect:  true,
				OrderNum:   1,
			},
			{
				ID:         uuid.New(),
				AnswerText: "A backend framework",
				IsCorrect:  false,
				OrderNum:   2,
			},
		},
	}

	// Export
	result, err := exporter.Export(test, questions, answersMap)
	require.NoError(t, err)

	// Verify result metadata
	assert.Equal(t, "application/json", result.ContentType)
	assert.Equal(t, "json", result.FileExt)
	assert.Equal(t, FormatJSON, result.Format)
	assert.NotEmpty(t, result.Content)

	// Verify JSON is valid
	var exportData map[string]interface{}
	err = json.Unmarshal([]byte(result.Content), &exportData)
	require.NoError(t, err)

	// Verify structure
	assert.Contains(t, exportData, "test")
	assert.Contains(t, exportData, "questions")

	testData := exportData["test"].(map[string]interface{})
	assert.Equal(t, "Test JSON Export", testData["title"])
	assert.Equal(t, "Testing JSON export functionality", testData["description"])
	assert.Equal(t, float64(2), testData["total_questions"])

	questionsData := exportData["questions"].([]interface{})
	assert.Len(t, questionsData, 2)

	// Verify first question
	q1Data := questionsData[0].(map[string]interface{})
	assert.Equal(t, "What is Go?", q1Data["question_text"])
	assert.Equal(t, "single_choice", q1Data["question_type"])
	assert.Equal(t, "medium", q1Data["difficulty"])

	q1Answers := q1Data["answers"].([]interface{})
	assert.Len(t, q1Answers, 2)

	ans1 := q1Answers[0].(map[string]interface{})
	assert.Equal(t, "A programming language", ans1["answer_text"])
	assert.Equal(t, true, ans1["is_correct"])
}

func TestJSONExporter_Format(t *testing.T) {
	exporter := NewJSONExporter()
	assert.Equal(t, FormatJSON, exporter.Format())
}

func TestJSONExporter_Name(t *testing.T) {
	exporter := NewJSONExporter()
	assert.Equal(t, "JSON", exporter.Name())
}

func TestJSONExporter_Description(t *testing.T) {
	exporter := NewJSONExporter()
	assert.NotEmpty(t, exporter.Description())
}

func TestJSONExporter_SupportedQuestionTypes(t *testing.T) {
	exporter := NewJSONExporter()
	types := exporter.SupportedQuestionTypes()
	assert.Len(t, types, 4)
	assert.Contains(t, types, entity.QuestionTypeSingleChoice)
	assert.Contains(t, types, entity.QuestionTypeMultipleChoice)
	assert.Contains(t, types, entity.QuestionTypeTrueFalse)
	assert.Contains(t, types, entity.QuestionTypeShortAnswer)
}
