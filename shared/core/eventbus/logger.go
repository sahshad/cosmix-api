package eventbus

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLogger(log *zap.Logger) {
	logger = log
}

func getLogger() *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}

	return logger
}
