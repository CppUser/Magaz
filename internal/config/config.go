package config

import (
	"Magaz/pkg/utils/parser"
	"go.uber.org/zap"
	"log"
	"time"
)

type APIConfig struct {
	Version   string       `mapstructure:"version"`
	Env       string       `mapstructure:"env"`
	Server    ServerConfig `mapstructure:"server"`
	Logger    *zap.Logger
	BotConfig TGBotConfig `mapstructure:"tg_bot"`
}

// ServerConfig holds the HTTP server configuration.
type ServerConfig struct {
	Host          string        `mapstructure:"host"`
	Port          string        `mapstructure:"port"`
	TimeoutMS     time.Duration `mapstructure:"timeout_ms"`
	IdleTimeoutMS time.Duration `mapstructure:"idle_timeout_ms"`
}

// TGBotConfig holds the Telegram bot configuration.
type TGBotConfig struct {
	Token       string `mapstructure:"token"`
	WebhookLink string `mapstructure:"webhook_link"`
	WebhookPath string `mapstructure:"webhook_path"`
}

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
