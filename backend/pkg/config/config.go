package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	File     FileConfig
	LLM      LLMConfig
	Moodle   MoodleConfig
	Logger   LoggerConfig
	Admin    AdminConfig
	MinIO    MinIOConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port          string
	Environment   string
	EnableMetrics bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration string
}

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

// FileConfig holds file upload configuration
type FileConfig struct {
	MaxFileSize int64
	UploadDir   string
}

// LLMConfig holds LLM API configuration
type LLMConfig struct {
	Provider         string
	PerplexityAPIKey string
	OpenAIAPIKey     string
	YandexAPIKey     string
	YandexFolderID   string
	YandexModel      string
}

// MoodleConfig holds Moodle integration configuration
type MoodleConfig struct {
	URL         string
	Token       string
	ImportToken string
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string
	Format string
}

// AdminConfig holds default admin user configuration
type AdminConfig struct {
	Email    string
	Password string
	FullName string
}

// MinIO настройки для подключения к MinIO
type MinIOConfig struct {
	// Endpoint - адрес MinIO сервера
	// Формат: "host:port" (например: "localhost:9000" или "minio.example.com:9000")
	// Не включает протокол (http:// или https://)
	Endpoint string

	// AccessKey - ключ доступа (аналог логина)
	// В MinIO по умолчанию: "minioadmin"
	// В production должен быть уникальным и сложным
	AccessKey string

	// SecretKey - секретный ключ (аналог пароля)
	// В MinIO по умолчанию: "minioadmin"
	// В production должен быть длинным случайным значением
	SecretKey string

	// BucketName - имя bucket для хранения файлов
	// Должен быть создан заранее или создастся автоматически в NewMinIOStorage
	BucketName string

	// UseSSL - использовать ли HTTPS/TLS для подключения
	// false - для локальной разработки (http://)
	// true - для production (https://)
	UseSSL bool

	// Region - регион S3 (опционально)
	// Нужен для AWS S3, для MinIO обычно не используется
	// Можно оставить пустым для MinIO
	Region string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {

	// Загружаем переменные из .env файла
	// godotenv.Load() читает файл .env и добавляет переменные в os.Getenv()
	// Если файл не найден - не критично, можно использовать системные env переменные

	if err := godotenv.Load(); err != nil {
		// Not a fatal error - .env file is optional in production
		// where system environment variables are typically used
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:          getEnv("PORT", "8080"),
			Environment:   getEnv("ENV", "development"),
			EnableMetrics: getEnvBool("ENABLE_METRICS", true),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "testgen_user"),
			Password: getEnv("DB_PASSWORD", "testgen_password"),
			DBName:   getEnv("DB_NAME", "testgen_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		Cookie: CookieConfig{
			Name:     getEnv("COOKIE_NAME", "testgen_token"),
			Domain:   getEnv("COOKIE_DOMAIN", ""),
			Path:     getEnv("COOKIE_PATH", "/"),
			Secure:   getEnvBool("COOKIE_SECURE", false),
			HTTPOnly: getEnvBool("COOKIE_HTTP_ONLY", true),
			SameSite: getEnv("COOKIE_SAME_SITE", "Lax"),
		},
		File: FileConfig{
			MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 52428800), // 50MB
			UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		},
		LLM: LLMConfig{
			Provider:         getEnv("LLM_PROVIDER", "yandexgpt"),
			PerplexityAPIKey: getEnv("PERPLEXITY_API_KEY", ""),
			OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
			YandexAPIKey:     getEnv("YANDEX_GPT_API_KEY", ""),
			YandexFolderID:   getEnv("YANDEX_GPT_FOLDER_ID", ""),
			YandexModel:      getEnv("YANDEX_GPT_MODEL", "yandexgpt-lite"),
		},
		Moodle: MoodleConfig{
			URL:         getEnv("MOODLE_URL", ""),
			Token:       getEnv("MOODLE_TOKEN", ""),
			ImportToken: getEnv("MOODLE_IMPORT_TOKEN", "testgen_secret_token_change_in_production"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "console"),
		},
		Admin: AdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@testgen.local"),
			Password: getEnv("ADMIN_PASSWORD", "admin123"),
			FullName: getEnv("ADMIN_FULL_NAME", "System Administrator"),
		},
		MinIO: MinIOConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
			BucketName: getEnv("MINIO_BUCKET_NAME", "testgen-documents"),
			UseSSL:     getEnvBool("MINIO_USE_SSL", false),
			Region:     getEnv("MINIO_REGION", ""),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// Validate validates the configuration with environment-aware checks
func (c *Config) Validate() error {
	isDev := c.Server.Environment == "development"

	// JWT secret validation - always required to be secure
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("JWT_SECRET must be set and not default 'your-secret-key'")
	}

	// In production, also reject the "change-in-production" default
	if !isDev && c.JWT.Secret == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed from default in production")
	}

	// Database validation
	if c.Database.Host == "" || c.Database.DBName == "" {
		return fmt.Errorf("database configuration is incomplete (missing host or dbname)")
	}

	// MinIO validation
	if c.MinIO.Endpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT must be set")
	}

	if c.MinIO.AccessKey == "" || c.MinIO.SecretKey == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY and MINIO_SECRET_KEY must be set")
	}

	if c.MinIO.BucketName == "" {
		return fmt.Errorf("MINIO_BUCKET_NAME must be set")
	}

	// Production-specific strict checks
	if !isDev {
		if c.MinIO.AccessKey == "minioadmin" || c.MinIO.SecretKey == "minioadmin" {
			return fmt.Errorf("MINIO credentials must not be 'minioadmin' in production")
		}

		if !c.MinIO.UseSSL {
			fmt.Println("WARNING: MinIO SSL is disabled in production!")
		}
	}

	return nil
}
