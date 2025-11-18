package test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
)

// GenerateUseCase handles test generation using LLM
type GenerateUseCase struct {
	documentRepo repository.DocumentRepository
	llmFactory   *llm.LLMFactory
}

// NewGenerateUseCase creates a new generate use case
func NewGenerateUseCase(documentRepo repository.DocumentRepository, llmFactory *llm.LLMFactory) *GenerateUseCase {
	return &GenerateUseCase{
		documentRepo: documentRepo,
		llmFactory:   llmFactory,
	}
}

// GenerateParams contains generation parameters
type GenerateParams struct {
	UserID       uuid.UUID
	DocumentID   uuid.UUID
	NumQuestions int
	Difficulty   string
	LLMProvider  string
}

// Execute executes the generate use case
func (uc *GenerateUseCase) Execute(ctx context.Context, params GenerateParams) ([]llm.GeneratedQuestion, error) {
	// Get document
	document, err := uc.documentRepo.FindByID(ctx, params.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Verify ownership
	if document.UserID != params.UserID {
		return nil, fmt.Errorf("unauthorized access to document")
	}

	// Check if document is parsed
	if !document.IsParsed() {
		return nil, fmt.Errorf("document not parsed yet")
	}

	// Default to perplexity if no provider specified
	provider := params.LLMProvider
	if provider == "" {
		provider = "perplexity"
	}

	// Create LLM strategy
	strategy, err := uc.llmFactory.CreateStrategy(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM strategy: %w", err)
	}

	// Create LLM context and generate questions
	llmContext := llm.NewLLMContext(strategy)
	questions, err := llmContext.GenerateQuestions(ctx, llm.GenerationParams{
		Text:         document.ParsedText,
		NumQuestions: params.NumQuestions,
		Difficulty:   params.Difficulty,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	return questions, nil
}
