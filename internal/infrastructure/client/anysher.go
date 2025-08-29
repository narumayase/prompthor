package client

import (
	"context"
	"net/http"

	anysherhttp "github.com/narumayase/anysher/http"
)

// AnysherHTTPClient is a wrapper around anysherhttp.Client to implement the domain.HTTPClient interface
type AnysherHTTPClient struct {
	Client *anysherhttp.Client
}

// NewAnysherHTTPClient creates a new AnysherHTTPClient
func NewAnysherHTTPClient(client *anysherhttp.Client) *AnysherHTTPClient {
	return &AnysherHTTPClient{Client: client}
}

// Post sends a POST request using the anysherhttp.Client
func (c *AnysherHTTPClient) Post(url, contentType string, body []byte) (*http.Response, error) {
	payload := anysherhttp.Payload{
		URL:     url,
		Headers: map[string]string{"Content-Type": contentType},
		Content: body,
	}
	return c.Client.Post(context.Background(), payload)
}
