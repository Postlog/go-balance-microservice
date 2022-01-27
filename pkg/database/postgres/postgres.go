package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"github.com/postlog/go-balance-microservice/pkg/types"
	"github.com/qustavo/sqlhooks/v2"
)

type postgres struct {
	db *sql.DB
}

func NewPostgresDatabase(dsn string, logger logger.Logger) (*postgres, error) {
	name := "postgresWithLogging"
	sql.Register(name, sqlhooks.Wrap(&pq.Driver{}, &logHook{logger}))
	db, err := sql.Open(name, dsn)
	if err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

type ctxKey string

const txKey = ctxKey("transaction")

func (pg *postgres) Close() error {
	return pg.db.Close()
}

func (pg *postgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	if tx := extractTx(ctx); tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}

	_, err := pg.db.ExecContext(ctx, query, args...)
	return err
}

func (pg *postgres) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return pg.db.QueryContext(ctx, query, args...)
}

func (pg *postgres) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return pg.db.QueryRowContext(ctx, query, args...)
}

func (pg *postgres) WithTransaction(ctx context.Context, level sql.IsolationLevel, f func(ctx context.Context) error) error {
	tx, err := pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
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

func (pg *postgres) GetTransactionWrapper() types.TransactionWrapper {
	return func(ctx context.Context, f func(ctx context.Context) error) error {
		return pg.WithTransaction(ctx, sql.LevelDefault, f)
	}
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}
