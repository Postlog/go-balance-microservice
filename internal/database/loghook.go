package database

import (
	"context"
	"github.com/postlog/go-balance-microservice/internal/logger"
	"time"
)

type LogHook struct {
	logger logger.Logger
}

func NewLogHook(logger logger.Logger) *LogHook {
	return &LogHook{logger}
}
func (h *LogHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "started", time.Now()), nil
}

func (h *LogHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	h.logger.Infof("Query: `%s`, Args: `%q`. took: %s", query, args, time.Since(ctx.Value("started").(time.Time)))
	return ctx, nil
}

func (h *LogHook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	h.logger.Errorf("Error: %v, Query: `%s`, Args: `%q`, Took: %s", err, query, args, time.Since(ctx.Value("started").(time.Time)))
	return err
}