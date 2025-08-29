package http

import "net/http"

// HTTPClient defines the interface for an HTTP client.
type HTTPClient interface {
	Post(url, contentType string, body []byte) (*http.Response, error)
}
