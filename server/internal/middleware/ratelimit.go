package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/pkg/response"
)

type bucket struct {
	tokens float64
	last   time.Time
}

var (
	mu      sync.Mutex
	buckets = map[string]*bucket{}
)

// RateLimit 按 IP 限流：每分钟最多 n 次
func RateLimit(n int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: float64(n), last: now}
			buckets[ip] = b
		}
		elapsed := now.Sub(b.last).Minutes()
		b.tokens += elapsed * float64(n)
		if b.tokens > float64(n) {
			b.tokens = float64(n)
		}
		b.last = now
		if b.tokens < 1 {
			mu.Unlock()
			response.Err(c, 42901, "请求太频繁，稍后再试")
			c.Abort()
			return
		}
		b.tokens--
		mu.Unlock()

		c.Next()
	}
}
