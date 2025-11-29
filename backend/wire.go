//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// ApplicationContainer holds all application dependencies
type ApplicationContainer struct {
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	DocumentHandler *handler.DocumentHandler
	TestHandler     *handler.TestHandler
	MoodleHandler   *handler.MoodleHandler
	StatsHandler    *handler.StatsHandler
	JWTManager      *utils.JWTManager
}

// InitializeApplication sets up all dependencies using Wire
func InitializeApplication(cfg *config.Config, db *gorm.DB) (*ApplicationContainer, error) {
	wire.Build(
		// Repositories
		postgres.NewUserRepository,
		postgres.NewRoleRepository,
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
		provideAuthHandler,
		handler.NewUserHandler,
		handler.NewDocumentHandler,
		handler.NewTestHandler,
		handler.NewMoodleHandler,
		handler.NewStatsHandler,

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
		cfg.LLM.YandexFolderID,
		cfg.LLM.YandexModel,
	)
}

func provideMoodleClient(cfg *config.Config) *moodle.Client {
	if cfg.Moodle.URL != "" && cfg.Moodle.Token != "" {
		return moodle.NewClient(cfg.Moodle.URL, cfg.Moodle.Token, cfg.Moodle.ImportToken)
	}
	return nil
}

func provideUploadDir(cfg *config.Config) string {
	return cfg.File.UploadDir
}

func provideMaxFileSize(cfg *config.Config) int64 {
	return cfg.File.MaxFileSize
}

func provideAuthHandler(
	cfg *config.Config,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtManager *utils.JWTManager,
) *handler.AuthHandler {
	return handler.NewAuthHandler(
		userRepo,
		roleRepo,
		jwtManager,
		cfg.Cookie.Name,
		cfg.Cookie.Domain,
		cfg.Cookie.Path,
		cfg.Cookie.SameSite,
		cfg.JWT.Expiration,
		cfg.Cookie.Secure,
		cfg.Cookie.HTTPOnly,
	)
}
