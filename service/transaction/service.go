package transaction

import (
	"context"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/dataservice/transaction"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/pkg/utils"
	"time"
)

type Service interface {
	CreateTransaction(ctx context.Context, senderId, receiverId uuid.NullUUID, amount float64, description string) error
	GetTransactions(ctx context.Context, userId uuid.UUID, limit, offset int, orderBy, orderDir string) ([]models.Transaction, error)
}

// NewService returns the service, that implements Service interface
func NewService(repo transaction.Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo transaction.Repository
}

func (s *service) CreateTransaction(ctx context.Context, senderId, receiverId uuid.NullUUID, amount float64, description string) error {
	if !senderId.Valid && !receiverId.Valid {
		return errors.NewArgumentError("both users ids are invalid")
	}
	if amount <= 0 {
		return errors.NewArgumentError("amount cannot be less or equal to 0")
	}

	t := models.Transaction{
		SenderId: senderId, ReceiverId: receiverId, Amount: amount, Description: description, Date: time.Now().UTC(),
	}
	return s.repo.Create(ctx, t)
}

func (s *service) GetTransactions(
	ctx context.Context,
	userId uuid.UUID,
	count, startFrom int, orderBy, orderDirection string,
) ([]models.Transaction, error) {
	if count < 0 {
		return nil, errors.NewArgumentError("count cannot be less or equal to zero")
	}
	if startFrom < 0 {
		return nil, errors.NewArgumentError("start index cannot be less or equal to zero")
	}
	if orderBy != "" && !utils.StringInCollection(orderBy, "amount", "date") {
		return nil, errors.NewArgumentError("ordering can only be done by fields \"amount\" and \"date\"")
	}

	if orderDirection != "" && !utils.StringInCollection(orderDirection, "desc", "asc") {
		return nil, errors.NewArgumentError("order direction value can only be \"desc\" or \"asc\"")
	}

	return s.repo.Get(ctx, userId, count, startFrom, orderBy, orderDirection)
}
