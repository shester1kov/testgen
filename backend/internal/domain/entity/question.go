package entity

import (
	"time"

	"github.com/google/uuid"
)

type QuestionType string

const (
	QuestionTypeSingleChoice   QuestionType = "single_choice"
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeShortAnswer    QuestionType = "short_answer"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Question struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TestID       uuid.UUID    `json:"test_id" gorm:"type:uuid;not null;index"`
	QuestionText string       `json:"question_text" gorm:"type:text;not null"`
	QuestionType QuestionType `json:"question_type" gorm:"type:varchar(50);default:'single_choice'"`
	Difficulty   Difficulty   `json:"difficulty" gorm:"type:varchar(50);default:'medium'"`
	Points       float64      `json:"points" gorm:"type:decimal(5,2);default:1.0"`
	OrderNum     int          `json:"order_num" gorm:"not null"`
	CreatedAt    time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Test    Test     `json:"test,omitempty" gorm:"foreignKey:TestID"`
	Answers []Answer `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
}

// TableName specifies the table name for GORM
func (Question) TableName() string {
	return "questions"
}

// IsSingleChoice checks if question is single choice type
func (q *Question) IsSingleChoice() bool {
	return q.QuestionType == QuestionTypeSingleChoice
}

// IsMultipleChoice checks if question is multiple choice type
func (q *Question) IsMultipleChoice() bool {
	return q.QuestionType == QuestionTypeMultipleChoice
}

// IsTrueFalse checks if question is true/false type
func (q *Question) IsTrueFalse() bool {
	return q.QuestionType == QuestionTypeTrueFalse
}

// IsShortAnswer checks if question is short answer type
func (q *Question) IsShortAnswer() bool {
	return q.QuestionType == QuestionTypeShortAnswer
}
