package moodle

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

func TestExportGeneratesXML(t *testing.T) {
	exporter := NewMoodleXMLExporter()

	testID := uuid.New()
	questionID := uuid.New()
	questions := []*entity.Question{{
		ID:           questionID,
		TestID:       testID,
		QuestionText: " What is 2 + 2? ",
		QuestionType: entity.QuestionTypeSingleChoice,
		Points:       1,
	}, {
		ID:           uuid.New(),
		TestID:       testID,
		QuestionText: "True is correct",
		QuestionType: entity.QuestionTypeTrueFalse,
		Points:       1,
	}}

	answers := map[string][]*entity.Answer{
		questionID.String(): {
			{AnswerText: "4", IsCorrect: true},
			{AnswerText: "3", IsCorrect: false},
		},
		questions[1].ID.String(): {
			{AnswerText: "true", IsCorrect: true},
			{AnswerText: "false", IsCorrect: false},
		},
	}

	xmlContent, err := exporter.Export(&entity.Test{ID: testID}, questions, answers)
	if err != nil {
		t.Fatalf("expected export to succeed, got %v", err)
	}

	if !strings.Contains(xmlContent, "<quiz>") || !strings.Contains(xmlContent, "multichoice") || !strings.Contains(xmlContent, "truefalse") {
		t.Fatalf("exported xml missing expected structures: %s", xmlContent)
	}

	if !strings.Contains(xmlContent, "What is 2 + 2?") {
		t.Fatalf("expected sanitized question text in xml: %s", xmlContent)
	}
}

func TestExportUnsupportedType(t *testing.T) {
	exporter := NewMoodleXMLExporter()
	question := &entity.Question{
		ID:           uuid.New(),
		TestID:       uuid.New(),
		QuestionText: "unsupported",
		QuestionType: "fill_blank",
	}

	_, err := exporter.Export(&entity.Test{}, []*entity.Question{question}, map[string][]*entity.Answer{})
	if err == nil {
		t.Fatalf("expected error for unsupported question type")
	}
}

func TestConvertAnswersSanitize(t *testing.T) {
	exporter := NewMoodleXMLExporter()
	answers := []*entity.Answer{
		{AnswerText: "  correct  ", IsCorrect: true},
		{AnswerText: strings.Repeat("x", 260), IsCorrect: false},
	}

	converted := exporter.convertAnswers(answers)
	if len(converted) != 2 {
		t.Fatalf("expected two converted answers")
	}

	if converted[0].Fraction != 100 || converted[1].Fraction != 0 {
		t.Fatalf("unexpected scoring fractions: %+v", converted)
	}

	if converted[0].Text != "correct" {
		t.Fatalf("expected sanitized text, got %q", converted[0].Text)
	}

	if len(converted[1].Text) != 255 || !strings.HasSuffix(converted[1].Text, "...") {
		t.Fatalf("expected long answer text to be truncated, got length %d", len(converted[1].Text))
	}
}

func TestConvertQuestionVariants(t *testing.T) {
	exporter := NewMoodleXMLExporter()

	multi := &entity.Question{ID: uuid.New(), QuestionText: "multi", QuestionType: entity.QuestionTypeMultipleChoice, Points: 2}
	short := &entity.Question{ID: uuid.New(), QuestionText: "short", QuestionType: entity.QuestionTypeShortAnswer, Points: 1}

	multiConverted, err := exporter.convertQuestion(multi, []*entity.Answer{{AnswerText: "a", IsCorrect: true}, {AnswerText: "b"}})
	if err != nil {
		t.Fatalf("expected multiple choice conversion to succeed: %v", err)
	}
	if multiConverted.Type != "multichoice" || multiConverted.Single == nil || *multiConverted.Single {
		t.Fatalf("expected multichoice with single=false, got %+v", multiConverted)
	}

	shortConverted, err := exporter.convertQuestion(short, []*entity.Answer{{AnswerText: "ans", IsCorrect: true}})
	if err != nil {
		t.Fatalf("expected short answer conversion to succeed: %v", err)
	}
	if shortConverted.Type != "shortanswer" || len(shortConverted.Answers) != 1 {
		t.Fatalf("unexpected short answer conversion: %+v", shortConverted)
	}
}

func TestSanitizeText(t *testing.T) {
	exporter := NewMoodleXMLExporter()
	longText := strings.Repeat("a", 300)

	sanitized := exporter.sanitizeText("  trimmed ")
	if sanitized != "trimmed" {
		t.Fatalf("expected whitespace to be trimmed")
	}

	truncated := exporter.sanitizeText(longText)
	if len(truncated) != 255 || !strings.HasSuffix(truncated, "...") {
		t.Fatalf("expected truncation to 255 chars, got len=%d", len(truncated))
	}
}
