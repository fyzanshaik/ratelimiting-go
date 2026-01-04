package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   float64
	rate       float64
	tokens     float64
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tokensToAdd := elapsed * tb.rate

	tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
	tb.lastRefill = now
}

func (tb *TokenBucket) Take(tokens float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()

	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}

	return false
}

func (tb *TokenBucket) TakeWithTimeout(tokens float64, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		tb.mutex.Lock()
		tb.refill()

		if tb.tokens >= tokens {
			tb.tokens -= tokens
			tb.mutex.Unlock()
			return true
		}

		tokensNeeded := tokens - tb.tokens
		timeNeeded := time.Duration(tokensNeeded/tb.rate*1000) * time.Millisecond
		tb.mutex.Unlock()

		waitTime := min(timeNeeded, time.Until(deadline))
		if waitTime > 0 {
			time.Sleep(waitTime)
		}
	}

	return false
}

func (tb *TokenBucket) TakeWithContext(ctx context.Context, tokens float64) bool {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if tb.Take(tokens) {
				return true
			}
		}
	}
}

func (tb *TokenBucket) TakeWithBurstLimit(tokens, maxBurst float64) bool {
	if tokens > maxBurst {
		tokens = maxBurst
	}
	return tb.Take(tokens)
}

func (tb *TokenBucket) Metrics() (capacity, rate, currentTokens float64) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()
	return tb.capacity, tb.rate, tb.tokens
}

func (tb *TokenBucket) SetRate(newRate float64) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()
	tb.rate = newRate
}
