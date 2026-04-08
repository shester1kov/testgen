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

// AcademicProfile represents the educational profile of a discipline
type AcademicProfile string

const (
	ProfileTechnological  AcademicProfile = "technological"
	ProfileNaturalScience AcademicProfile = "natural_science"
	ProfileHumanities     AcademicProfile = "humanities"
	ProfileSocialEconomic AcademicProfile = "social_economic"
	ProfileCreative       AcademicProfile = "creative"
	ProfileUniversal      AcademicProfile = "universal"
)

// ExampleAnswer is a single answer option within a ProfileData question example.
type ExampleAnswer struct {
	Text      string
	IsCorrect bool
}

// QuestionExample is a domain-specific example question loaded from the DB for use in prompts.
type QuestionExample struct {
	QuestionType QuestionType
	QuestionText string
	Answers      []ExampleAnswer
	Explanation  string
}

// FormulationExample is a ПЛОХО/ХОРОШО phrasing pair loaded from the DB.
type FormulationExample struct {
	Bad  string
	Good string
}

// ProfileData contains all academic profile configuration loaded from the DB.
// When set on GenerationParams, BuildSystemPrompt uses this data instead of the
// hardcoded fallback functions.
type ProfileData struct {
	Code             AcademicProfile
	Temperature      float64
	Instruction      string // profile block text; empty for universal
	Formulations     []FormulationExample
	QuestionExamples []QuestionExample
}

// GenerationParams holds parameters for question generation
type GenerationParams struct {
	Text                string
	NumQuestions        int
	QuestionTypes       []QuestionType
	Difficulty          string
	Language            string
	MultipleChoiceCount int             // number of multiple_choice questions; the rest are single_choice
	Profile             AcademicProfile // academic profile of the discipline
	// Temperature overrides the hardcoded LLM temperature when set (non-zero).
	// Populated from ProfileData.Temperature by the generate use case.
	Temperature float64
	// ProfileData carries the full profile loaded from the DB.
	// When non-nil, BuildSystemPrompt uses this data instead of the hardcoded fallback.
	ProfileData *ProfileData
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
