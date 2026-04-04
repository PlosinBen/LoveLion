package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	tokens    int
	lastReset time.Time
}

// RateLimit returns a middleware that limits requests per IP.
// maxRequests: max requests allowed within the window.
// window: time window for the rate limit.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Cleanup stale entries periodically
	go func() {
		for {
			time.Sleep(window * 2)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastReset) > window*2 {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Since(v.lastReset) > window {
			visitors[ip] = &visitor{tokens: maxRequests - 1, lastReset: time.Now()}
			mu.Unlock()
			c.Next()
			return
		}

		if v.tokens <= 0 {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please try again later"})
			c.Abort()
			return
		}

		v.tokens--
		mu.Unlock()
		c.Next()
	}
}
