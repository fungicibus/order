package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App App `json:"app" envPrefix:"APP_"`
	Log Log `json:"log" envPrefix:"LOG_"`

	Server      Server   `json:"server" envPrefix:"SERVER_"`
	OpenapiPath string   `json:"openapi_path" env:"OPENAPI_PATH"`
	Security    Security `json:"security" envPrefix:"SECURITY_"`

	Postgres Postgres `json:"postgres" envPrefix:"POSTGRES_"`
	Kafka    Kafka    `json:"kafka" envPrefix:"KAFKA_"`
}

type App struct {
	Env             string        `json:"env" env:"ENV" envDefault:"dev"`
	Name            string        `json:"name" env:"NAME" envDefault:"order"`
	Version         string        `json:"version" env:"-"`
	InitTimeout     time.Duration `json:"init_timeout" env:"INIT_TIMEOUT" envDefault:"10s"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
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

type Security struct {
	AdminKey string `json:"admin_key" env:"ADMIN_KEY"`
}

type Postgres struct {
	RWDSN       string        `json:"rw_dsn" env:"RW_DSN"`
	RODSN       string        `json:"ro_dsn" env:"RO_DSN"`
	PingTimeout time.Duration `json:"ping_timeout" env:"PING_TIMEOUT" envDefault:"5s"`
}

type Kafka struct {
	CommitInterval      time.Duration `json:"commit_interval" env:"COMMIT_INTERVAL" envDefault:"5s"`
	Brokers             []string      `json:"brokers" env:"BROKERS" envSeparator:","`
	TopicOrderCreated   string        `json:"topic_order_created" env:"TOPIC_ORDER_CREATED"`
	TopicOrderConfirmed string        `json:"topic_order_confirmed" env:"TOPIC_ORDER_CONFIRMED"`
	TopicOrderRejected  string        `json:"topic_order_rejected" env:"TOPIC_ORDER_REJECTED"`
}

func GetDefault() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to parse: %w", err)
	}

	return &cfg, nil
}
