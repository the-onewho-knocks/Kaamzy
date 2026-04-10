package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
		PoolSize: 10,
	}
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

func RedisHealthCheck(client *redis.Client, ctx context.Context) error {
	return client.Ping(ctx).Err()
}
