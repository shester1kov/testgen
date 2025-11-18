//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yourusername/testgen-backend/internal/infrastructure/llm"
	"github.com/yourusername/testgen-backend/internal/infrastructure/moodle"
	"github.com/yourusername/testgen-backend/internal/infrastructure/parser"
	"github.com/yourusername/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/yourusername/testgen-backend/internal/interfaces/http/handler"
	"github.com/yourusername/testgen-backend/pkg/config"
	"github.com/yourusername/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// ApplicationContainer holds all application dependencies
type ApplicationContainer struct {
	AuthHandler     *handler.AuthHandler
	DocumentHandler *handler.DocumentHandler
	TestHandler     *handler.TestHandler
	MoodleHandler   *handler.MoodleHandler
	JWTManager      *utils.JWTManager
}

// InitializeApplication sets up all dependencies using Wire
func InitializeApplication(cfg *config.Config, db *gorm.DB) (*ApplicationContainer, error) {
	wire.Build(
		// Repositories
		postgres.NewUserRepository,
		postgres.NewDocumentRepository,
		postgres.NewTestRepository,
		postgres.NewQuestionRepository,
		postgres.NewAnswerRepository,

		// JWT Manager
		provideJWTManager,

		// Document Parser Factory
		parser.NewDocumentParserFactory,

		// LLM Factory
		provideLLMFactory,

		// Moodle components
		moodle.NewMoodleXMLExporter,
		provideMoodleClient,

		// Handlers
		handler.NewAuthHandler,
		handler.NewDocumentHandler,
		handler.NewTestHandler,
		handler.NewMoodleHandler,

		// File config providers
		provideUploadDir,
		provideMaxFileSize,

		// Wire the ApplicationContainer
		wire.Struct(new(ApplicationContainer), "*"),
	)

	return &ApplicationContainer{}, nil
}

// Provider functions for Wire

func provideJWTManager(cfg *config.Config) (*utils.JWTManager, error) {
	return utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
}

func provideLLMFactory(cfg *config.Config) *llm.LLMFactory {
	return llm.NewLLMFactory(
		cfg.LLM.PerplexityAPIKey,
		cfg.LLM.OpenAIAPIKey,
		cfg.LLM.YandexAPIKey,
	)
}

func provideMoodleClient(cfg *config.Config) *moodle.Client {
	if cfg.Moodle.URL != "" && cfg.Moodle.Token != "" {
		return moodle.NewClient(cfg.Moodle.URL, cfg.Moodle.Token)
	}
	return nil
}

func provideUploadDir(cfg *config.Config) string {
	return cfg.File.UploadDir
}

func provideMaxFileSize(cfg *config.Config) int64 {
	return cfg.File.MaxFileSize
}
