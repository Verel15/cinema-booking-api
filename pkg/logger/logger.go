package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()

	// Customize encoder config to match user requirements
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"

	// Create logger with global fields
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	logger, err := config.Build(zap.Fields(
		zap.String("service", "cinema-booking-api"),
		zap.String("env", env),
	))
	if err != nil {
		panic(err)
	}

	Log = logger
}