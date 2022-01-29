package currency

import (
	"context"
	"github.com/postlog/go-balance-microservice/dataservice/currency"
	currencyErrors "github.com/postlog/go-balance-microservice/dataservice/currency/errors"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/service/constants"
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
		if err == currencyErrors.RateNotFoundErr {
			return 0, errors.NewServiceError(constants.CurrencyNotSupportedErrCode, "specified currency not supported")
		}
		return 0, err
	}

	return (amount / rates[from]) * rates[to], nil
}
