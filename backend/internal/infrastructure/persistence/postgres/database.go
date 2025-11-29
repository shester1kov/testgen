package postgres

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	customlogger "github.com/shester1kov/testgen-backend/pkg/logger"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Logger   *customlogger.Logger // Add logger to config
}

// NewDatabase creates a new database connection
func NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// First, connect with standard library for migrations
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Run migrations
	if err := runMigrations(sqlDB, config.DBName, config.Logger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create GORM logger from Zap logger
	gormLogger := customlogger.NewGormLogger(config.Logger)

	// Now open with GORM
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	gormDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	gormDB.SetMaxIdleConns(10)
	gormDB.SetMaxOpenConns(100)

	// Database connection established - logged by caller
	return db, nil
}

// runMigrations automatically applies database migrations
func runMigrations(db *sql.DB, dbName string, logger *customlogger.Logger) error {
	if logger != nil {
		logger.Info("Running database migrations...")
	}

	// Create postgres driver instance
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations directory path
	migrationsPath := filepath.Join("internal", "infrastructure", "persistence", "migrations")
	// Convert Windows backslashes to forward slashes for file:// URL
	migrationsPath = filepath.ToSlash(migrationsPath)
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	if logger != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new migrations to apply")
		} else {
			logger.Info("Database migrations completed successfully")
		}
	}

	return nil
}
