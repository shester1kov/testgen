package exporter

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// MoodleQuiz represents the root element of Moodle XML
type MoodleQuiz struct {
	XMLName   xml.Name        `xml:"quiz"`
	Questions []MoodleQuestion `xml:"question"`
}

// MoodleQuestion represents a Moodle question
type MoodleQuestion struct {
	Type              string        `xml:"type,attr"`
	Name              MoodleName    `xml:"name"`
	QuestionText      MoodleText    `xml:"questiontext"`
	GeneralFeedback   MoodleText    `xml:"generalfeedback"`
	DefaultGrade      float64       `xml:"defaultgrade"`
	Penalty           float64       `xml:"penalty"`
	Hidden            int           `xml:"hidden"`
	Single            *bool         `xml:"single,omitempty"`
	ShuffleAnswers    *bool         `xml:"shuffleanswers,omitempty"`
	AnswerNumbering   *string       `xml:"answernumbering,omitempty"`
	CorrectFeedback   *MoodleText   `xml:"correctfeedback,omitempty"`
	IncorrectFeedback *MoodleText   `xml:"incorrectfeedback,omitempty"`
	Answers           []MoodleAnswer `xml:"answer,omitempty"`
}

// MoodleName represents question name
type MoodleName struct {
	Text string `xml:"text"`
}

// MoodleText represents formatted text in Moodle
type MoodleText struct {
	Text   string `xml:"text"`
	Format string `xml:"format,attr"`
}

// MoodleAnswer represents a question answer
type MoodleAnswer struct {
	Fraction float64    `xml:"fraction,attr"`
	Format   string     `xml:"format,attr"`
	Text     string     `xml:"text"`
	Feedback MoodleText `xml:"feedback"`
}

// MoodleXMLExporter exports tests to Moodle XML format
type MoodleXMLExporter struct{}

// NewMoodleXMLExporter creates a new Moodle XML exporter
func NewMoodleXMLExporter() *MoodleXMLExporter {
	return &MoodleXMLExporter{}
}

// Format returns the export format identifier
func (e *MoodleXMLExporter) Format() ExportFormat {
	return FormatMoodleXML
}

// Name returns human-readable format name
func (e *MoodleXMLExporter) Name() string {
	return "Moodle XML"
}

// Description returns format description
func (e *MoodleXMLExporter) Description() string {
	return "Native Moodle quiz import format (XML)"
}

// SupportedQuestionTypes returns list of supported question types
func (e *MoodleXMLExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
		entity.QuestionTypeTrueFalse,
		entity.QuestionTypeShortAnswer,
	}
}

// Export converts a test with questions and answers to Moodle XML
func (e *MoodleXMLExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	quiz := MoodleQuiz{
		Questions: make([]MoodleQuestion, 0, len(questions)),
	}

	for _, q := range questions {
		moodleQuestion, err := e.convertQuestion(q, answers[q.ID.String()])
		if err != nil {
			return nil, fmt.Errorf("failed to convert question %s: %w", q.ID, err)
		}
		quiz.Questions = append(quiz.Questions, moodleQuestion)
	}

	output, err := xml.MarshalIndent(quiz, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}

	xmlContent := xml.Header + string(output)

	return &ExportResult{
		Content:     xmlContent,
		ContentType: "application/xml",
		FileExt:     "xml",
		Format:      FormatMoodleXML,
	}, nil
}

func (e *MoodleXMLExporter) convertQuestion(q *entity.Question, answers []*entity.Answer) (MoodleQuestion, error) {
	moodleQuestion := MoodleQuestion{
		Name: MoodleName{
			Text: e.sanitizeText(q.QuestionText),
		},
		QuestionText: MoodleText{
			Text:   e.sanitizeText(q.QuestionText),
			Format: "html",
		},
		GeneralFeedback: MoodleText{
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
		moodleQuestion.CorrectFeedback = &MoodleText{Text: "Correct!", Format: "html"}
		moodleQuestion.IncorrectFeedback = &MoodleText{Text: "Incorrect.", Format: "html"}
		moodleQuestion.Answers = e.convertAnswers(answers)

	case entity.QuestionTypeMultipleChoice:
		moodleQuestion.Type = "multichoice"
		single := false
		shuffle := true
		numbering := "abc"
		moodleQuestion.Single = &single
		moodleQuestion.ShuffleAnswers = &shuffle
		moodleQuestion.AnswerNumbering = &numbering
		moodleQuestion.CorrectFeedback = &MoodleText{Text: "Correct!", Format: "html"}
		moodleQuestion.IncorrectFeedback = &MoodleText{Text: "Incorrect.", Format: "html"}
		moodleQuestion.Answers = e.convertAnswers(answers)

	case entity.QuestionTypeTrueFalse:
		moodleQuestion.Type = "truefalse"
		trueAnswer := MoodleAnswer{
			Fraction: 0,
			Format:   "moodle_auto_format",
			Text:     "True",
			Feedback: MoodleText{Text: "", Format: "html"},
		}
		falseAnswer := MoodleAnswer{
			Fraction: 0,
			Format:   "moodle_auto_format",
			Text:     "False",
			Feedback: MoodleText{Text: "", Format: "html"},
		}

		for _, ans := range answers {
			if ans.IsCorrect {
				if strings.ToLower(strings.TrimSpace(ans.AnswerText)) == "true" {
					trueAnswer.Fraction = 100
				} else {
					falseAnswer.Fraction = 100
				}
			}
		}

		moodleQuestion.Answers = []MoodleAnswer{trueAnswer, falseAnswer}

	case entity.QuestionTypeShortAnswer:
		moodleQuestion.Type = "shortanswer"
		moodleQuestion.Answers = e.convertAnswers(answers)

	default:
		return moodleQuestion, fmt.Errorf("unsupported question type: %s", q.QuestionType)
	}

	return moodleQuestion, nil
}

func (e *MoodleXMLExporter) convertAnswers(answers []*entity.Answer) []MoodleAnswer {
	moodleAnswers := make([]MoodleAnswer, 0, len(answers))

	for _, ans := range answers {
		fraction := 0.0
		if ans.IsCorrect {
			fraction = 100.0
		}

		moodleAnswers = append(moodleAnswers, MoodleAnswer{
			Fraction: fraction,
			Format:   "html",
			Text:     e.sanitizeText(ans.AnswerText),
			Feedback: MoodleText{
				Text:   "",
				Format: "html",
			},
		})
	}

	return moodleAnswers
}

func (e *MoodleXMLExporter) sanitizeText(text string) string {
	text = strings.TrimSpace(text)
	if len(text) > 255 {
		text = text[:252] + "..."
	}
	return text
}
