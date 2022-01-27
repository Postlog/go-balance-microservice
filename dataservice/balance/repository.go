package balance

import (
	"context"
	"github.com/postlog/go-balance-microservice/dataservice/balance/databaserepository"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/pkg/database"
)

// Repository interface refers to Repository pattern and encapsulates logic
// of access to user balance in storage
type Repository interface {
	// Get returns models.UserBalance for specified user
	Get(ctx context.Context, balance models.UserBalance) (models.UserBalance, error)
	// GetForUpdate returns models.UserBalance for specified user with lock for update
	//
	// Makes sense only if called inside a transaction
	GetAndBlock(ctx context.Context, balance models.UserBalance) (models.UserBalance, error)
	// Create creates new balance record in storage
	Create(ctx context.Context, balance models.UserBalance) error
	// Update updates existing balance record in storage
	Update(ctx context.Context, balance models.UserBalance) error
	// Transaction allows performing operations with the repository within a transaction
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

// NewRepository returns repository, that implements Repository interface
func NewRepository(db database.Database) Repository {
	return databaserepository.New(db)
}
