package transaction

import (
	"context"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/dataservice/transaction/databaserepository"
	"github.com/postlog/go-balance-microservice/pkg/database"
)

type Repository interface {
	// Get returns an array of user's transactions
	Get(ctx context.Context, userId uuid.UUID, count, startFrom int, orderBy, orderDirection string) ([]models.Transaction, error)
	// Create creates new transaction record
	Create(ctx context.Context, t models.Transaction) error
}

func NewRepository(db database.Database) Repository {
	return databaserepository.New(db)
}
