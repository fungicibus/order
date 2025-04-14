package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App App `json:"app" envPrefix:"APP_"`
	Log Log `json:"log" envPrefix:"LOG_"`

	Server Server `json:"server" envPrefix:"SERVER_"`
}

type App struct {
	Env     string `json:"env" env:"ENV" envDefault:"dev"`
	Name    string `json:"name" env:"NAME" envDefault:"order"`
	Version string `json:"version" env:"-"`
}

type Log struct {
	Level       int    `json:"level" env:"LEVEL" envDefault:"0"`
	VictoriaUrl string `json:"victoria_url" env:"VICTORIA_URL"`
}

type Server struct {
	Port         int           `json:"port" env:"PORT" envDefault:"8081"`
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
