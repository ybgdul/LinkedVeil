package models

import (
	"sync"
	"time"
)

type RateLimiter struct{
	mu sync.Mutex
	tokens float64
	maxTokens float64
	refillRate float64
	lastRefill time.Time
}

func NewRateLimiter(maxTokens float64, refillRate float64) *RateLimiter { 
	return &RateLimiter{
		tokens: maxTokens,
		maxTokens: maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (r *RateLimiter) refillTokens() { 
	now := time.Now()
	duration := now.Sub(r.lastRefill).Seconds()
	tokensToAdd := duration * r.refillRate

	r.tokens += tokensToAdd
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}

	r.lastRefill = time.Now()
} 

func (r *RateLimiter) Allow() bool { 
	r.mu.Lock()

	defer r.mu.Unlock()

	r.refillTokens()
	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}