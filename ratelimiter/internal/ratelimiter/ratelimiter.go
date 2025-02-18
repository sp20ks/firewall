package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens         float64   // Current number of tokens
	maxTokens      float64   // Maximum tokens allowed
	refillRate     float64   // Tokens added per second
	lastRefillTime time.Time // Last time tokens were refilled
	mutex          sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate float64) *RateLimiter {
	return &RateLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (r *RateLimiter) refillTokens() {
	now := time.Now()
	duration := now.Sub(r.lastRefillTime).Seconds()
	tokensToAdd := duration * r.refillRate

	r.tokens += tokensToAdd
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefillTime = now
}

func (r *RateLimiter) Allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refillTokens()

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

type IPRateLimiter struct {
	limiters   map[string]*RateLimiter
	mutex      sync.Mutex
	maxTokens  float64
	refillRate float64
}

func NewIPRateLimiter(maxTokens, refillRate float64) *IPRateLimiter {
	return &IPRateLimiter{
		limiters:   make(map[string]*RateLimiter),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = NewRateLimiter(i.maxTokens, i.refillRate)
		i.limiters[ip] = limiter
	}

	return limiter
}
