package llm

import (
	"context"
	"fmt"
)

// PerplexityStrategy implements LLM strategy for Perplexity API
type PerplexityStrategy struct {
	apiKey string
}

// NewPerplexityStrategy creates a new Perplexity strategy
func NewPerplexityStrategy(apiKey string) *PerplexityStrategy {
	return &PerplexityStrategy{
		apiKey: apiKey,
	}
}

// GenerateQuestions generates questions using Perplexity API
func (s *PerplexityStrategy) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	// TODO: Implement actual Perplexity API call
	// For MVP, return mock data

	if s.apiKey == "" {
		return nil, fmt.Errorf("perplexity API key not configured")
	}

	// Mock implementation - replace with actual API call
	questions := make([]GeneratedQuestion, 0, params.NumQuestions)

	for i := 0; i < params.NumQuestions && i < 3; i++ {
		question := GeneratedQuestion{
			QuestionText: fmt.Sprintf("Generated question %d based on the provided text", i+1),
			QuestionType: SingleChoice,
			Difficulty:   params.Difficulty,
			Answers: []GeneratedAnswer{
				{Text: "Correct answer", IsCorrect: true},
				{Text: "Incorrect answer 1", IsCorrect: false},
				{Text: "Incorrect answer 2", IsCorrect: false},
				{Text: "Incorrect answer 3", IsCorrect: false},
			},
			Explanation: "This is a mock question. Real implementation will use Perplexity API.",
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// GetProviderName returns the provider name
func (s *PerplexityStrategy) GetProviderName() string {
	return "perplexity"
}
