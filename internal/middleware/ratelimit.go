package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]visitor),
		limit:    limit,
		window:   window,
	}
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		now := time.Now()

		r.mu.Lock()

		entry, exists := r.visitors[key]

		if !exists || now.After(entry.resetAt) {
			entry = visitor{
				count:   0,
				resetAt: now.Add(r.window),
			}
		}

		entry.count++
		r.visitors[key] = entry

		exceeded := entry.count > r.limit

		r.mu.Unlock()

		if exceeded {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{"error": "too many requests"},
			)
			return
		}

		c.Next()
	}
}
