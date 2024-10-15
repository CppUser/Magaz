package config

import (
	"log"
	"tg/internal/utlis/parser"
)

// RedisConfig holds the Redis configuration.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	//TODO: add more redis configuration
}

type BotConf struct {
	Host  string      `mapstructure:"host"`
	Port  string      `mapstructure:"port"`
	Bot   botsConfig  `mapstructure:"tg_bot"`
	Redis RedisConfig `mapstructure:"redis"`
}

type botsConfig struct {
	WebhookLink     string      `mapstructure:"webhook_link"`
	WebhookBasePath string      `mapstructure:"webhook_base_path"`
	Tokens          []BotTokens `mapstructure:"tokens"`
	Groups          []botGroups `mapstructure:"groups"`
}

type BotTokens struct {
	Type  string `mapstructure:"type"`
	Token string `mapstructure:"id"`
}

type botGroups struct {
	Type    string `mapstructure:"type"`
	GroupID int64  `mapstructure:"id"` // Updated to use `int64` for group IDs
}

// TODO: Refactor loading bot config function
func LoadBotConfig() (*BotConf, error) {
	var cfg BotConf

	configPaths := []string{
		".",
		"config/",
		"/app/config/",
		//TODO: Refactor (find other way to avoid using hard coded path)
	}

	if err := parser.Load("bot_config", "yaml", configPaths, &cfg); err != nil {
		log.Fatalf("Failed to load Bot configs: %v", err)
	}

	return &cfg, nil
}
