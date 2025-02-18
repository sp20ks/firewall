package ratelimiter

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redis_rate/v10"
)

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

func NewIPRateLimiter(maxTokens, refillRate float64, redisAddr string) *IPRateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &IPRateLimiter{
		redisClient: rdb,
		maxTokens:   maxTokens,
		refillRate:  refillRate,
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
