package ratelimiter

import (
	"sync"
)

type RateLimiter struct {
	buckets  map[string]*TokenBucket
	capacity float64
	rate     float64
	mu       sync.RWMutex
}

func NewRateLimiter(capacity, rate float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

func (rl *RateLimiter) GetBucket(clientID string) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, exists := rl.buckets[clientID]; exists {
		return bucket
	}

	bucket = NewTokenBucket(rl.capacity, rl.rate)
	rl.buckets[clientID] = bucket
	return bucket
}

func (rl *RateLimiter) Allow(clientID string) bool {
	return rl.GetBucket(clientID).Take(1)
}

func (rl *RateLimiter) AllowN(clientID string, n float64) bool {
	return rl.GetBucket(clientID).Take(n)
}

func (rl *RateLimiter) GetMetrics(clientID string) (capacity, rate, currentTokens float64) {
	return rl.GetBucket(clientID).Metrics()
}
