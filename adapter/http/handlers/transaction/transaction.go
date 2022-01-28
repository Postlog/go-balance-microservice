package transaction

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/handlers"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
	"github.com/postlog/go-balance-microservice/service/transaction"
)

func Register(router fiber.Router, transactionService transaction.Service) {
	api := provider{transactionService}
	router.Post("list", api.getTransactions)
}

type provider struct {
	transactionService transaction.Service
}

func (api *provider) getTransactions(c *fiber.Ctx) error {
	req := protocol.GetTransactionsRequest{}
	err := handlers.ParseBody(c, &req)
	if err != nil {
		return err
	}

	userId := *req.UserId

	ts, err := api.transactionService.GetTransactions(c.Context(), userId, req.Limit, req.Offset, req.OrderBy, req.OrderDirection)
	if err != nil {
		return err
	}

	transactionsPayload := protocol.TransactionsPayload{Transactions: ts}
	bytes, err := json.Marshal(&transactionsPayload)
	if err != nil {
		return err
	}

	resp := protocol.Response{Error: nil, Payload: bytes}
	return c.Status(fiber.StatusOK).JSON(resp)
}
