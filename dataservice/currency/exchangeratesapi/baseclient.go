package exchangeratesapi

import (
	"context"
	"encoding/json"
	"errors"
	e "github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/pkg/utils"
	"net/http"
	"strings"
	"time"
)

const BaseURLV1 = "http://api.exchangeratesapi.io/v1"

type baseExchangesRatesClient struct {
	BaseURL, apiKey string
	HTTPClient      *http.Client
}

func NewBaseClient(apikey string, timeout int64) *baseExchangesRatesClient {
	return &baseExchangesRatesClient{
		BaseURL: BaseURLV1,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
		apiKey: apikey,
	}
}

// ratesResponse struct, containing the response of API endpoint
//
// The implemented schema is presented on https://exchangeratesapi.io/documentation/
type ratesResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

// GetRates fetches rates for provided currencies.
//
// If there is no provided currencies, method returns rates for all available currencies.
//
// If provided not supported currency, then its rate in returned map is equals to 0
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
		return nil, errors.New("API request has not succeeded")
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

	// Successful are considered only 2xx and 3xx status codes
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return e.APIBadStatusCode(resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&v); err != nil {
		return e.APIUnexpectedSchema(err)
	}
	return nil
}
