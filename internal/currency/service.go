package currency

import (
	"context"
	"github.com/postlog/go-balance-microservice/internal/errors"
	"strings"
)

type Service interface {
	Convert(ctx context.Context, from, to string, amount float64) (float64, error)
}

type service struct {
	client ExchangeRatesClient
}

func NewService(client ExchangeRatesClient) Service {
	return &service{client}
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
