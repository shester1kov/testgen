// go:build wireinject
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/auth"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/document"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/stats"
	testuc "github.com/shester1kov/testgen-backend/internal/application/usecase/test"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/user"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/cache"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/exporter"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/persistence/postgres"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/storage"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/handler"
	"github.com/shester1kov/testgen-backend/internal/interfaces/http/middleware"
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
	StatsHandler    *handler.StatsHandler

	// infra
	JWTManager *utils.JWTManager

	// repositories (for seeder)
	UserRepo    repository.UserRepository
	RoleRepo    repository.RoleRepository
	RedisClient *cache.RedisClient

	// activity logging
	ActivityLogger *middleware.ActivityLogger
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
		postgres.NewActivityLogRepository,
		postgres.NewAcademicProfileRepository,
		postgres.NewTestGenerationConfigRepository,
		middleware.NewActivityLogger,

		// JWT Manager
		provideJWTManager,

		// Document Parser Factory
		parser.NewDocumentParserFactory,

		// LLM Factory
		provideLLMFactory,

		// Exporter Factory (for multiple LMS formats)
		exporter.NewExporterFactory,

		// Auth Use Cases
		provideRegisterUseCase,
		provideLoginUseCase,
		provideGetMeUseCase,

		// Document Use Cases
		provideUploadUseCase,
		provideDeleteUseCase,
		provideParseUseCase,
		provideListDocumentUseCase,
		provideGetDocumentUseCase,

		// User Use Cases
		provideListUserUseCase,
		provideUpdateRoleUseCase,

		// Test Use Cases
		provideCreateTestUseCase,
		provideGenerateTestUseCase,
		provideListTestUseCase,
		provideGetTestUseCase,
		provideUpdateTestUseCase,
		provideDeleteTestUseCase,
		provideUpdateQuestionUseCase,
		provideExportTestUseCase,

		// Stats Use Cases
		provideDashboardUseCase,

		// Handlers
		provideAuthHandler,
		handler.NewUserHandler,
		handler.NewDocumentHandler,
		handler.NewTestHandler,
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

func provideLLMFactory(cfg *config.Config, log *logger.Logger) *llm.LLMFactory {
	return llm.NewLLMFactory(
		cfg.LLM.PerplexityAPIKey,
		cfg.LLM.OpenAIAPIKey,
		cfg.LLM.YandexAPIKey,
		cfg.LLM.YandexFolderID,
		cfg.LLM.YandexModel,
		log.Logger,
	)
}

func provideAuthHandler(
	cfg *config.Config,
	registerUseCase *auth.RegisterUseCase,
	loginUseCase *auth.LoginUseCase,
	getMeUseCase *auth.GetMeUseCase,
) *handler.AuthHandler {
	return handler.NewAuthHandler(
		registerUseCase,
		loginUseCase,
		getMeUseCase,
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

// Auth Use Case providers

func provideRegisterUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtManager *utils.JWTManager,
) *auth.RegisterUseCase {
	return auth.NewRegisterUseCase(userRepo, roleRepo, jwtManager)
}

func provideLoginUseCase(
	userRepo repository.UserRepository,
	jwtManager *utils.JWTManager,
) *auth.LoginUseCase {
	return auth.NewLoginUseCase(userRepo, jwtManager)
}

func provideGetMeUseCase(
	userRepo repository.UserRepository,
) *auth.GetMeUseCase {
	return auth.NewGetMeUseCase(userRepo)
}

// Document Use Case providers

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
	userRepo repository.UserRepository,
	log *logger.Logger,
) *document.DeleteUseCase {
	return document.NewDeleteUseCase(
		documentRepo,
		storage,
		userRepo,
		log.Logger,
	)
}

func provideParseUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	parserFactory *parser.DocumentParserFactory,
	userRepo repository.UserRepository,
	log *logger.Logger,
) *document.ParseUseCase {
	return document.NewParseUseCase(
		documentRepo,
		storage,
		parserFactory,
		userRepo,
		log.Logger,
	)
}

func provideListDocumentUseCase(
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

func provideGetDocumentUseCase(
	documentRepo repository.DocumentRepository,
	userRepo repository.UserRepository,
	log *logger.Logger,
) *document.GetUseCase {
	return document.NewGetUseCase(
		documentRepo,
		userRepo,
		log.Logger,
	)
}

// User Use Case providers

func provideListUserUseCase(
	userRepo repository.UserRepository,
) *user.ListUseCase {
	return user.NewListUseCase(userRepo)
}

func provideUpdateRoleUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) *user.UpdateRoleUseCase {
	return user.NewUpdateRoleUseCase(userRepo, roleRepo)
}

// Test Use Case providers

func provideCreateTestUseCase(
	testRepo repository.TestRepository,
) *testuc.CreateUseCase {
	return testuc.NewCreateUseCase(testRepo)
}

func provideGenerateTestUseCase(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
	academicProfileRepo repository.AcademicProfileRepository,
	genConfigRepo repository.TestGenerationConfigRepository,
	llmFactory *llm.LLMFactory,
	log *logger.Logger,
) *testuc.GenerateUseCase {
	return testuc.NewGenerateUseCase(testRepo, documentRepo, questionRepo, answerRepo, userRepo, academicProfileRepo, genConfigRepo, llmFactory, log.Logger)
}

func provideListTestUseCase(
	testRepo repository.TestRepository,
	userRepo repository.UserRepository,
) *testuc.ListUseCase {
	return testuc.NewListUseCase(testRepo, userRepo)
}

func provideGetTestUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
) *testuc.GetUseCase {
	return testuc.NewGetUseCase(testRepo, questionRepo, answerRepo, userRepo)
}

func provideUpdateTestUseCase(
	testRepo repository.TestRepository,
	userRepo repository.UserRepository,
) *testuc.UpdateUseCase {
	return testuc.NewUpdateUseCase(testRepo, userRepo)
}

func provideDeleteTestUseCase(
	testRepo repository.TestRepository,
	userRepo repository.UserRepository,
) *testuc.DeleteUseCase {
	return testuc.NewDeleteUseCase(testRepo, userRepo)
}

func provideUpdateQuestionUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
) *testuc.UpdateQuestionUseCase {
	return testuc.NewUpdateQuestionUseCase(testRepo, questionRepo, answerRepo, userRepo)
}

func provideExportTestUseCase(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	exporterFactory *exporter.ExporterFactory,
	userRepo repository.UserRepository,
) *testuc.ExportUseCase {
	return testuc.NewExportUseCase(testRepo, questionRepo, answerRepo, exporterFactory, userRepo)
}

// Stats Use Case providers

func provideDashboardUseCase(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	userRepo repository.UserRepository,
) *stats.DashboardUseCase {
	return stats.NewDashboardUseCase(testRepo, documentRepo, questionRepo, userRepo)
}
