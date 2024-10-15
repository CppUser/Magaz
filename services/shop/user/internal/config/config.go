package config

import (
	"go.uber.org/zap"
	"log"
	"time"
	"user/internal/utils/parser"
)

type APIConfig struct {
	Version  string         `mapstructure:"version"`
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   *zap.Logger    `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
}

// DatabaseConfig holds the database configuration for PostgreSQL.
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// ServerConfig holds the HTTP server configuration.
type ServerConfig struct {
	Host          string        `mapstructure:"host"`
	Port          string        `mapstructure:"port"`
	TimeoutMS     time.Duration `mapstructure:"timeout_ms"`
	IdleTimeoutMS time.Duration `mapstructure:"idle_timeout_ms"`
}

func LoadAPIConfig() (*APIConfig, error) {
	var cfg APIConfig

	configPaths := []string{
		".",
		"config/",
		"/app/config/",
	}

	if err := parser.Load("api_config", "yaml", configPaths, &cfg); err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	return &cfg, nil
}
