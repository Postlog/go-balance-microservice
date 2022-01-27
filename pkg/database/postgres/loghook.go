package postgres

import (
	"context"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"time"
)

type logHook struct {
	logger logger.Logger
}

func (h *logHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "started", time.Now()), nil
}

func (h *logHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	h.logger.Infof("Query: `%s`, Args: `%q`. took: %s", query, args, time.Since(ctx.Value("started").(time.Time)))
	return ctx, nil
}

func (h *logHook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	h.logger.Errorf("ErrorMessage: %v, Query: `%s`, Args: `%q`, Took: %s", err, query, args, time.Since(ctx.Value("started").(time.Time)))
	return err
}
