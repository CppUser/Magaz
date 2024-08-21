package tgconfig

import (
	"Magaz/pkg/utils/parser"
	"log"
)

//TODO: move bot configs from api to bot package to separate concerns

type Transition struct {
	Event string `mapstructure:"event"`
	To    string `mapstructure:"to"`
}

type State struct {
	Name        string       `mapstructure:"name"`
	Transitions []Transition `mapstructure:"transitions"`
}

type Handler struct {
	Event   string `mapstructure:"event"`
	Handler string `mapstructure:"handler"`
}

type FSMConfig struct {
	States   []State   `mapstructure:"states"`
	Handlers []Handler `mapstructure:"handlers"`
}

// TODO: make more generic to load any config , from any package call (i.e. api uses same logic to load config)
func LoadConfig(name string, configType string, path []string) (*FSMConfig, error) {
	var cfg FSMConfig

	if err := parser.Load(name, configType, path, &cfg); err != nil {
		log.Fatalf("Failed to load Bot configs: %v", err)
	}

	return &cfg, nil
}
