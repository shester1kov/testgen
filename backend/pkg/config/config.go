package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	File     FileConfig
	LLM      LLMConfig
	Moodle   MoodleConfig
	Logger   LoggerConfig
	Admin    AdminConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port           string
	Environment    string
	EnableMetrics  bool
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
}

// MoodleConfig holds Moodle integration configuration
type MoodleConfig struct {
	URL   string
	Token string
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

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
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
			Provider:         getEnv("LLM_PROVIDER", "perplexity"),
			PerplexityAPIKey: getEnv("PERPLEXITY_API_KEY", ""),
			OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
			YandexAPIKey:     getEnv("YANDEX_GPT_API_KEY", ""),
		},
		Moodle: MoodleConfig{
			URL:   getEnv("MOODLE_URL", ""),
			Token: getEnv("MOODLE_TOKEN", ""),
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
	}
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
