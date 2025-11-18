package llm

import (
	"context"
	"fmt"
)

// YandexGPTStrategy implements LLM strategy for YandexGPT API
type YandexGPTStrategy struct {
	apiKey string
}

// NewYandexGPTStrategy creates a new YandexGPT strategy
func NewYandexGPTStrategy(apiKey string) *YandexGPTStrategy {
	return &YandexGPTStrategy{
		apiKey: apiKey,
	}
}

// GenerateQuestions generates questions using YandexGPT API
func (s *YandexGPTStrategy) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	// TODO: Implement actual YandexGPT API call
	// For MVP, return mock data

	if s.apiKey == "" {
		return nil, fmt.Errorf("yandexgpt API key not configured")
	}

	// Mock implementation
	questions := make([]GeneratedQuestion, 0, params.NumQuestions)

	for i := 0; i < params.NumQuestions && i < 3; i++ {
		question := GeneratedQuestion{
			QuestionText: fmt.Sprintf("Вопрос %d сгенерированный YandexGPT", i+1),
			QuestionType: SingleChoice,
			Difficulty:   params.Difficulty,
			Answers: []GeneratedAnswer{
				{Text: "Правильный ответ", IsCorrect: true},
				{Text: "Неправильный ответ 1", IsCorrect: false},
				{Text: "Неправильный ответ 2", IsCorrect: false},
				{Text: "Неправильный ответ 3", IsCorrect: false},
			},
			Explanation: "Тестовый вопрос из YandexGPT",
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// GetProviderName returns the provider name
func (s *YandexGPTStrategy) GetProviderName() string {
	return "yandexgpt"
}
