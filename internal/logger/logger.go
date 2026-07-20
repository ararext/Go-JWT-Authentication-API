package logger

import "go.uber.org/zap"

func New() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return logger
}