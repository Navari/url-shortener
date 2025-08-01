package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func GetLogger() *zap.Logger {
	return Logger
}

func Info(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	Logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	Logger.Fatal(message, fields...)
}

func Sync() {
	Logger.Sync()
}
