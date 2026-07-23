// Package logger builds the zap.Logger used across services: a colorized,
// human-readable console encoder for local development, and structured JSON
// once APP_ENV=production so log aggregators still get parseable output.
package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var registerPrettyEncoderOnce sync.Once

func New(service string) (*zap.Logger, error) {
	var (
		log *zap.Logger
		err error
	)

	if os.Getenv("APP_ENV") == "production" {
		log, err = zap.NewProduction()
	} else {
		log, err = newDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return log.Named(service), nil
}

func newDevelopment() (*zap.Logger, error) {
	registerPrettyEncoderOnce.Do(func() {
		_ = zap.RegisterEncoder("pretty", newPrettyEncoder)
	})

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableStacktrace: true,
		Encoding:          "pretty",
		EncoderConfig:     zapcore.EncoderConfig{},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return cfg.Build()
}
