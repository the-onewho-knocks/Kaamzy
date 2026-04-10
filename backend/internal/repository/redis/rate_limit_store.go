package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRule struct {
	MaxRequests int
	Window      time.Duration
}

type RateLimitStore struct {
	client *redis.Client
	rules  map[string]RateLimitRule
}

func NewRateLimitStore(client *redis.Client) *RateLimitStore {
	return &RateLimitStore{
		client: client,
		rules:  make(map[string]RateLimitRule),
	}
}

func (s *RateLimitStore) AddRule(key string, maxRequests int, window time.Duration) {
	s.rules[key] = RateLimitRule{
		MaxRequests: maxRequests,
		Window:      window,
	}
}

func (s *RateLimitStore) Allow(ctx context.Context, key string) (bool, int, error) {
	rule, exists := s.rules[key]
	if !exists {
		return true, 0, nil
	}

	now := time.Now().UnixMilli()
	windowStart := now - rule.Window.Milliseconds()

	pipe := s.client.Pipeline()

	pipe.ZRemRangeByScore(ctx, s.key(key), "0", strconv.FormatInt(windowStart, 10))

	countCmd := pipe.ZCard(ctx, s.key(key))

	pipe.ZAdd(ctx, s.key(key), redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	pipe.Expire(ctx, s.key(key), rule.Window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := int(countCmd.Val())
	remaining := rule.MaxRequests - count - 1
	if remaining < 0 {
		remaining = 0
	}

	return count < rule.MaxRequests, remaining, nil
}

func (s *RateLimitStore) GetLimit(ctx context.Context, key string) (int, error) {
	rule, exists := s.rules[key]
	if !exists {
		return 0, nil
	}
	return rule.MaxRequests, nil
}

func (s *RateLimitStore) Reset(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.key(key)).Err()
}

func (s *RateLimitStore) GetUsage(ctx context.Context, key string) (int, error) {
	rule, exists := s.rules[key]
	if !exists {
		return 0, nil
	}

	now := time.Now().UnixMilli()
	windowStart := now - rule.Window.Milliseconds()

	count, err := s.client.ZCount(ctx, s.key(key), strconv.FormatInt(windowStart, 10), "+").Result()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *RateLimitStore) key(identifier string) string {
	return fmt.Sprintf("ratelimit:%s", identifier)
}
