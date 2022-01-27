package logger

import (
	"context"
	"encoding/json"
	"github.com/postlog/go-balance-microservice/pkg/utils"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	// With returns new logger, that holds passed arguments and request id, fetched from the context
	With(ctx context.Context, args ...interface{}) Logger
	// Flush flushes logger's buffers
	Flush() error
}

func New(config json.RawMessage) (Logger, error) {
	l, err := newZapLogger(config)
	if err != nil {
		return nil, err
	}
	return l, nil
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

func newZapLogger(config json.RawMessage) (Logger, error) {
	zapLogger, err := initLoggerFromConfig(config)
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
