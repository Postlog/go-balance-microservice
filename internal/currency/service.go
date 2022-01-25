package currency

import (
	"context"
	"github.com/postlog/go-balance-microservice/internal/errors"
	"github.com/postlog/go-balance-microservice/internal/logger"
	"strings"
)

type Service interface {
	Convert(ctx context.Context, from, to string, amount float64) (float64, error)
}

type service struct {
	client ExchangeRatesClient
	logger logger.Logger
}

func NewService(client ExchangeRatesClient, logger logger.Logger) Service {
	return &service{client, logger}
}

func (s *service) Convert(ctx context.Context, from, to string, amount float64) (float64, error) {
	from, to = strings.ToUpper(from), strings.ToUpper(to)

	rates, err := s.client.GetRates(ctx, from, to)
	if err != nil {
		if err == UnsupportedCurrencyErr {
			return 0, errors.NewArgumentError(err.Error())
		}
		return 0, err
	}

	return (amount / rates[from]) * rates[to], nil
}
