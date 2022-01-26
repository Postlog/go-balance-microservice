package database

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/postlog/go-balance-microservice/internal/logger"
	"github.com/qustavo/sqlhooks/v2"
)

type ctxKey string

const txKey = ctxKey("transaction")

type Database struct {
	db *sql.DB
}

func New(dsn string, logger logger.Logger) (*Database, error) {
	sql.Register("psql", sqlhooks.Wrap(&pq.Driver{}, NewLogHook(logger)))
	db, err := sql.Open("psql", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) error {
	if tx := extractTx(ctx); tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}

	_, err := db.db.ExecContext(ctx, query, args...)
	return err
}

func (db *Database) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *Database) WithTransaction(ctx context.Context, level sql.IsolationLevel, f func(ctx context.Context) error) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
	if err != nil {
		return err
	}

	err = f(context.WithValue(ctx, txKey, tx))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}
