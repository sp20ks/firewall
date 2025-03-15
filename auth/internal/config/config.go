package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	AuthServer `yaml:"auth_server"`
}

type AuthServer struct {
	Address   string `yaml:"address" env-default:"localhost:8083"`
	SecretKey string `yaml:"secret" env-default:"secret-key"`
	KeyTTL    int    `yaml:"key_ttl"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("AUTH_CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("AUTH_CONFIG_PATH is not set")
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
