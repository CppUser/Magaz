package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: need to add the logger configuration
// TODO: handle the error
// InitLogger initializes and returns a Zap logger.
func InitLogger(env string) (*zap.Logger, error) {
	var logger *zap.Logger

	switch env {
	case "dev":
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	case "prod":
		//TODO: Need to convert to JSON format
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		logger, _ = config.Build()

	}

	defer logger.Sync()

	return logger, nil
}
