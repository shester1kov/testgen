package llm

import "fmt"

// LLMFactory creates LLM strategies based on provider name
type LLMFactory struct {
	perplexityKey string
	openaiKey     string
	yandexKey     string
}

// NewLLMFactory creates a new LLM factory
func NewLLMFactory(perplexityKey, openaiKey, yandexKey string) *LLMFactory {
	return &LLMFactory{
		perplexityKey: perplexityKey,
		openaiKey:     openaiKey,
		yandexKey:     yandexKey,
	}
}

// CreateStrategy creates an LLM strategy for the specified provider
func (f *LLMFactory) CreateStrategy(provider string) (LLMStrategy, error) {
	switch provider {
	case "perplexity":
		return NewPerplexityStrategy(f.perplexityKey), nil
	case "openai":
		return NewOpenAIStrategy(f.openaiKey), nil
	case "yandexgpt", "yandex":
		return NewYandexGPTStrategy(f.yandexKey), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}

// GetAvailableProviders returns list of available providers
func (f *LLMFactory) GetAvailableProviders() []string {
	providers := make([]string, 0)

	if f.perplexityKey != "" {
		providers = append(providers, "perplexity")
	}
	if f.openaiKey != "" {
		providers = append(providers, "openai")
	}
	if f.yandexKey != "" {
		providers = append(providers, "yandexgpt")
	}

	return providers
}
