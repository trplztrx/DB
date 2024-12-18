package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig `yaml:"postgres"`
}

type DBConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password     string `yaml:"password" env:"DB_USER" env-default:"postgres"`
	DatabaseName string `yaml:"db" env:"DB_NAME" env-default:"postgres"`
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
