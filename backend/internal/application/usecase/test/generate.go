package test

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shester1kov/testgen-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/pkg/security"
	"go.uber.org/zap"
)

// GenerateUseCase handles test generation using LLM
type GenerateUseCase struct {
	testRepo            repository.TestRepository
	documentRepo        repository.DocumentRepository
	questionRepo        repository.QuestionRepository
	answerRepo          repository.AnswerRepository
	userRepo            repository.UserRepository
	academicProfileRepo repository.AcademicProfileRepository
	genConfigRepo       repository.TestGenerationConfigRepository
	llmFactory          *llm.LLMFactory
	logger              *zap.Logger
}

// NewGenerateUseCase creates a new generate use case.
// The logger parameter is optional; if omitted, a no-op logger is used.
func NewGenerateUseCase(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
	academicProfileRepo repository.AcademicProfileRepository,
	genConfigRepo repository.TestGenerationConfigRepository,
	llmFactory *llm.LLMFactory,
	logger ...*zap.Logger,
) *GenerateUseCase {
	var l *zap.Logger
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	} else {
		l = zap.NewNop()
	}

	return &GenerateUseCase{
		testRepo:            testRepo,
		documentRepo:        documentRepo,
		questionRepo:        questionRepo,
		answerRepo:          answerRepo,
		userRepo:            userRepo,
		academicProfileRepo: academicProfileRepo,
		genConfigRepo:       genConfigRepo,
		llmFactory:          llmFactory,
		logger:              l,
	}
}

// GenerateParams contains generation parameters
type GenerateParams struct {
	UserID          uuid.UUID
	DocumentID      uuid.UUID
	Title           string
	NumQuestions    int
	Difficulty      string
	LLMProvider     string
	AcademicProfile llm.AcademicProfile
}

// Execute executes the generate use case — generates questions via LLM and saves the full test
func (uc *GenerateUseCase) Execute(ctx context.Context, params GenerateParams) (*dto.TestResponse, error) {
	// Get document
	document, err := uc.documentRepo.FindByID(ctx, params.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDocumentNotFound, err)
	}

	// Check access: admin sees all, others only own documents
	user, err := uc.userRepo.FindByID(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsAdmin() && document.UserID != params.UserID {
		return nil, fmt.Errorf("%w: document %s", domain.ErrAccessDenied, params.DocumentID)
	}

	// Check if document is parsed
	if !document.IsParsed() {
		return nil, domain.ErrNotParsedYet
	}

	// Default provider
	provider := params.LLMProvider
	if provider == "" {
		provider = "perplexity"
	}

	// Create LLM strategy
	strategy, err := uc.llmFactory.CreateStrategy(provider)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidLLMProvider, err)
	}

	// Compute multiple_choice count based on difficulty
	// easy: 0%, medium: 30% (ceil), hard: 60% (ceil)
	multipleChoiceCount := 0
	switch params.Difficulty {
	case "medium":
		multipleChoiceCount = int(math.Ceil(float64(params.NumQuestions) * 0.3))
	case "hard":
		multipleChoiceCount = int(math.Ceil(float64(params.NumQuestions) * 0.6))
	}

	// Default profile
	profile := params.AcademicProfile
	if profile == "" {
		profile = llm.ProfileUniversal
	}

	// Load profile data from DB to get per-profile temperature and prompt content.
	// Falls back to hardcoded prompt builder data when repo is unavailable (e.g. in tests).
	var profileEntity *entity.AcademicProfile
	if uc.academicProfileRepo != nil {
		profileEntity, err = uc.academicProfileRepo.FindByCode(ctx, string(profile))
		if err != nil {
			uc.logger.Warn("academic profile not found in DB, falling back to hardcoded data",
				zap.String("profile", string(profile)),
				zap.Error(err),
			)
			profileEntity = nil
		}
	}

	genParams := llm.GenerationParams{
		Text:                document.ParsedText,
		NumQuestions:        params.NumQuestions,
		Difficulty:          params.Difficulty,
		MultipleChoiceCount: multipleChoiceCount,
		Profile:             profile,
	}

	if profileEntity != nil {
		genParams.Temperature = profileEntity.Temperature
		genParams.ProfileData = toProfileData(profileEntity, profile)
	}

	// Generate questions via LLM
	llmContext := llm.NewLLMContext(strategy)
	questions, err := llmContext.GenerateQuestions(ctx, genParams)
	if err != nil {
		uc.logger.Error("LLM generation failed",
			zap.String("provider", params.LLMProvider),
			zap.String("document_id", params.DocumentID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToGenerate, err)
	}

	// Sanitize title
	sanitizedTitle := security.SanitizeInput(params.Title)

	// Create test entity
	test := &entity.Test{
		ID:             uuid.New(),
		UserID:         params.UserID,
		DocumentID:     &params.DocumentID,
		Title:          sanitizedTitle,
		TotalQuestions: len(questions),
		Status:         entity.TestStatusDraft,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.testRepo.Create(ctx, test); err != nil {
		return nil, fmt.Errorf("failed to save test: %w", err)
	}

	// Save generation config when profile was loaded from DB
	if profileEntity != nil {
		genConfig := &entity.TestGenerationConfig{
			TestID:            test.ID,
			AcademicProfileID: profileEntity.ID,
			Difficulty:        params.Difficulty,
			NumQuestions:      params.NumQuestions,
			LLMProvider:       provider,
			Language:          "ru",
			CreatedAt:         time.Now(),
		}
		if err := uc.genConfigRepo.Create(ctx, genConfig); err != nil {
			uc.logger.Warn("failed to save test generation config",
				zap.String("test_id", test.ID.String()),
				zap.Error(err),
			)
		}
	}

	// Save questions and answers
	for i, q := range questions {
		sanitizedQuestionText := security.SanitizeMultiline(q.QuestionText)

		question := &entity.Question{
			ID:           uuid.New(),
			TestID:       test.ID,
			QuestionText: sanitizedQuestionText,
			QuestionType: entity.QuestionType(q.QuestionType),
			Difficulty:   entity.Difficulty(q.Difficulty),
			Points:       1.0,
			OrderNum:     i + 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := uc.questionRepo.Create(ctx, question); err != nil {
			return nil, fmt.Errorf("failed to save question: %w", err)
		}

		for j, a := range q.Answers {
			sanitizedAnswerText := security.SanitizeInput(a.Text)

			answer := &entity.Answer{
				ID:         uuid.New(),
				QuestionID: question.ID,
				AnswerText: sanitizedAnswerText,
				IsCorrect:  a.IsCorrect,
				OrderNum:   j + 1,
				CreatedAt:  time.Now(),
			}

			if err := uc.answerRepo.Create(ctx, answer); err != nil {
				return nil, fmt.Errorf("failed to save answer: %w", err)
			}
		}
	}

	return &dto.TestResponse{
		ID:             test.ID.String(),
		UserID:         test.UserID.String(),
		Title:          test.Title,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   false,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
	}, nil
}

// toProfileData converts a DB entity to the llm.ProfileData used by the prompt builder.
func toProfileData(p *entity.AcademicProfile, code llm.AcademicProfile) *llm.ProfileData {
	formulations := make([]llm.FormulationExample, len(p.Formulations))
	for i, f := range p.Formulations {
		formulations[i] = llm.FormulationExample{
			Bad:  f.BadExample,
			Good: f.GoodExample,
		}
	}

	examples := make([]llm.QuestionExample, len(p.QuestionExamples))
	for i, qe := range p.QuestionExamples {
		answers := make([]llm.ExampleAnswer, len(qe.Answers))
		for j, a := range qe.Answers {
			answers[j] = llm.ExampleAnswer{
				Text:      a.Text,
				IsCorrect: a.IsCorrect,
			}
		}
		examples[i] = llm.QuestionExample{
			QuestionType: llm.QuestionType(qe.QuestionType),
			QuestionText: qe.QuestionText,
			Answers:      answers,
			Explanation:  qe.Explanation,
		}
	}

	return &llm.ProfileData{
		Code:             code,
		Temperature:      p.Temperature,
		Instruction:      p.ProfileInstruction,
		Formulations:     formulations,
		QuestionExamples: examples,
	}
}
