package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App App `json:"app" envPrefix:"APP_"`
	Log Log `json:"log" envPrefix:"LOG_"`

	Server      Server `json:"server" envPrefix:"SERVER_"`
	OpenapiPath string `json:"openapi_path" env:"OPENAPI_PATH"`

	Postgres Postgres `json:"postgres" envPrefix:"POSTGRES_"`
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

type Postgres struct {
	RWDSN       string        `json:"rw_dsn" env:"RW_DSN"`
	RODSN       string        `json:"ro_dsn" env:"RO_DSN"`
	PingTimeout time.Duration `json:"ping_timeout" env:"PING_TIMEOUT" envDefault:"5s"`
}

func GetDefault() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to parse: %w", err)
	}

	return &cfg, nil
}
