package ratelimiter

import (
	"context"
	"fmt"
	"time"

	redis_rate "github.com/go-redis/redis_rate/v10"
	redis "github.com/redis/go-redis/v9"
)

type IPRateLimiter struct {
	limiter *redis_rate.Limiter
	rdb     *redis.Client
}

func NewIPRateLimiter(redisAddr string) (*IPRateLimiter, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	limiter := redis_rate.NewLimiter(rdb)
	return &IPRateLimiter{rdb: rdb, limiter: limiter}, nil
}

func (r *IPRateLimiter) Allow(ctx context.Context, clientIP string) error {
	res, err := r.limiter.Allow(ctx, clientIP, redis_rate.PerMinute(10))
	if err != nil {
		return err
	}
	if res.Remaining == 0 {
		return fmt.Errorf("rate limit exceeded")
	}

	return nil
}

func (r *IPRateLimiter) Close() error {
	return r.rdb.Close()
}
