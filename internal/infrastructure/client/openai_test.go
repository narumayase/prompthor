package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAIClientImpl_CreateChatCompletion(t *testing.T) {
	// Test the actual OpenAI client implementation
	client := NewOpenAIClient("test-key")
	assert.NotNil(t, client)

	// We can't test the actual API call without a real key,
	// but we can verify the client is created properly
	clientImpl, ok := client.(*OpenAIClientImpl)
	assert.True(t, ok)
	assert.NotNil(t, clientImpl.client)
}
