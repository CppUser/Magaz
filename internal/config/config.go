package config

import (
	"Magaz/pkg/utils/parser"
	"go.uber.org/zap"
	"log"
	"time"
)

// TODO:FIX: some of the code uses config as storage like Logger
type APIConfig struct {
	Version  string         `mapstructure:"version"`
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   *zap.Logger    `mapstructure:"logger"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Database DatabaseConfig `mapstructure:"database"`
	Bot      TGBotConfig    `mapstructure:"tg_bot"`
	CacheDir TemplateCache  `mapstructure:"cache_dir"`
	ScrKey   string         `mapstructure:"scr_key"`
}

// RedisConfig holds the Redis configuration.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	//TODO: add more redis configuration
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

type TemplateCache struct {
	Layouts string `mapstructure:"layouts"`
	Pages   string `mapstructure:"pages"`
}

// TODO: Move to bot config logic
// TGBotConfig holds the Telegram bot configuration.
type TGBotConfig struct {
	Token       string `mapstructure:"token"`
	WebhookLink string `mapstructure:"webhook_link"`
	WebhookPath string `mapstructure:"webhook_path"`
}

// TODO: make more generic to load any config , from any package call (i.e. bot telegram uses same logic to load config)
func LoadConfig() (*APIConfig, error) {
	var cfg APIConfig

	configPaths := []string{
		".",
		"config/",
	}

	if err := parser.Load("api_config", "yaml", configPaths, &cfg); err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	return &cfg, nil
}
