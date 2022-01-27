package currency

import (
	"context"
	"github.com/postlog/go-balance-microservice/dataservice/currency"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"strings"
)

type Service interface {
	Convert(ctx context.Context, from, to string, amount float64) (float64, error)
}

type service struct {
	client currency.Client
}

// NewService returns the service, that implements Service interface
func NewService(client currency.Client) Service {
	return &service{client}
}

// Convert converts provided amount from one currency to another.
//
// If provided not supported currency,  then error returns
func (s *service) Convert(ctx context.Context, from, to string, amount float64) (float64, error) {
	from, to = strings.ToUpper(from), strings.ToUpper(to)
	if from == to {
		return amount, nil
	}

	rates, err := s.client.GetRates(ctx, from, to)
	if err != nil {
		return 0, err
	}

	if rates[from] == 0 || rates[to] == 0 {
		return 0, errors.NewArgumentError("specified currency not supported")
	}

	return (amount / rates[from]) * rates[to], nil
}
