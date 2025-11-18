package moodle

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// Quiz represents the root element of Moodle XML
type Quiz struct {
	XMLName   xml.Name   `xml:"quiz"`
	Questions []Question `xml:"question"`
}

// Question represents a Moodle question
type Question struct {
	Type              string   `xml:"type,attr"`
	Name              Name     `xml:"name"`
	QuestionText      Text     `xml:"questiontext"`
	GeneralFeedback   Text     `xml:"generalfeedback"`
	DefaultGrade      float64  `xml:"defaultgrade"`
	Penalty           float64  `xml:"penalty"`
	Hidden            int      `xml:"hidden"`
	Single            *bool    `xml:"single,omitempty"`            // For multiple choice
	ShuffleAnswers    *bool    `xml:"shuffleanswers,omitempty"`    // For multiple choice
	AnswerNumbering   *string  `xml:"answernumbering,omitempty"`   // For multiple choice
	CorrectFeedback   *Text    `xml:"correctfeedback,omitempty"`   // For multiple choice
	IncorrectFeedback *Text    `xml:"incorrectfeedback,omitempty"` // For multiple choice
	Answers           []Answer `xml:"answer,omitempty"`
}

// Name represents question name
type Name struct {
	Text string `xml:"text"`
}

// Text represents formatted text in Moodle
type Text struct {
	Text   string `xml:"text"`
	Format string `xml:"format,attr"`
}

// Answer represents a question answer
type Answer struct {
	Fraction float64 `xml:"fraction,attr"`
	Format   string  `xml:"format,attr"`
	Text     string  `xml:"text"`
	Feedback Text    `xml:"feedback"`
}

// MoodleXMLExporter exports tests to Moodle XML format
type MoodleXMLExporter struct{}

// NewMoodleXMLExporter creates a new Moodle XML exporter
func NewMoodleXMLExporter() *MoodleXMLExporter {
	return &MoodleXMLExporter{}
}

// Export converts a test with questions and answers to Moodle XML
func (e *MoodleXMLExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (string, error) {
	quiz := Quiz{
		Questions: make([]Question, 0, len(questions)),
	}

	for _, q := range questions {
		moodleQuestion, err := e.convertQuestion(q, answers[q.ID.String()])
		if err != nil {
			return "", fmt.Errorf("failed to convert question %s: %w", q.ID, err)
		}
		quiz.Questions = append(quiz.Questions, moodleQuestion)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(quiz, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header
	xmlContent := xml.Header + string(output)
	return xmlContent, nil
}

// convertQuestion converts domain question to Moodle question
func (e *MoodleXMLExporter) convertQuestion(q *entity.Question, answers []*entity.Answer) (Question, error) {
	moodleQuestion := Question{
		Name: Name{
			Text: e.sanitizeText(q.QuestionText),
		},
		QuestionText: Text{
			Text:   e.sanitizeText(q.QuestionText),
			Format: "html",
		},
		GeneralFeedback: Text{
			Text:   "",
			Format: "html",
		},
		DefaultGrade: float64(q.Points),
		Penalty:      0.3333333,
		Hidden:       0,
	}

	switch q.QuestionType {
	case entity.QuestionTypeSingleChoice:
		moodleQuestion.Type = "multichoice"
		single := true
		shuffle := true
		numbering := "abc"
		moodleQuestion.Single = &single
		moodleQuestion.ShuffleAnswers = &shuffle
		moodleQuestion.AnswerNumbering = &numbering
		moodleQuestion.CorrectFeedback = &Text{Text: "Correct!", Format: "html"}
		moodleQuestion.IncorrectFeedback = &Text{Text: "Incorrect.", Format: "html"}
		moodleQuestion.Answers = e.convertAnswers(answers)

	case entity.QuestionTypeMultipleChoice:
		moodleQuestion.Type = "multichoice"
		single := false
		shuffle := true
		numbering := "abc"
		moodleQuestion.Single = &single
		moodleQuestion.ShuffleAnswers = &shuffle
		moodleQuestion.AnswerNumbering = &numbering
		moodleQuestion.CorrectFeedback = &Text{Text: "Correct!", Format: "html"}
		moodleQuestion.IncorrectFeedback = &Text{Text: "Incorrect.", Format: "html"}
		moodleQuestion.Answers = e.convertAnswers(answers)

	case entity.QuestionTypeTrueFalse:
		moodleQuestion.Type = "truefalse"
		// For true/false, we need exactly 2 answers
		trueAnswer := Answer{
			Fraction: 0,
			Format:   "moodle_auto_format",
			Text:     "True",
			Feedback: Text{Text: "", Format: "html"},
		}
		falseAnswer := Answer{
			Fraction: 0,
			Format:   "moodle_auto_format",
			Text:     "False",
			Feedback: Text{Text: "", Format: "html"},
		}

		// Determine which is correct
		for _, ans := range answers {
			if ans.IsCorrect {
				if strings.ToLower(strings.TrimSpace(ans.AnswerText)) == "true" {
					trueAnswer.Fraction = 100
				} else {
					falseAnswer.Fraction = 100
				}
			}
		}

		moodleQuestion.Answers = []Answer{trueAnswer, falseAnswer}

	case entity.QuestionTypeShortAnswer:
		moodleQuestion.Type = "shortanswer"
		moodleQuestion.Answers = e.convertAnswers(answers)

	default:
		return moodleQuestion, fmt.Errorf("unsupported question type: %s", q.QuestionType)
	}

	return moodleQuestion, nil
}

// convertAnswers converts domain answers to Moodle answers
func (e *MoodleXMLExporter) convertAnswers(answers []*entity.Answer) []Answer {
	moodleAnswers := make([]Answer, 0, len(answers))

	for _, ans := range answers {
		fraction := 0.0
		if ans.IsCorrect {
			fraction = 100.0
		}

		moodleAnswers = append(moodleAnswers, Answer{
			Fraction: fraction,
			Format:   "html",
			Text:     e.sanitizeText(ans.AnswerText),
			Feedback: Text{
				Text:   "",
				Format: "html",
			},
		})
	}

	return moodleAnswers
}

// sanitizeText cleans text for XML export
func (e *MoodleXMLExporter) sanitizeText(text string) string {
	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Truncate long texts for question names (max 255 chars)
	if len(text) > 255 {
		text = text[:252] + "..."
	}

	return text
}
