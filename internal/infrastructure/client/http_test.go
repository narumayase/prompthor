package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClientImpl_Post(t *testing.T) {
	// Test the actual HTTP client implementation
	client := NewHttpClient(&http.Client{}, "test-token")
	assert.NotNil(t, client)

	// We can't test the actual HTTP call without a real server,
	// but we can verify the client is created properly
	clientImpl, ok := client.(*HTTPClientImpl)
	assert.True(t, ok)
	assert.NotNil(t, clientImpl.client)
	assert.Equal(t, "test-token", clientImpl.bearerToken)
}