package balance

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/internal/database"
)

type UserBalance struct {
	UserId uuid.UUID `json:"userId"`
	Value  float64   `json:"balance"`
}

type Repository interface {
	Get(ctx context.Context, balance UserBalance) (UserBalance, error)
	GetWithSpec(ctx context.Context, balance UserBalance, spec Specification) (UserBalance, error)
	Create(ctx context.Context, balance UserBalance) error
	Update(ctx context.Context, balance UserBalance) error
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

func NewRepository(db *database.Database) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db *database.Database
}

var NotFoundErr = errors.New("user with specified id not found")

func (r *postgresRepository) Get(ctx context.Context, balance UserBalance) (UserBalance, error) {
	query := "SELECT value from balance WHERE user_id=$1"
	return r.get(ctx, balance, query)
}

func (r *postgresRepository) GetWithSpec(ctx context.Context, balance UserBalance, spec Specification) (UserBalance, error) {
	query := "SELECT value from balance WHERE user_id=$1" + " " + spec.ToSQLClosure()
	return r.get(ctx, balance, query)
}

func (r *postgresRepository) get(ctx context.Context, balance UserBalance, query string) (UserBalance, error) {
	err := r.db.QueryRow(ctx, query, balance.UserId).Scan(&balance.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return balance, NotFoundErr
		}
		return balance, err
	}

	return balance, nil
}

func (r *postgresRepository) Create(ctx context.Context, balance UserBalance) error {
	return r.db.Exec(ctx, "INSERT INTO balance VALUES ($1, $2)", balance.UserId, balance.Value)
}

func (r *postgresRepository) Update(ctx context.Context, balance UserBalance) error {
	return r.db.Exec(ctx, "UPDATE balance SET value=$1 WHERE user_id=$2", balance.Value, balance.UserId)
}

func (r *postgresRepository) Transaction(ctx context.Context, f func(ctx context.Context) error) error {
	return r.db.WithTransaction(ctx, sql.LevelReadCommitted, f)
}
