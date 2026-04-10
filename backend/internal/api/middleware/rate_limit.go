package middleware

import (
	"net/http"
	"strings"
	"time"

	"kaamzy/internal/repository/redis"
	"kaamzy/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RateLimiter struct {
	store *redis.RateLimitStore
}

func NewRateLimiter(store *redis.RateLimitStore) *RateLimiter {
	store.AddRule("global:default", 100, time.Minute)
	store.AddRule("global:strict", 10, time.Minute)
	store.AddRule("auth:login", 5, time.Minute)
	store.AddRule("auth:register", 3, time.Minute)
	store.AddRule("api:jobs", 30, time.Minute)
	store.AddRule("api:messages", 60, time.Minute)

	return &RateLimiter{store: store}
}

func (r *RateLimiter) Middleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := r.getIdentifier(c, key)
		allowed, remaining, err := r.store.Allow(c.Request.Context(), identifier)

		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Remaining", formatRemaining(remaining))
		c.Header("X-RateLimit-Limit", key)

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(
				"rate limit exceeded, please try again later",
			))
			return
		}

		c.Next()
	}
}

func (r *RateLimiter) IPRateLimit() gin.HandlerFunc {
	return r.Middleware("global:default")
}

func (r *RateLimiter) AuthRateLimit() gin.HandlerFunc {
	return r.Middleware("auth:login")
}

func (r *RateLimiter) UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			r.IPRateLimit()(c)
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			r.IPRateLimit()(c)
			return
		}

		identifier := "user:" + uid.String() + ":default"
		allowed, remaining, err := r.store.Allow(c.Request.Context(), identifier)

		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Remaining", formatRemaining(remaining))
		c.Header("X-RateLimit-Limit", "user:default")

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(
				"rate limit exceeded, please try again later",
			))
			return
		}

		c.Next()
	}
}

func (r *RateLimiter) EndpointRateLimit(endpoint string) gin.HandlerFunc {
	return r.Middleware("api:" + endpoint)
}

func (r *RateLimiter) getIdentifier(c *gin.Context, key string) string {
	if userID, exists := c.Get("userID"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			return "user:" + uid.String() + ":" + key
		}
	}

	ip := c.ClientIP()
	return "ip:" + ip + ":" + key
}

func formatRemaining(remaining int) string {
	if remaining < 0 {
		return "0"
	}
	return string(rune('0' + remaining%10))
}

func getClientIP(c *gin.Context) string {
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		return xri
	}

	return c.ClientIP()
}
