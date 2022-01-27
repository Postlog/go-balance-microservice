package databaserepository

import (
	"context"
	"database/sql"
	error2 "github.com/postlog/go-balance-microservice/dataservice/balance/error"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/pkg/database"
)

// databaseRepository implements balance.Repository interface using Database
type databaseRepository struct {
	db database.Database
}

func New(db database.Database) *databaseRepository {
	return &databaseRepository{db}
}

func (r *databaseRepository) Get(ctx context.Context, balance models.UserBalance) (models.UserBalance, error) {
	query := "SELECT value from balance WHERE user_id=$1"
	return r.get(ctx, balance, query)
}

func (r *databaseRepository) GetAndBlock(ctx context.Context, balance models.UserBalance) (models.UserBalance, error) {
	query := "SELECT value from balance WHERE user_id=$1 FOR UPDATE"
	return r.get(ctx, balance, query)
}

func (r *databaseRepository) get(ctx context.Context, b models.UserBalance, query string) (models.UserBalance, error) {
	err := r.db.QueryRow(ctx, query, b.UserId).Scan(&b.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return b, error2.NotFoundErr
		}
		return b, err
	}

	return b, nil
}

func (r *databaseRepository) Create(ctx context.Context, balance models.UserBalance) error {
	return r.db.Exec(ctx, "INSERT INTO balance VALUES ($1, $2)", balance.UserId, balance.Value)
}

func (r *databaseRepository) Update(ctx context.Context, balance models.UserBalance) error {
	return r.db.Exec(ctx, "UPDATE balance SET value=$1 WHERE user_id=$2", balance.Value, balance.UserId)
}

func (r *databaseRepository) Transaction(ctx context.Context, f func(ctx context.Context) error) error {
	return r.db.WithTransaction(ctx, sql.LevelReadCommitted, f)
}
