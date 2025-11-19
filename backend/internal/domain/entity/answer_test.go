package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnswer_TableName(t *testing.T) {
	answer := Answer{}
	assert.Equal(t, "answers", answer.TableName())
}

func TestAnswer_IsCorrect(t *testing.T) {
	answer := &Answer{IsCorrect: true}
	assert.True(t, answer.IsCorrect)

	answer.IsCorrect = false
	assert.False(t, answer.IsCorrect)
}
