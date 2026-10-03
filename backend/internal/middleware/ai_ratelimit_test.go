package middleware

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIRateLimiterReserveQueuesBeyondDailyCapacity(t *testing.T) {
	limiter := NewAIRateLimiter(2)
	userID := uuid.New()

	_, queued := limiter.Reserve(userID)
	assert.False(t, queued)
	_, queued = limiter.Reserve(userID)
	assert.False(t, queued)

	processAfter, queued := limiter.Reserve(userID)
	require.True(t, queued)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), processAfter, time.Second)
}

func TestAIRateLimiterReserveReleasesElapsedCapacity(t *testing.T) {
	limiter := NewAIRateLimiter(2)
	userID := uuid.New()
	limiter.hits[userID] = &aiCounter{
		count:      3,
		windowFrom: time.Now().Add(-24*time.Hour - time.Minute),
	}

	_, queued := limiter.Reserve(userID)
	assert.False(t, queued)
}
