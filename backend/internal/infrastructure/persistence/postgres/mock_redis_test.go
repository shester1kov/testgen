package postgres

import (
	"github.com/alicebob/miniredis/v2"
	"go.uber.org/zap"

	"github.com/shester1kov/testgen-backend/internal/infrastructure/cache"
	"github.com/shester1kov/testgen-backend/pkg/config"
)

// newMockRedisClient creates a mock Redis client using miniredis (in-memory Redis)
// This allows tests to run without a real Redis server
func newMockRedisClient() *cache.RedisClient {
	// Create an in-memory Redis server
	mr, err := miniredis.Run()
	if err != nil {
		panic("failed to start miniredis: " + err.Error())
	}

	// Create Redis client config pointing to miniredis
	cfg := config.RedisConfig{
		Host:     mr.Host(),
		Port:     mr.Port(),
		Password: "",
		DB:       0,
		PoolSize: 1,
	}

	// Create a noop logger
	logger := zap.NewNop()

	// Create real Redis client connected to miniredis
	client, err := cache.NewRedisClient(cfg, logger)
	if err != nil {
		panic("failed to create Redis client: " + err.Error())
	}

	return client
}
