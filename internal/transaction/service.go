package transaction

import (
	"context"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/internal/errors"
	"time"
)

type Service interface {
	CreateTransaction(ctx context.Context, senderId, receiverId uuid.NullUUID, amount float64, description string) error
	GetTransactions(ctx context.Context, userId uuid.UUID, limit, offset int, orderBy, orderDir string) ([]Transaction, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

var (
	BothIdsInvalidErr       = errors.NewArgumentError("both users ids are invalid")
	EmptyDescriptionBodyErr = errors.NewArgumentError("description body is empty")
	AmountValueErr          = errors.NewArgumentError("amount cannot be less or equal to 0")
)

type service struct {
	repo Repository
}

func (s *service) CreateTransaction(ctx context.Context, senderId, receiverId uuid.NullUUID, amount float64, description string) error {
	if !senderId.Valid && !receiverId.Valid {
		return BothIdsInvalidErr
	}
	if amount <= 0 {
		return AmountValueErr
	}
	if description == "" {
		return EmptyDescriptionBodyErr
	}

	t := Transaction{
		SenderId: senderId, ReceiverId: receiverId, Amount: amount, Description: description, Date: time.Now().UTC(),
	}
	return s.repo.Create(ctx, t)
}

func (s *service) GetTransactions(ctx context.Context, userId uuid.UUID, limit, offset int, orderBy, orderDir string) ([]Transaction, error) {
	opts, err := NewPaginationOptions(limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, errors.NewArgumentError(err.Error())
	}

	return s.repo.Get(ctx, userId, opts)
}
