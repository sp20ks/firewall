package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env               string `yaml:"env" env:"ENV" env-required:"true"`
	RulesEngineServer `yaml:"rules_engine_server"`
	RulesEngineDB     `yaml:"rules_engine_db"`
	AuthURL           string `yaml:"auth_url"`
}

type RulesEngineServer struct {
	Address string `yaml:"address" env-default:"localhost:8084"`
}

type RulesEngineDB struct {
	Host     string `yaml:"host" env-default:"postgres-rules-engine"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"admin"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-default:"postgres-rules-engine"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("RULES_ENGINE_CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("RULES_ENGINE_CONFIG_PATH is not set")
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
