package config

import (
	"Magaz/pkg/utils/parser"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
	"log"
	"time"
)

type APIConfig struct {
	Version string       `mapstructure:"version"`
	Env     string       `mapstructure:"env"`
	Server  ServerConfig `mapstructure:"server"`
	Logger  *zap.Logger
	Redis   RedisConfig `mapstructure:"redis"`
	Bot     TGBotConfig `mapstructure:"tg_bot"`
}

// RedisConfig holds the Redis configuration.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	//TODO: add more redis configuration
}

// ServerConfig holds the HTTP server configuration.
type ServerConfig struct {
	Host          string        `mapstructure:"host"`
	Port          string        `mapstructure:"port"`
	TimeoutMS     time.Duration `mapstructure:"timeout_ms"`
	IdleTimeoutMS time.Duration `mapstructure:"idle_timeout_ms"`
}

// TODO: Move to bot config logic
// TGBotConfig holds the Telegram bot configuration.
type TGBotConfig struct {
	API         *telego.Bot
	Logger      *zap.Logger //TODO: figure out how to use API logger initialization
	Token       string      `mapstructure:"token"`
	WebhookLink string      `mapstructure:"webhook_link"`
	WebhookPath string      `mapstructure:"webhook_path"`
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
