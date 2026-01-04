package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	buckets  map[string]*TokenBucket
	capacity float64
	rate     float64
	mu       sync.RWMutex
}

func NewRateLimiter(capacity, rate float64) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}

	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for id, bucket := range rl.buckets {
			if time.Since(bucket.LastAccessed()) > 5*time.Minute {
				delete(rl.buckets, id)
			}
		}
		rl.mu.Unlock()
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
