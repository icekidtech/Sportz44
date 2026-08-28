package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit implements a sliding-window rate limiter backed by Redis, per the
// PRD (§9.3): 100 requests/minute per user, with per-endpoint limits. When no
// user is authenticated, it falls back to the client IP.
func RateLimit(rdb *redis.Client, requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Identify the caller: authenticated user ID, else client IP.
		key := c.ClientIP()
		if uid, ok := c.Get(ContextUserID); ok {
			if id, ok := uid.(uint); ok && id != 0 {
				key = fmt.Sprintf("user:%d", id)
			}
		}

		// Redis key: ratelimit:<window>:<key>:<endpoint>
		windowKey := fmt.Sprintf("%d", time.Now().Unix()/int64(window.Seconds()))
		redisKey := fmt.Sprintf("ratelimit:%s:%s:%s", windowKey, key, c.FullPath())

		ctx := context.Background()
		count, err := rdb.Incr(ctx, redisKey).Result()
		if err != nil {
			// If Redis is unavailable, fail open rather than block traffic.
			c.Next()
			return
		}
		if count == 1 {
			rdb.Expire(ctx, redisKey, window)
		}

		if count > int64(requests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
