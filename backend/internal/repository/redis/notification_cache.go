package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CachedNotification struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
}

type NotificationCache struct {
	client    *redis.Client
	ttl       time.Duration
	maxUnread int
}

func NewNotificationCache(client *redis.Client, ttl time.Duration) *NotificationCache {
	return &NotificationCache{
		client:    client,
		ttl:       ttl,
		maxUnread: 50,
	}
}

func (c *NotificationCache) CacheNotification(ctx context.Context, notification *CachedNotification) error {
	b, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	key := c.unreadKey(notification.UserID)
	pipe := c.client.Pipeline()

	pipe.LPush(ctx, key, b)
	pipe.LTrim(ctx, key, 0, int64(c.maxUnread-1))
	pipe.Expire(ctx, key, c.ttl)

	_, err = pipe.Exec(ctx)
	return err
}

func (c *NotificationCache) GetUnread(ctx context.Context, userID uuid.UUID, limit int) ([]CachedNotification, error) {
	key := c.unreadKey(userID)

	if limit <= 0 {
		limit = c.maxUnread
	}

	vals, err := c.client.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	notifications := make([]CachedNotification, 0, len(vals))
	for _, v := range vals {
		var n CachedNotification
		if err := json.Unmarshal([]byte(v), &n); err != nil {
			continue
		}
		if !n.IsRead {
			notifications = append(notifications, n)
		}
	}

	return notifications, nil
}

func (c *NotificationCache) GetAll(ctx context.Context, userID uuid.UUID, limit int) ([]CachedNotification, error) {
	key := c.unreadKey(userID)

	if limit <= 0 {
		limit = c.maxUnread
	}

	vals, err := c.client.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	notifications := make([]CachedNotification, 0, len(vals))
	for _, v := range vals {
		var n CachedNotification
		if err := json.Unmarshal([]byte(v), &n); err != nil {
			continue
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (c *NotificationCache) MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error {
	key := c.unreadKey(userID)

	vals, err := c.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, v := range vals {
		var n CachedNotification
		if err := json.Unmarshal([]byte(v), &n); err != nil {
			continue
		}
		if n.ID == notificationID {
			return c.client.LRem(ctx, key, 1, v).Err()
		}
	}

	return nil
}

func (c *NotificationCache) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	key := c.unreadKey(userID)

	vals, err := c.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}

	if len(vals) == 0 {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}

func (c *NotificationCache) InvalidateUser(ctx context.Context, userID uuid.UUID) error {
	return c.client.Del(ctx, c.unreadKey(userID)).Err()
}

func (c *NotificationCache) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	key := c.unreadKey(userID)
	vals, err := c.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, v := range vals {
		var n CachedNotification
		if err := json.Unmarshal([]byte(v), &n); err != nil {
			continue
		}
		if !n.IsRead {
			count++
		}
	}

	return count, nil
}

func (c *NotificationCache) unreadKey(userID uuid.UUID) string {
	return fmt.Sprintf("notifications:unread:%s", userID.String())
}
