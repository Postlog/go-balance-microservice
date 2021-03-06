package balance

import (
	"context"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/dataservice/balance"
	balanceError "github.com/postlog/go-balance-microservice/dataservice/balance/errors"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/service/constants"
	"math"
)

// Service interface defines methods, that allows to manipulate with user balance
type Service interface {
	AddToBalance(ctx context.Context, userId uuid.UUID, amount float64) error
	ReduceBalance(ctx context.Context, userId uuid.UUID, amount float64) error
	GetBalance(ctx context.Context, userId uuid.UUID) (models.UserBalance, error)
	TransferFounds(ctx context.Context, senderId, receiverId uuid.UUID, amount float64) error
}

// NewService returns the service, that implements Service interface
func NewService(repo balance.Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo balance.Repository
}

var (
	AmountValueErr     = errors.NewArgumentError("amount cannot be less or equal to 0")
	IdsAreEqualErr     = errors.NewArgumentError("sender id and receiver id are equal")
	NotEnoughFoundsErr = errors.NewServiceError(constants.NotEnoughFoundsErrCode, "user has not enough money")
)

func (s *service) AddToBalance(ctx context.Context, userId uuid.UUID, amount float64) error {
	if amount <= 0 {
		return AmountValueErr
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		b, err := s.repo.GetAndBlock(ctx, models.UserBalance{UserId: userId})
		if err != nil {
			if err == balanceError.NotFoundErr {
				return s.repo.Create(ctx, models.UserBalance{UserId: userId, Value: amount})
			}
			return err
		}

		b.Value += amount
		return s.repo.Update(ctx, b)
	})
}

func (s *service) ReduceBalance(ctx context.Context, userId uuid.UUID, amount float64) error {
	if amount <= 0 {
		return AmountValueErr
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		b, err := s.repo.GetAndBlock(ctx, models.UserBalance{UserId: userId})
		if err != nil {
			if err == balanceError.NotFoundErr {
				return NotEnoughFoundsErr
			}
			return err
		}

		amount = math.Round(amount*100)/100 // round to 2 decimal places
		if b.Value < amount {
			return NotEnoughFoundsErr
		}

		b.Value -= amount
		return s.repo.Update(ctx, b)
	})
}

func (s *service) TransferFounds(ctx context.Context, senderId, receiverId uuid.UUID, amount float64) error {
	if amount <= 0 {
		return AmountValueErr
	}

	if senderId == receiverId {
		return IdsAreEqualErr
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		senderBalance, err := s.repo.GetAndBlock(ctx, models.UserBalance{UserId: senderId})
		if err != nil {
			if err == balanceError.NotFoundErr {
				return NotEnoughFoundsErr
			}
			return err
		}

		amount = math.Round(amount*100)/100 // round to 2 decimal places
		if senderBalance.Value < amount {
			return NotEnoughFoundsErr
		}

		receiverBalance, err := s.repo.GetAndBlock(ctx, models.UserBalance{UserId: receiverId})
		if err != nil {
			if err == balanceError.NotFoundErr {
				receiverBalance = models.UserBalance{UserId: receiverId, Value: 0}
				err = s.repo.Create(ctx, receiverBalance)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		senderBalance.Value -= amount
		receiverBalance.Value += amount

		err = s.repo.Update(ctx, senderBalance)
		if err != nil {
			return err
		}
		return s.repo.Update(ctx, receiverBalance)
	})
}

func (s *service) GetBalance(ctx context.Context, userId uuid.UUID) (models.UserBalance, error) {
	b := models.UserBalance{UserId: userId}
	b, err := s.repo.Get(ctx, b)
	if err != nil {
		if err == balanceError.NotFoundErr {
			b.Value = 0
			return b, nil
		}
		return b, err
	}

	return b, nil
}
