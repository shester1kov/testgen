package llm

import (
	"fmt"

	"go.uber.org/zap"
)

// LLMFactory creates LLM strategies based on provider name
type LLMFactory struct {
	perplexityKey  string
	openaiKey      string
	yandexKey      string
	yandexFolderID string
	yandexModel    string
	logger         *zap.Logger
}

// NewLLMFactory creates a new LLM factory.
// The logger parameter is optional; if omitted, a no-op logger is used.
func NewLLMFactory(perplexityKey, openaiKey, yandexKey, yandexFolderID, yandexModel string, logger ...*zap.Logger) *LLMFactory {
	var l *zap.Logger
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	} else {
		l = zap.NewNop()
	}

	return &LLMFactory{
		perplexityKey:  perplexityKey,
		openaiKey:      openaiKey,
		yandexKey:      yandexKey,
		yandexFolderID: yandexFolderID,
		yandexModel:    yandexModel,
		logger:         l,
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
		return NewYandexGPTStrategy(f.yandexKey, f.yandexFolderID, f.yandexModel, f.logger), nil
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
	if f.yandexKey != "" && f.yandexFolderID != "" {
		providers = append(providers, "yandexgpt")
	}

	return providers
}
