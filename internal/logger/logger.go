package logger

import (
	"context"
	"encoding/json"
	"github.com/postlog/go-balance-microservice/internal/config"
	"github.com/postlog/go-balance-microservice/internal/utils"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	With(ctx context.Context, args ...interface{}) Logger
	Flush() error
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) With(ctx context.Context, args ...interface{}) Logger {
	if id, ok := utils.GetRequestID(ctx); ok {
		args = append(args, "request_id", id)
	}
	return &logger{l.SugaredLogger.With(args...)}
}

func (l *logger) Flush() error {
	return l.Sync()
}

func New(cfg *config.Config) (Logger, error) {
	l, err := newZapLogger(cfg)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func newZapLogger(cfg *config.Config) (Logger, error) {
	zapLogger, err := initLoggerFromConfig(cfg.Logger)
	if err != nil {
		return nil, err
	}
	return zapLogger, nil
}

func initLoggerFromConfig(jsonCfg json.RawMessage) (*logger, error) {
	cfg := zap.Config{}
	err := json.Unmarshal(jsonCfg, &cfg)
	if err != nil {
		return nil, err
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &logger{l.Sugar()}, nil
}
