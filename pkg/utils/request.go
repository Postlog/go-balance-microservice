package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func PrepareGETRequest(ctx context.Context, baseURL, resourcePath string, params ...string) (*http.Request, error) {
	if len(params)%2 != 0 {
		return nil, errors.New("params must be key value pairs")
	}

	parsedBaseURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	parsedPath, err := url.Parse(resourcePath)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s://%s", parsedBaseURL.Scheme, path.Join(parsedBaseURL.Host, parsedPath.Path))

	req, err := http.NewRequest("GET", reqURL, nil)
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
