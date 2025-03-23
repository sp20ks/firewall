package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `env:"ENV"               env-required:"true" yaml:"env"`
	HTTPServer     `yaml:"http_server"`
	IpsFilePath    string `env:"IPS_PATH"`
	RateLimiterURL string `yaml:"ratelimiter_url"`
	CacherURL      string `yaml:"cacher_url"`
	RulesEngineURL string `yaml:"rules_engine_url"`
}

type HTTPServer struct {
	Address     string        `env-default:"localhost:8080" yaml:"address"`
	Timeout     time.Duration `env-default:"4s"             yaml:"timeout"`
	IdleTimeout time.Duration `env-default:"60s"            yaml:"idle_timeout"`
	User        string        `env-required:"true"          yaml:"user"`
	Password    string        `env:"HTTP_SERVER_PASSWORD"   env-required:"true" yaml:"password"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("PROXY_CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("PROXY_CONFIG_PATH is not set")
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
