package utils

import (
	"context"
	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey = ctxKey("request_id")

func GetRequestID(ctx context.Context) (string, bool) {
	if id, ok := ctx.Value(requestIDKey).(string); !ok {
		return "", false
	} else {
		return id, true
	}
}

func WithRequestID(ctx context.Context) context.Context {
	if ctx.Value(requestIDKey) != nil {
		return ctx
	}
	return context.WithValue(ctx, requestIDKey, uuid.New().String())
}
