package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type ctxKey string

const txKey = ctxKey("transaction")

type Database struct {
	db *sql.DB
}

func New(url string) (*Database, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) executeTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) error {
	if tx := db.executeTx(ctx); tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}

	_, err := db.db.ExecContext(ctx, query, args...)
	return err
}

func (db *Database) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := db.executeTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := db.executeTx(ctx); tx != nil {
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
