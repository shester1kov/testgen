package entity

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:uuid;not null;index"`
	AnswerText string    `json:"answer_text" gorm:"type:text;not null"`
	IsCorrect  bool      `json:"is_correct" gorm:"default:false"`
	OrderNum   int       `json:"order_num" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	Question Question `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

// TableName specifies the table name for GORM
func (Answer) TableName() string {
	return "answers"
}

// MarkCorrect marks answer as correct
func (a *Answer) MarkCorrect() {
	a.IsCorrect = true
}

// MarkIncorrect marks answer as incorrect
func (a *Answer) MarkIncorrect() {
	a.IsCorrect = false
}
