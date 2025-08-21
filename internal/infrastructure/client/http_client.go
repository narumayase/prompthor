package client

import "net/http"

// HTTPClient interface for dependency injection
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient wraps the standard http.Client
type DefaultHTTPClient struct {
	client *http.Client
}

func NewDefaultHTTPClient() HTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{},
	}
}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
