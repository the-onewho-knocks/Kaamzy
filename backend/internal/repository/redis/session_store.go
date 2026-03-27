package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionStore(client *redis.Client, ttl time.Duration) *SessionStore {
	return &SessionStore{client: client, ttl: ttl}
}

func (s *SessionStore) Set(ctx context.Context, token string, data SessionData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	key := s.key(token)
	return s.client.Set(ctx, key, b, s.ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, token string) (*SessionData, error) {
	key := s.key(token)
	val, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var data SessionData
	if err := json.Unmarshal(val, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	return s.client.Del(ctx, s.key(token)).Err()
}

func (s *SessionStore) Refresh(ctx context.Context, token string) error {
	return s.client.Expire(ctx, s.key(token), s.ttl).Err()
}

func (s *SessionStore) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("session:user:%s:*", userID.String())
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}

func (s *SessionStore) Exists(ctx context.Context, token string) (bool, error) {
	count, err := s.client.Exists(ctx, s.key(token)).Result()
	return count > 0, err
}

func (s *SessionStore) key(token string) string {
	return fmt.Sprintf("session:%s", token)
}