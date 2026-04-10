package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CachedJobFeed struct {
	Jobs       []JobFeedItem `json:"jobs"`
	Filters    string        `json:"filters"`
	CachedAt   time.Time     `json:"cached_at"`
	TotalCount int           `json:"total_count"`
}

type JobFeedItem struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	BusinessName  string    `json:"business_name"`
	City          string    `json:"city"`
	PaymentAmount float64   `json:"payment_amount"`
	PaymentType   string    `json:"payment_type"`
	JobType       string    `json:"job_type"`
	CreatedAt     time.Time `json:"created_at"`
}

type JobFeedCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewJobFeedCache(client *redis.Client, ttl time.Duration) *JobFeedCache {
	return &JobFeedCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *JobFeedCache) CacheFeed(ctx context.Context, userID uuid.UUID, filters string, items []JobFeedItem, totalCount int) error {
	feed := CachedJobFeed{
		Jobs:       items,
		Filters:    filters,
		CachedAt:   time.Now(),
		TotalCount: totalCount,
	}

	b, err := json.Marshal(feed)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.key(userID, filters), b, c.ttl).Err()
}

func (c *JobFeedCache) GetFeed(ctx context.Context, userID uuid.UUID, filters string) (*CachedJobFeed, error) {
	val, err := c.client.Get(ctx, c.key(userID, filters)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var feed CachedJobFeed
	if err := json.Unmarshal(val, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

func (c *JobFeedCache) InvalidateFeed(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("jobfeed:%s:*", userID.String())
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

func (c *JobFeedCache) InvalidateAll(ctx context.Context) error {
	pattern := "jobfeed:*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

func (c *JobFeedCache) IsStale(cachedAt time.Time) bool {
	return time.Since(cachedAt) > c.ttl
}

func (c *JobFeedCache) key(userID uuid.UUID, filters string) string {
	return fmt.Sprintf("jobfeed:%s:%s", userID.String(), filters)
}
