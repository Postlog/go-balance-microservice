package currency

import (
	"context"
	"encoding/json"
	"errors"
	e "github.com/postlog/go-balance-microservice/internal/errors"
	"github.com/postlog/go-balance-microservice/internal/utils"
	"net/http"
	"strings"
	"time"
)

const BaseURLV1 = "http://api.exchangeratesapi.io/v1"

type ExchangeRatesClient interface {
	GetRates(ctx context.Context, currencies ...string) (map[string]float64, error)
}

type baseExchangesRatesClient struct {
	BaseURL, apiKey string
	HTTPClient      *http.Client
}

func NewClient(apikey string, timeout int64) ExchangeRatesClient {
	return &baseExchangesRatesClient{
		BaseURL: BaseURLV1,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
		apiKey: apikey,
	}
}

type ratesResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

var (
	UnsupportedCurrencyErr = errors.New("specified currency not supported")
)

func (c *baseExchangesRatesClient) GetRates(ctx context.Context, currencies ...string) (map[string]float64, error) {
	req, err := utils.PrepareGETRequest(ctx, c.BaseURL, "latest",
		"access_key", c.apiKey,
		"symbols", strings.Join(currencies, ","),
	)

	if err != nil {
		return nil, err
	}

	var rates ratesResponse
	err = c.sendRequest(req, &rates)
	if err != nil {
		return nil, err
	}

	if !rates.Success {
		return nil, UnsupportedCurrencyErr
	}
	return rates.Rates, nil
}

func (c *baseExchangesRatesClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return e.APIRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return e.APIBadStatusCode(resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&v); err != nil {
		return e.APIUnexpectedSchema(err)
	}
	return nil
}
