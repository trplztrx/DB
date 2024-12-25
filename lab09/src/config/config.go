package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig `yaml:"postgres"`
	RedisConfig `yaml:"redis"`
}

type DBConfig struct {
	Host         string `yaml:"host" env:"DB_HOST"`
	Port         int    `yaml:"port" env:"DB_PORT"`
	User         string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password     string `yaml:"password" env:"DB_PASS" env-default:"postgres"`
	DatabaseName string `yaml:"db" env:"DB_NAME" env-default:"postgres"`
}

type RedisConfig struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
}

func ReadConfig() (*Config, error) {
	cfg := Config{}
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %v", err)
	}

	configPath := filepath.Join(dir, "config", "config.yaml")

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read env: %v", err)
	}

	return &cfg, nil
}
