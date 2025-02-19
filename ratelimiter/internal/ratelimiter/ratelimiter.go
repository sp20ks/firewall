package ratelimiter

import (
	"context"
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type IPRateLimiter struct {
	limiter *redis_rate.Limiter
	rdb     *redis.Client
}

func NewIPRateLimiter(redisAddr string) *IPRateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	limiter := redis_rate.NewLimiter(rdb)

	return &IPRateLimiter{rdb: rdb, limiter: limiter}
}

func (r *IPRateLimiter) Allow(clientIP string) error {
	ctx := context.Background()

	res, err := r.limiter.Allow(ctx, clientIP, redis_rate.PerMinute(10))
	if err != nil {
		return err
	}
	if res.Remaining == 0 {
		return fmt.Errorf("rate limit exceeded")
	}

	return nil
}
