package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env:"ENV" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
	IpsFilePath    string     `env:"IPS_PATH"`
	Resources      []resource `yaml:"resources"`
	RateLimiterURL string     `yaml:"ratelimiter_url"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type resource struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Host     string `yaml:"host"`
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
