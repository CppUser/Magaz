package tgconfig

import (
	"Magaz/pkg/utils/parser"
	"log"
)

//TODO: move bot configs from api to bot package to separate concerns

// FSMConfig represents the structure of the FSM configuration in YAML
type FSMConfig struct {
	InitialState string                    `mapstructure:"initial_state"`
	States       map[string]FSMStateConfig `mapstructure:"states"`
}

// FSMStateConfig represents the configuration for each state
type FSMStateConfig struct {
	Handler     string            `mapstructure:"handler"`
	Transitions map[string]string `mapstructure:"transitions"`
}

// TODO: make more generic to load any config , from any package call (i.e. api uses same logic to load config)
func LoadConfig(name string, configType string, path []string) (*FSMConfig, error) {
	var cfg FSMConfig

	if err := parser.Load(name, configType, path, &cfg); err != nil {
		log.Fatalf("Failed to load Bot configs: %v", err)
	}

	return &cfg, nil
}
