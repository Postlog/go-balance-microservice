package currency

import (
	"context"
	"github.com/postlog/go-balance-microservice/dataservice/currency/exchangeratesapi"
)

type Client interface {
	GetRates(ctx context.Context, currencies ...string) (map[string]float64, error)
}

// NewClient returns the client, that implements Client interface
// timeout - waiting time for a response from the API in milliseconds
func NewClient(apikey string, timeout int64) Client {
	return exchangeratesapi.NewClient(apikey, timeout)
}
