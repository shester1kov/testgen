package postgres

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations_InvalidConnection(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=0 user=invalid password=invalid dbname=invalid sslmode=disable")
	require.NoError(t, err)

	err = runMigrations(db, "invalid", nil) // Pass nil logger for test
	assert.Error(t, err)
}

func TestNewDatabase_InvalidConfig(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     "0",
		User:     "invalid",
		Password: "invalid",
		DBName:   "invalid",
		SSLMode:  "disable",
		Logger:   nil, // Pass nil logger for test
	}

	db, err := NewDatabase(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}
