package middleware

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// AIRateLimiter schedules AI receipt extraction calls under a per-user daily cap.
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

// Reserve returns the earliest time at which a submitted AI job may run.
// Jobs beyond the current daily allowance are assigned to subsequent 24-hour
// windows instead of being rejected, so uploads can be persisted and queued.
func (rl *AIRateLimiter) Reserve(userID uuid.UUID) (time.Time, bool) {
	if rl == nil || rl.perDay <= 0 {
		return time.Time{}, false
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, ok := rl.hits[userID]
	if !ok {
		c = &aiCounter{windowFrom: now}
		rl.hits[userID] = c
	} else {
		advanceAICounter(c, now, rl.perDay)
	}

	window := c.count / rl.perDay
	c.count++
	if window == 0 {
		return time.Time{}, false
	}
	return c.windowFrom.Add(time.Duration(window) * 24 * time.Hour), true
}

func advanceAICounter(c *aiCounter, now time.Time, perDay int) {
	elapsedWindows := int(now.Sub(c.windowFrom) / (24 * time.Hour))
	if elapsedWindows <= 0 {
		return
	}
	c.count -= elapsedWindows * perDay
	if c.count < 0 {
		c.count = 0
	}
	c.windowFrom = c.windowFrom.Add(time.Duration(elapsedWindows) * 24 * time.Hour)
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

// Allow is kept for callers that only need an immediate yes/no decision.
// New intake paths should use Reserve so over-capacity work can be queued.
func (rl *AIRateLimiter) Allow(userID uuid.UUID) bool {
	if rl == nil || rl.perDay <= 0 {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, ok := rl.hits[userID]
	if !ok {
		c = &aiCounter{windowFrom: now}
		rl.hits[userID] = c
	} else {
		advanceAICounter(c, now, rl.perDay)
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
			advanceAICounter(c, now, rl.perDay)
			if c.count == 0 {
				delete(rl.hits, id)
			}
		}
		rl.mu.Unlock()
	}
}
