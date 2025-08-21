package infrastructure

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	// Test creating a new HTTP client
	httpClient := &http.Client{}
	bearerToken := "test-bearer-token"

	client := NewHttpClient(httpClient, bearerToken)

	assert.NotNil(t, client)
	
	// Verify it's the correct type
	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
}

func TestNewHttpClient_WithEmptyToken(t *testing.T) {
	// Test creating a client with empty token
	httpClient := &http.Client{}
	bearerToken := ""

	client := NewHttpClient(httpClient, bearerToken)

	assert.NotNil(t, client)
	
	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Empty(t, clientImpl.bearerToken)
}

func TestNewHttpClient_WithNilHttpClient(t *testing.T) {
	// Test creating a client with nil http.Client
	bearerToken := "test-token"

	client := NewHttpClient(nil, bearerToken)

	assert.NotNil(t, client)
	
	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Nil(t, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
}

func TestHttpClientImpl_Structure(t *testing.T) {
	// Test the structure of HttpClientImpl
	httpClient := &http.Client{}
	bearerToken := "test-structure-token"

	client := NewHttpClient(httpClient, bearerToken)
	clientImpl := client.(*HttpClientImpl)

	// Verify fields are accessible and correct
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
	assert.NotEmpty(t, clientImpl.bearerToken)
}
