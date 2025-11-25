package llm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// simple mock strategy to avoid API calls
type mockStrategy struct{}

func (mockStrategy) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	return []GeneratedQuestion{{QuestionText: "mock", QuestionType: SingleChoice}}, nil
}

func (mockStrategy) GetProviderName() string { return "mock" }

func TestLLMFactory_CreateStrategy(t *testing.T) {
	factory := NewLLMFactory("perplexity", "openai", "yandex", "folder123", "yandexgpt-lite")

	strategy, err := factory.CreateStrategy("perplexity")
	require.NoError(t, err)
	require.Equal(t, "perplexity", strategy.GetProviderName())

	strategy, err = factory.CreateStrategy("openai")
	require.NoError(t, err)
	require.Equal(t, "openai", strategy.GetProviderName())

	strategy, err = factory.CreateStrategy("yandex")
	require.NoError(t, err)
	require.Equal(t, "yandexgpt", strategy.GetProviderName())

	_, err = factory.CreateStrategy("unknown")
	require.Error(t, err)
}

func TestLLMFactory_GetAvailableProviders(t *testing.T) {
	factory := NewLLMFactory("a", "", "b", "folder", "yandexgpt-lite")
	providers := factory.GetAvailableProviders()
	require.ElementsMatch(t, []string{"perplexity", "yandexgpt"}, providers)
}

func TestLLMContext_GenerateQuestions(t *testing.T) {
	ctx := NewLLMContext(nil)
	_, err := ctx.GenerateQuestions(context.Background(), GenerationParams{})
	require.EqualError(t, err, "no LLM strategy set")

	ctx.SetStrategy(mockStrategy{})
	questions, err := ctx.GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 1})
	require.NoError(t, err)
	require.Len(t, questions, 1)
	require.Equal(t, "mock", questions[0].QuestionText)
	require.Equal(t, "mock", ctx.GetProviderName())
}

func TestStrategies_RequireAPIKeys(t *testing.T) {
	_, err := NewPerplexityStrategy("").GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 1})
	require.Error(t, err)

	_, err = NewOpenAIStrategy("").GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 1})
	require.Error(t, err)

	_, err = NewYandexGPTStrategy("", "", "").GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 1})
	require.Error(t, err)
}

func TestStrategies_ProduceMockData(t *testing.T) {
	// Perplexity and OpenAI still use mock data in tests
	q, err := NewPerplexityStrategy("key").GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 2, Difficulty: "hard"})
	require.NoError(t, err)
	require.Len(t, q, 2)
	require.Equal(t, "perplexity", NewPerplexityStrategy("key").GetProviderName())

	q, err = NewOpenAIStrategy("key").GenerateQuestions(context.Background(), GenerationParams{NumQuestions: 1, Difficulty: "easy"})
	require.NoError(t, err)
	require.Len(t, q, 1)
	require.Equal(t, "openai", NewOpenAIStrategy("key").GetProviderName())
}

func TestYandexGPTStrategy_GetProviderName(t *testing.T) {
	// YandexGPT now makes real API calls, so we only test the provider name
	strategy := NewYandexGPTStrategy("key", "folder123", "yandexgpt-lite")
	require.Equal(t, "yandexgpt", strategy.GetProviderName())
}
