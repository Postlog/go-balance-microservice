package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func PrepareGETRequest(ctx context.Context, baseURL, path string, params ...string) (*http.Request, error) {
	if len(params)%2 != 0 {
		return nil, errors.New("params must be key value pairs")
	}

	reqUrl := fmt.Sprintf("%s/%s", baseURL, path)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for i := 0; i < len(params); i += 2 {
		q.Add(params[i], params[i+1])
	}
	req.URL.RawQuery = q.Encode()
	return req.WithContext(ctx), nil
}
