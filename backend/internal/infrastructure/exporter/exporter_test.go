package exporter

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for creating test data

func createTestData() (*entity.Test, []*entity.Question, map[string][]*entity.Answer) {
	testID := uuid.New()
	q1ID := uuid.New()
	q2ID := uuid.New()
	q3ID := uuid.New()

	test := &entity.Test{
		ID:          testID,
		UserID:      uuid.New(),
		Title:       "Sample Test",
		Description: "A test for export testing",
		Status:      entity.TestStatusDraft,
	}

	questions := []*entity.Question{
		{
			ID:           q1ID,
			TestID:       testID,
			QuestionText: "What is 2 + 2?",
			QuestionType: entity.QuestionTypeSingleChoice,
			Difficulty:   entity.DifficultyEasy,
			Points:       1.0,
			OrderNum:     1,
		},
		{
			ID:           q2ID,
			TestID:       testID,
			QuestionText: "Is the sky blue?",
			QuestionType: entity.QuestionTypeTrueFalse,
			Difficulty:   entity.DifficultyEasy,
			Points:       1.0,
			OrderNum:     2,
		},
		{
			ID:           q3ID,
			TestID:       testID,
			QuestionText: "Select all prime numbers",
			QuestionType: entity.QuestionTypeMultipleChoice,
			Difficulty:   entity.DifficultyMedium,
			Points:       2.0,
			OrderNum:     3,
		},
	}

	answers := map[string][]*entity.Answer{
		q1ID.String(): {
			{ID: uuid.New(), QuestionID: q1ID, AnswerText: "3", IsCorrect: false, OrderNum: 1},
			{ID: uuid.New(), QuestionID: q1ID, AnswerText: "4", IsCorrect: true, OrderNum: 2},
			{ID: uuid.New(), QuestionID: q1ID, AnswerText: "5", IsCorrect: false, OrderNum: 3},
		},
		q2ID.String(): {
			{ID: uuid.New(), QuestionID: q2ID, AnswerText: "True", IsCorrect: true, OrderNum: 1},
			{ID: uuid.New(), QuestionID: q2ID, AnswerText: "False", IsCorrect: false, OrderNum: 2},
		},
		q3ID.String(): {
			{ID: uuid.New(), QuestionID: q3ID, AnswerText: "2", IsCorrect: true, OrderNum: 1},
			{ID: uuid.New(), QuestionID: q3ID, AnswerText: "3", IsCorrect: true, OrderNum: 2},
			{ID: uuid.New(), QuestionID: q3ID, AnswerText: "4", IsCorrect: false, OrderNum: 3},
			{ID: uuid.New(), QuestionID: q3ID, AnswerText: "5", IsCorrect: true, OrderNum: 4},
		},
	}

	return test, questions, answers
}

// Factory tests

func TestExporterFactory_NewExporterFactory(t *testing.T) {
	factory := NewExporterFactory()
	require.NotNil(t, factory)

	formats := factory.GetAvailableFormats()
	assert.Len(t, formats, 7) // json, moodle_xml, stepik_csv, gift, aiken, blackboard, qti
}

func TestExporterFactory_GetExporter(t *testing.T) {
	factory := NewExporterFactory()

	tests := []struct {
		format    ExportFormat
		expectErr bool
	}{
		{FormatJSON, false},
		{FormatMoodleXML, false},
		{FormatStepikCSV, false},
		{FormatGIFT, false},
		{FormatAiken, false},
		{FormatBlackboard, false},
		{FormatQTI, false},
		{ExportFormat("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			exp, err := factory.GetExporter(tt.format)
			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, exp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, exp)
				assert.Equal(t, tt.format, exp.Format())
			}
		})
	}
}

func TestExporterFactory_GetAvailableFormats(t *testing.T) {
	factory := NewExporterFactory()
	formats := factory.GetAvailableFormats()

	formatMap := make(map[ExportFormat]bool)
	for _, f := range formats {
		formatMap[f.Format] = true
		assert.NotEmpty(t, f.Name)
		assert.NotEmpty(t, f.Description)
	}

	assert.True(t, formatMap[FormatJSON])
	assert.True(t, formatMap[FormatMoodleXML])
	assert.True(t, formatMap[FormatStepikCSV])
	assert.True(t, formatMap[FormatGIFT])
	assert.True(t, formatMap[FormatAiken])
	assert.True(t, formatMap[FormatBlackboard])
	assert.True(t, formatMap[FormatQTI])

	// Verify we have all 7 formats
	assert.Len(t, formats, 7)
}

// Moodle XML Exporter tests

func TestMoodleXMLExporter_Export(t *testing.T) {
	exporter := NewMoodleXMLExporter()
	test, questions, answers := createTestData()

	result, err := exporter.Export(test, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatMoodleXML, result.Format)
	assert.Equal(t, "application/xml", result.ContentType)
	assert.Equal(t, "xml", result.FileExt)
	assert.Contains(t, result.Content, "<?xml")
	assert.Contains(t, result.Content, "<quiz>")
	assert.Contains(t, result.Content, "What is 2 + 2?")
}

func TestMoodleXMLExporter_Metadata(t *testing.T) {
	exporter := NewMoodleXMLExporter()
	assert.Equal(t, FormatMoodleXML, exporter.Format())
	assert.Equal(t, "Moodle XML", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.NotEmpty(t, exporter.SupportedQuestionTypes())
}

// Stepik CSV Exporter tests

func TestStepikCSVExporter_Export(t *testing.T) {
	exporter := NewStepikCSVExporter()
	test, questions, answers := createTestData()

	result, err := exporter.Export(test, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatStepikCSV, result.Format)
	assert.Equal(t, "text/csv", result.ContentType)
	assert.Equal(t, "csv", result.FileExt)

	// Check CSV format: question,answer,y/n,...
	lines := strings.Split(result.Content, "\n")
	assert.GreaterOrEqual(t, len(lines), 1)

	// First line should contain question and y/n markers
	assert.Contains(t, lines[0], ",y")
	assert.Contains(t, lines[0], ",n")
}

func TestStepikCSVExporter_Metadata(t *testing.T) {
	exporter := NewStepikCSVExporter()
	assert.Equal(t, FormatStepikCSV, exporter.Format())
	assert.Equal(t, "Stepik CSV", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.NotEmpty(t, exporter.SupportedQuestionTypes())
}

func TestStepikCSVExporter_EscapeCSV(t *testing.T) {
	exporter := NewStepikCSVExporter()

	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"text with, comma", "\"text with, comma\""},
		{"text with \"quote\"", "\"text with \"\"quote\"\"\""},
		{"text\nwith\nnewlines", "text with newlines"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := exporter.escapeCSV(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// GIFT Exporter tests

func TestGIFTExporter_Export(t *testing.T) {
	exporter := NewGIFTExporter()
	test, questions, answers := createTestData()

	result, err := exporter.Export(test, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatGIFT, result.Format)
	assert.Equal(t, "text/plain", result.ContentType)
	assert.Equal(t, "txt", result.FileExt)

	// Check GIFT format markers
	assert.Contains(t, result.Content, "$CATEGORY:")
	assert.Contains(t, result.Content, "::") // Question title markers
	assert.Contains(t, result.Content, "=")  // Correct answer marker
	assert.Contains(t, result.Content, "~")  // Incorrect answer marker
}

func TestGIFTExporter_TrueFalse(t *testing.T) {
	exporter := NewGIFTExporter()

	questions := []*entity.Question{
		{
			ID:           uuid.New(),
			QuestionText: "The sky is blue",
			QuestionType: entity.QuestionTypeTrueFalse,
		},
	}

	answers := map[string][]*entity.Answer{
		questions[0].ID.String(): {
			{AnswerText: "True", IsCorrect: true},
			{AnswerText: "False", IsCorrect: false},
		},
	}

	result, err := exporter.Export(nil, questions, answers)
	require.NoError(t, err)
	assert.Contains(t, result.Content, "{TRUE}")
}

func TestGIFTExporter_Metadata(t *testing.T) {
	exporter := NewGIFTExporter()
	assert.Equal(t, FormatGIFT, exporter.Format())
	assert.Equal(t, "GIFT", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.NotEmpty(t, exporter.SupportedQuestionTypes())
}

// Aiken Exporter tests

func TestAikenExporter_Export(t *testing.T) {
	exporter := NewAikenExporter()

	// Create single choice only (Aiken only supports single choice)
	q1ID := uuid.New()
	questions := []*entity.Question{
		{
			ID:           q1ID,
			QuestionText: "What is 2 + 2?",
			QuestionType: entity.QuestionTypeSingleChoice,
		},
	}

	answers := map[string][]*entity.Answer{
		q1ID.String(): {
			{AnswerText: "3", IsCorrect: false, OrderNum: 1},
			{AnswerText: "4", IsCorrect: true, OrderNum: 2},
			{AnswerText: "5", IsCorrect: false, OrderNum: 3},
		},
	}

	result, err := exporter.Export(nil, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatAiken, result.Format)
	assert.Equal(t, "text/plain", result.ContentType)
	assert.Equal(t, "txt", result.FileExt)

	// Check Aiken format
	assert.Contains(t, result.Content, "What is 2 + 2?")
	assert.Contains(t, result.Content, "A. 3")
	assert.Contains(t, result.Content, "B. 4")
	assert.Contains(t, result.Content, "C. 5")
	assert.Contains(t, result.Content, "ANSWER: B")
}

func TestAikenExporter_SkipsNonSingleChoice(t *testing.T) {
	exporter := NewAikenExporter()

	// Only multiple choice question - should fail
	questions := []*entity.Question{
		{
			ID:           uuid.New(),
			QuestionText: "Select all",
			QuestionType: entity.QuestionTypeMultipleChoice,
		},
	}

	answers := map[string][]*entity.Answer{
		questions[0].ID.String(): {
			{AnswerText: "A", IsCorrect: true},
		},
	}

	_, err := exporter.Export(nil, questions, answers)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no compatible questions")
}

func TestAikenExporter_Metadata(t *testing.T) {
	exporter := NewAikenExporter()
	assert.Equal(t, FormatAiken, exporter.Format())
	assert.Equal(t, "Aiken", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.Len(t, exporter.SupportedQuestionTypes(), 1) // Only single choice
}

// Blackboard Exporter tests

func TestBlackboardExporter_Export(t *testing.T) {
	exporter := NewBlackboardExporter()
	test, questions, answers := createTestData()

	result, err := exporter.Export(test, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatBlackboard, result.Format)
	assert.Equal(t, "text/tab-separated-values", result.ContentType)
	assert.Equal(t, "txt", result.FileExt)

	// Check Blackboard format - tab separated
	lines := strings.Split(result.Content, "\n")
	assert.GreaterOrEqual(t, len(lines), 1)

	// First line should be MC (multiple choice)
	assert.True(t, strings.HasPrefix(lines[0], "MC\t"))
	assert.Contains(t, lines[0], "correct")
	assert.Contains(t, lines[0], "incorrect")
}

func TestBlackboardExporter_TrueFalse(t *testing.T) {
	exporter := NewBlackboardExporter()

	questions := []*entity.Question{
		{
			ID:           uuid.New(),
			QuestionText: "The sky is blue",
			QuestionType: entity.QuestionTypeTrueFalse,
		},
	}

	answers := map[string][]*entity.Answer{
		questions[0].ID.String(): {
			{AnswerText: "True", IsCorrect: true},
			{AnswerText: "False", IsCorrect: false},
		},
	}

	result, err := exporter.Export(nil, questions, answers)
	require.NoError(t, err)
	assert.Contains(t, result.Content, "TF\t")
	assert.Contains(t, result.Content, "\ttrue")
}

func TestBlackboardExporter_Metadata(t *testing.T) {
	exporter := NewBlackboardExporter()
	assert.Equal(t, FormatBlackboard, exporter.Format())
	assert.Equal(t, "Blackboard", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.NotEmpty(t, exporter.SupportedQuestionTypes())
}

// QTI Exporter tests

func TestQTIExporter_Export(t *testing.T) {
	exporter := NewQTIExporter()
	test, questions, answers := createTestData()

	result, err := exporter.Export(test, questions, answers)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, FormatQTI, result.Format)
	assert.Equal(t, "application/xml", result.ContentType)
	assert.Equal(t, "xml", result.FileExt)

	// Check QTI format
	assert.Contains(t, result.Content, "<?xml")
	assert.Contains(t, result.Content, "questestinterop")
	assert.Contains(t, result.Content, "assessment")
	assert.Contains(t, result.Content, "<item")
}

func TestQTIExporter_Metadata(t *testing.T) {
	exporter := NewQTIExporter()
	assert.Equal(t, FormatQTI, exporter.Format())
	assert.Equal(t, "QTI 2.1", exporter.Name())
	assert.NotEmpty(t, exporter.Description())
	assert.NotEmpty(t, exporter.SupportedQuestionTypes())
}

// Edge cases

func TestExporters_EmptyQuestions(t *testing.T) {
	factory := NewExporterFactory()

	formats := []ExportFormat{FormatMoodleXML, FormatStepikCSV, FormatGIFT, FormatAiken, FormatBlackboard, FormatQTI}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			exp, _ := factory.GetExporter(format)
			_, err := exp.Export(nil, []*entity.Question{}, map[string][]*entity.Answer{})

			// All exporters should handle empty questions gracefully
			// Most should return an error
			if format != FormatMoodleXML {
				assert.Error(t, err)
			}
		})
	}
}

func TestExporters_SpecialCharacters(t *testing.T) {
	factory := NewExporterFactory()

	q1ID := uuid.New()
	test := &entity.Test{Title: "Test with <special> & \"chars\""}
	questions := []*entity.Question{
		{
			ID:           q1ID,
			QuestionText: "What is <html> & \"quotes\"?",
			QuestionType: entity.QuestionTypeSingleChoice,
		},
	}
	answers := map[string][]*entity.Answer{
		q1ID.String(): {
			{AnswerText: "Option <a>", IsCorrect: true},
			{AnswerText: "Option & more", IsCorrect: false},
		},
	}

	formats := []ExportFormat{FormatMoodleXML, FormatStepikCSV, FormatGIFT, FormatAiken, FormatBlackboard, FormatQTI}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			exp, _ := factory.GetExporter(format)
			result, err := exp.Export(test, questions, answers)

			// Should handle special characters without error
			require.NoError(t, err)
			assert.NotEmpty(t, result.Content)
		})
	}
}

// Benchmarks

func BenchmarkMoodleXMLExporter(b *testing.B) {
	exporter := NewMoodleXMLExporter()
	test, questions, answers := createTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exporter.Export(test, questions, answers)
	}
}

func BenchmarkStepikCSVExporter(b *testing.B) {
	exporter := NewStepikCSVExporter()
	test, questions, answers := createTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exporter.Export(test, questions, answers)
	}
}

func BenchmarkGIFTExporter(b *testing.B) {
	exporter := NewGIFTExporter()
	test, questions, answers := createTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exporter.Export(test, questions, answers)
	}
}
