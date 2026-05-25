package eventbus

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLogger(
	l *zap.Logger,
) {
	logger = l
}

func getLogger() *zap.Logger {

	if logger == nil {
		return zap.NewNop()
	}

	return logger
}