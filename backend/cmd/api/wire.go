// go:build wireinject
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/document"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/cache"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/exporter"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/storage"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"github.com/shester1kov/testgen-backend/pkg/logger"
	"github.com/shester1kov/testgen-backend/pkg/utils"
	"gorm.io/gorm"
)

// ApplicationContainer holds all application dependencies
type ApplicationContainer struct {
	// handlers
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	DocumentHandler *handler.DocumentHandler
	TestHandler     *handler.TestHandler
	MoodleHandler   *handler.MoodleHandler
	StatsHandler    *handler.StatsHandler

	// infra
	JWTManager *utils.JWTManager

	// repositories (для seeder)
	UserRepo    repository.UserRepository
	RoleRepo    repository.RoleRepository
	RedisClient *cache.RedisClient
}

// InitializeApplication sets up all dependencies using Wire
func InitializeApplication(cfg *config.Config, db *gorm.DB, log *logger.Logger) (*ApplicationContainer, error) {
	wire.Build(
		// Redis Cache
		provideRedisClient,

		// Storage
		provideMinIOStorage,

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

		// Exporter Factory (for multiple LMS formats)
		exporter.NewExporterFactory,

		// Moodle components
		provideMoodleClient,

		// Document Use Cases
		provideUploadUseCase,
		provideDeleteUseCase,
		provideParseUseCase,
		provideListUseCase,
		provideGetUseCase,

		// Handlers
		provideAuthHandler,
		handler.NewUserHandler,
		handler.NewDocumentHandler,
		handler.NewTestHandler,
		handler.NewMoodleHandler,
		handler.NewStatsHandler,

		// Wire the ApplicationContainer
		wire.Struct(new(ApplicationContainer), "*"),
	)

	return &ApplicationContainer{}, nil
}

// Provider functions for Wire

func provideRedisClient(cfg *config.Config, log *logger.Logger) (*cache.RedisClient, error) {
	return cache.NewRedisClient(cfg.Redis, log.Logger)
}

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

func provideMinIOStorage(cfg *config.Config, log *logger.Logger) (storage.Storage, error) {
	return storage.NewMinIOStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.BucketName,
		cfg.MinIO.UseSSL,
		log.Logger,
	)
}

func provideUploadUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	cfg *config.Config,
	log *logger.Logger,
) *document.UploadUseCase {
	return document.NewUploadUseCase(
		documentRepo,
		storage,
		cfg.File.MaxFileSize,
		log.Logger,
	)
}

func provideDeleteUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	log *logger.Logger,
) *document.DeleteUseCase {
	return document.NewDeleteUseCase(
		documentRepo,
		storage,
		log.Logger,
	)
}

func provideParseUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	parserFactory *parser.DocumentParserFactory,
	log *logger.Logger,
) *document.ParseUseCase {
	return document.NewParseUseCase(
		documentRepo,
		storage,
		parserFactory,
		log.Logger,
	)
}

func provideListUseCase(
	documentRepo repository.DocumentRepository,
	userRepo repository.UserRepository,
	log *logger.Logger,
) *document.ListUseCase {
	return document.NewListUseCase(
		documentRepo,
		userRepo,
		log.Logger,
	)
}

func provideGetUseCase(
	documentRepo repository.DocumentRepository,
	log *logger.Logger,
) *document.GetUseCase {
	return document.NewGetUseCase(
		documentRepo,
		log.Logger,
	)
}
