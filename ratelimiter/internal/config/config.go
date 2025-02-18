package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env               string `yaml:"env" env:"ENV" env-required:"true"`
	RateLimiterServer `yaml:"rate_limiter_server"`
}

type RateLimiterServer struct {
	Address    string  `yaml:"address" env-default:"localhost:8081"`
	MaxTokens  float64 `yaml:"max_tokens" env-default:"5"`
	RefillRate float64 `yaml:"refill_rate" env-default:"0.5"`
	RedisAddr  string  `yaml:"redis_address"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("RATELIMITER_CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("RATELIMITER_CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil
}
