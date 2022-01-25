package balance

import (
	"context"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/internal/errors"
)

type Service interface {
	AddToBalance(ctx context.Context, userId uuid.UUID, amount float64) error
	ReduceBalance(ctx context.Context, userId uuid.UUID, amount float64) error
	GetBalance(ctx context.Context, userId uuid.UUID) (float64, error)
	TransferFounds(ctx context.Context, senderId, receiverId uuid.UUID, amount float64) error
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

var (
	AmountValueErr     = errors.NewArgumentError("amount cannot be less or equal to 0")
	NotEnoughFoundsErr = errors.NewArgumentError("user has not enough money")
	IdsAreEqualErr     = errors.NewArgumentError("sender id and receiver id are equal")
)

func (s *service) AddToBalance(ctx context.Context, userId uuid.UUID, amount float64) error {
	if amount <= 0 {
		return AmountValueErr
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		b, err := s.repo.GetWithSpec(ctx, UserBalance{UserId: userId}, ForUpdateSpec{})
		if err != nil {
			if err == NotFoundErr {
				return s.repo.Create(ctx, UserBalance{UserId: userId, Value: amount})
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
		b, err := s.repo.GetWithSpec(ctx, UserBalance{UserId: userId}, ForUpdateSpec{})
		if err != nil {
			if err == NotFoundErr {
				return errors.NewArgumentError(err.Error())
			}
			return err
		}

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
		senderBalance, err := s.repo.GetWithSpec(ctx, UserBalance{UserId: senderId}, ForUpdateSpec{})
		if err != nil {
			if err == NotFoundErr {
				return errors.NewArgumentError(err.Error())
			}
			return err
		}
		if senderBalance.Value < amount {
			return NotEnoughFoundsErr
		}

		receiverBalance, err := s.repo.GetWithSpec(ctx, UserBalance{UserId: receiverId}, ForUpdateSpec{})
		if err != nil {
			if err == NotFoundErr {
				receiverBalance = UserBalance{UserId: receiverId, Value: 0}
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

func (s *service) GetBalance(ctx context.Context, userId uuid.UUID) (float64, error) {
	balance, err := s.repo.Get(ctx, UserBalance{UserId: userId})
	if err != nil {
		if err == NotFoundErr {
			return 0, nil
		}
		return 0, err
	}

	return balance.Value, nil
}
