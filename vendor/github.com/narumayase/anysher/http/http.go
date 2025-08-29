package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Payload struct {
	URL     string
	Token   string
	Headers map[string]string
	Content []byte
}

type Client struct {
	client *http.Client
	config Config
}

// NewClient creates a new HTTP client with bearer token authentication
func NewClient(client *http.Client, config Config) *Client {
	return &Client{
		client: client,
		config: config,
	}
}

// Post sends a POST request with JSON payload and bearer token authentication
func (c *Client) Post(ctx context.Context, payload Payload) (*http.Response, error) {
	payloadContent := payload.Content
	url := payload.URL
	headers := payload.Headers

	log.Debug().Msgf("payload to send: %s", string(payloadContent))
	log.Debug().Msgf("url %s", url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		// Set headers
		req.Header.Set(key, value)
	}
	log.Debug().Msgf("headers: to send to %s %+v", url, req.Header)

	// Set headers
	req.Header.Set("Authorization", "Bearer "+payload.Token)

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Err(err).Msgf("Failed to send message via HTTP: %v", resp)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	log.Debug().Msgf("API status code: %d", resp.StatusCode)
	return resp, nil
}
