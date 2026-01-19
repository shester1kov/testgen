package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shester1kov/testgen-backend/pkg/config"
	"go.uber.org/zap"
)

// RedisClient wraps redis.Client with additional helper methods
type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(cfg config.RedisConfig, logger *zap.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis client connected successfully",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

// Get retrieves a value from Redis and unmarshals it into dest
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		r.logger.Error("Redis GET error", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("redis get error: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		r.logger.Error("Failed to unmarshal cached value", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("unmarshal error: %w", err)
	}

	r.logger.Debug("Cache hit", zap.String("key", key))
	return nil
}

// Set stores a value in Redis with the specified TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("Failed to marshal value", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		r.logger.Error("Redis SET error", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("redis set error: %w", err)
	}

	r.logger.Debug("Cache set", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		r.logger.Error("Redis DEL error", zap.Strings("keys", keys), zap.Error(err))
		return fmt.Errorf("redis delete error: %w", err)
	}

	r.logger.Debug("Cache deleted", zap.Strings("keys", keys))
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Redis EXISTS error", zap.String("key", key), zap.Error(err))
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return count > 0, nil
}

// DeleteByPattern deletes all keys matching a pattern (use with caution!)
func (r *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var deletedCount int

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.logger.Error("Redis SCAN error", zap.String("pattern", pattern), zap.Error(err))
			return fmt.Errorf("redis scan error: %w", err)
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				r.logger.Error("Redis DEL error", zap.Strings("keys", keys), zap.Error(err))
				return fmt.Errorf("redis delete error: %w", err)
			}
			deletedCount += len(keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	r.logger.Debug("Cache deleted by pattern",
		zap.String("pattern", pattern),
		zap.Int("count", deletedCount),
	)
	return nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if err := r.client.Close(); err != nil {
		r.logger.Error("Failed to close Redis connection", zap.Error(err))
		return err
	}
	r.logger.Info("Redis connection closed")
	return nil
}

// Ping checks if Redis is reachable
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
