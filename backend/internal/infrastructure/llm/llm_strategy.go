package llm

import (
	"context"
	"fmt"
)

// QuestionType represents the type of question to generate
type QuestionType string

const (
	SingleChoice   QuestionType = "single_choice"
	MultipleChoice QuestionType = "multiple_choice"
	TrueFalse      QuestionType = "true_false"
	ShortAnswer    QuestionType = "short_answer"
)

// GenerationParams holds parameters for question generation
type GenerationParams struct {
	Text           string
	NumQuestions   int
	QuestionTypes  []QuestionType
	Difficulty     string
	Language       string
}

// GeneratedQuestion represents a generated question with answers
type GeneratedQuestion struct {
	QuestionText string
	QuestionType QuestionType
	Difficulty   string
	Answers      []GeneratedAnswer
	Explanation  string
}

// GeneratedAnswer represents a possible answer
type GeneratedAnswer struct {
	Text      string
	IsCorrect bool
}

// LLMStrategy defines the interface for LLM providers (Strategy Pattern)
type LLMStrategy interface {
	GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error)
	GetProviderName() string
}

// LLMContext manages LLM strategy selection
type LLMContext struct {
	strategy LLMStrategy
}

// NewLLMContext creates a new LLM context with the specified strategy
func NewLLMContext(strategy LLMStrategy) *LLMContext {
	return &LLMContext{
		strategy: strategy,
	}
}

// SetStrategy changes the current LLM strategy
func (c *LLMContext) SetStrategy(strategy LLMStrategy) {
	c.strategy = strategy
}

// GenerateQuestions generates questions using the current strategy
func (c *LLMContext) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	if c.strategy == nil {
		return nil, fmt.Errorf("no LLM strategy set")
	}
	return c.strategy.GenerateQuestions(ctx, params)
}

// GetProviderName returns the name of the current provider
func (c *LLMContext) GetProviderName() string {
	if c.strategy == nil {
		return "none"
	}
	return c.strategy.GetProviderName()
}
