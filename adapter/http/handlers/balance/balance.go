package balance

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/adapter/http/handlers"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
	"github.com/postlog/go-balance-microservice/pkg/types"
	"github.com/postlog/go-balance-microservice/pkg/utils"
	"github.com/postlog/go-balance-microservice/service/balance"
	"github.com/postlog/go-balance-microservice/service/currency"
	"github.com/postlog/go-balance-microservice/service/transaction"
)

func Register(
	router fiber.Router,
	balanceService balance.Service,
	transactionService transaction.Service,
	currencyService currency.Service,
	transaction types.TransactionWrapper,
	baseCurrency string,
) {
	api := provider{balanceService, transactionService, currencyService, transaction, baseCurrency}

	router.Post("add", api.addToBalance)
	router.Post("reduce", api.reduceBalance)
	router.Post("transfer", api.transferFounds)
	router.Post("get", api.getBalance)
}

type provider struct {
	balanceService     balance.Service
	transactionService transaction.Service
	currencyService    currency.Service
	transaction        types.TransactionWrapper
	baseCurrency       string
}

func (p *provider) addToBalance(c *fiber.Ctx) error {
	req := protocol.UpdateBalanceRequest{}
	err := handlers.ParseBody(c, &req)
	if err != nil {
		return err
	}

	userId, amount := *req.UserId, *req.Amount

	amount, err = p.convertCurrency(c.Context(), req.Currency, amount)
	if err != nil {
		return err
	}

	return p.transaction(c.Context(), func(ctx context.Context) error {
		err = p.balanceService.AddToBalance(ctx, userId, amount)
		if err != nil {
			return err
		}

		err = p.transactionService.CreateTransaction(ctx, uuid.NullUUID{}, utils.ToNullableUUID(userId), amount, req.Description)
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *provider) reduceBalance(c *fiber.Ctx) error {
	req := protocol.UpdateBalanceRequest{}
	err := handlers.ParseBody(c, &req)
	if err != nil {
		return err
	}

	userId, amount := *req.UserId, *req.Amount

	amount, err = p.convertCurrency(c.Context(), req.Currency, amount)
	if err != nil {
		return err
	}

	return p.transaction(c.Context(), func(ctx context.Context) error {
		err = p.balanceService.ReduceBalance(ctx, userId, amount)
		if err != nil {
			return err
		}

		err = p.transactionService.CreateTransaction(ctx, utils.ToNullableUUID(userId), uuid.NullUUID{}, amount, req.Description)
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *provider) transferFounds(c *fiber.Ctx) error {
	req := protocol.TransferFoundsRequest{}
	err := handlers.ParseBody(c, &req)
	if err != nil {
		return err
	}

	senderId, receiverId, amount := *req.SenderId, *req.ReceiverId, *req.Amount

	amount, err = p.convertCurrency(c.Context(), req.Currency, amount)
	if err != nil {
		return err
	}

	return p.transaction(c.Context(), func(ctx context.Context) error {
		err = p.balanceService.TransferFounds(ctx, senderId, receiverId, amount)
		if err != nil {
			return err
		}

		err = p.transactionService.CreateTransaction(ctx, utils.ToNullableUUID(senderId), utils.ToNullableUUID(receiverId), amount, req.Description)
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *provider) getBalance(c *fiber.Ctx) error {
	req := protocol.GetBalanceRequest{}
	err := handlers.ParseBody(c, &req)
	if err != nil {
		return err
	}

	userId := *req.UserId

	b, err := p.balanceService.GetBalance(c.Context(), userId)
	if err != nil {
		return err
	}
	balancePayload := protocol.BalancePayload{Balance: b}
	bytes, err := json.Marshal(&balancePayload)
	if err != nil {
		return err
	}

	resp := protocol.Response{Error: nil, Payload: bytes}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (p *provider) convertCurrency(ctx context.Context, from string, amount float64) (float64, error) {
	if from == "" || from == p.baseCurrency {
		return amount, nil
	}

	converted, err := p.currencyService.Convert(ctx, from, p.baseCurrency, amount)
	if err != nil {
		return 0, err
	}
	return converted, nil
}
