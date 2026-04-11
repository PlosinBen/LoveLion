package middleware

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// AIRateLimiter enforces a per-user daily cap on AI receipt extraction calls.
// The window is a rolling 24h; state is in-memory only (single-instance worker
// model) which matches the rest of the deployment.
type AIRateLimiter struct {
	perDay int
	mu     sync.Mutex
	hits   map[uuid.UUID]*aiCounter
}

type aiCounter struct {
	count      int
	windowFrom time.Time
}

// NewAIRateLimiter builds a limiter with the given per-user daily cap.
// A zero or negative cap disables the check.
func NewAIRateLimiter(perDay int) *AIRateLimiter {
	rl := &AIRateLimiter{
		perDay: perDay,
		hits:   make(map[uuid.UUID]*aiCounter),
	}

	if perDay > 0 {
		go rl.gcLoop()
	}
	return rl
}

// Allow reports whether userID has remaining quota and reserves one call if so.
func (rl *AIRateLimiter) Allow(userID uuid.UUID) bool {
	if rl == nil || rl.perDay <= 0 {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, ok := rl.hits[userID]
	if !ok || now.Sub(c.windowFrom) > 24*time.Hour {
		rl.hits[userID] = &aiCounter{count: 1, windowFrom: now}
		return true
	}
	if c.count >= rl.perDay {
		return false
	}
	c.count++
	return true
}

func (rl *AIRateLimiter) gcLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for id, c := range rl.hits {
			if now.Sub(c.windowFrom) > 24*time.Hour {
				delete(rl.hits, id)
			}
		}
		rl.mu.Unlock()
	}
}
