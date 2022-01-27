package database

import (
	"context"
	"database/sql"
	"github.com/postlog/go-balance-microservice/pkg/database/postgres"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"github.com/postlog/go-balance-microservice/pkg/types"
)

type Database interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	WithTransaction(ctx context.Context, level sql.IsolationLevel, f func(ctx context.Context) error) error
	GetTransactionWrapper() types.TransactionWrapper
	Close() error
}

func NewDatabase(dsn string, logger logger.Logger) (Database, error) {
	return postgres.NewPostgresDatabase(dsn, logger)
}
