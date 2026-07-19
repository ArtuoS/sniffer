package container

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func LoadLogger() error {
	l, err := zap.NewProduction()
	if err != nil {
		return err
	}

	logger = l
	return nil
}

func GetLogger() *zap.Logger {
	return logger
}

func SyncLogger() {
	logger.Sync()
}
