package config

import (
	"errors"
	"fmt"
	"maps"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host" env:"APP_URL"`
		Port string `mapstructure:"port" env:"APP_PORT"`
	} `mapstructure:"server"`
	Env                string   `mapstructure:"env" env:"APP_ENV"`
	Debug              bool     `mapstructure:"debug" env:"APP_DEBUG"`
	CorsAllowedOrigins []string `mapstructure:"cors_allowed_origins" env:"APP_CORS_ALLOWED_ORIGINS"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" env:"POSTGRES_HOST"`
	Port     string `mapstructure:"port" env:"POSTGRES_PORT"`
	User     string `mapstructure:"user" env:"POSTGRES_USER"`
	Password string `mapstructure:"password" env:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"dbname" env:"POSTGRES_DB"`
	SSLMode  string `mapstructure:"sslmode" env:"POSTGRES_SSLMODE"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

func NewConfig() *Config {
	cfg := &Config{
		App: AppConfig{
			Env:                "development",
			Debug:              true,
			CorsAllowedOrigins: []string{"*"},
		},
	}
	cfg.App.Server.Host = "localhost"
	cfg.App.Server.Port = "8080"
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = "5432"
	cfg.Postgres.User = "postgres"
	cfg.Postgres.Password = "password"
	cfg.Postgres.DBName = "postgres"
	cfg.Postgres.SSLMode = "disable"
	return cfg
}

// Load reads app settings from the optional .env file and the process environment.
func Load() (*Config, error) {
	values, err := godotenv.Read()
	if errors.Is(err, os.ErrNotExist) {
		values = make(map[string]string)
	} else if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return nil, fmt.Errorf("read .env: %w", err)
		}
		// Parser errors can include the original line, which may contain a secret.
		return nil, errors.New("invalid .env file syntax")
	}
	maps.Copy(values, env.ToMap(os.Environ()))

	cfg := NewConfig()
	if err := env.ParseWithOptions(cfg, env.Options{Environment: values}); err != nil {
		return nil, fmt.Errorf("parse app configuration: %w", err)
	}
	return cfg, nil
}
