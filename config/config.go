package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server   Server `json:"server" envPrefix:"SERVER_"`
	LogLevel int    `json:"log_level" env:"LOG_LEVEL"`
}

type Server struct {
	Port         int           `json:"port" env:"PORT"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `json:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"5s"`
}

func GetDefault() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to parse: %w", err)
	}

	return &cfg, nil
}
