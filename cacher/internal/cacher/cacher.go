package cacher

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher struct {
	rdb *redis.Client
}

func NewCacher(redisAddr string) (*Cacher, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cacher{rdb: rdb}, nil
}

func (c *Cacher) GetCache(key string) (*string, error) {
	val, err := c.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch cache: %w", err)
	}

	return &val, nil
}

func (c *Cacher) SetCahce(key string) (string, error) {

}

func (c *Cacher) Close() error {
	return c.rdb.Close()
}
