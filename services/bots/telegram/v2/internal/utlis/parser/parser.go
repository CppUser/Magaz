package parser

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

// Load loads the configuration from a file and unmarshals it into the provided struct.
func Load(configName string, configType string, configPaths []string, configStruct interface{}) error {
	viper.SetConfigName(configName) // Name of the configs file (without extension)
	viper.SetConfigType(configType) // Config file type

	// Add provided configs paths
	for _, path := range configPaths {
		viper.AddConfigPath(filepath.Clean(path))
	}

	// Read the configs file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading configs file: %w", err)
	}

	// Unmarshal the configs into the provided struct
	if err := viper.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("error unmarshalling configs: %w", err)
	}

	return nil
}
