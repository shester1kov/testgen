package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion_IsSingleChoice(t *testing.T) {
	q := &Question{
		QuestionType: QuestionTypeSingleChoice,
	}
	assert.True(t, q.IsSingleChoice())
	assert.False(t, q.IsMultipleChoice())
	assert.False(t, q.IsTrueFalse())
	assert.False(t, q.IsShortAnswer())
}

func TestQuestion_IsMultipleChoice(t *testing.T) {
	q := &Question{
		QuestionType: QuestionTypeMultipleChoice,
	}
	assert.False(t, q.IsSingleChoice())
	assert.True(t, q.IsMultipleChoice())
	assert.False(t, q.IsTrueFalse())
	assert.False(t, q.IsShortAnswer())
}

func TestQuestion_IsTrueFalse(t *testing.T) {
	q := &Question{
		QuestionType: QuestionTypeTrueFalse,
	}
	assert.False(t, q.IsSingleChoice())
	assert.False(t, q.IsMultipleChoice())
	assert.True(t, q.IsTrueFalse())
	assert.False(t, q.IsShortAnswer())
}

func TestQuestion_IsShortAnswer(t *testing.T) {
	q := &Question{
		QuestionType: QuestionTypeShortAnswer,
	}
	assert.False(t, q.IsSingleChoice())
	assert.False(t, q.IsMultipleChoice())
	assert.False(t, q.IsTrueFalse())
	assert.True(t, q.IsShortAnswer())
}

func TestQuestion_TableName(t *testing.T) {
	q := Question{}
	assert.Equal(t, "questions", q.TableName())
}
