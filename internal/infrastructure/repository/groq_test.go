package repository

import (
	"anyprompt/internal/infrastructure/client/mocks"
	"errors"
	"net/http"
	"testing"

	"anyprompt/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGroqRepository_Success(t *testing.T) {
	cfg := config.Config{
		GroqAPIKey: "test-groq-api-key",
		ChatModel:  "test-model",
	}

	// Create mock HTTP client
	mockClient := &mocks.MockHTTPClient{}
	repo, err := NewGroqRepository(cfg, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.IsType(t, &GroqRepository{}, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Equal(t, "test-groq-api-key", groqRepo.apiKey)
	assert.Equal(t, "test-model", groqRepo.model)
}

func TestNewGroqRepository_EmptyConfig(t *testing.T) {
	cfg := config.Config{}

	// Create mock HTTP client
	mockClient := &mocks.MockHTTPClient{}
	repo, err := NewGroqRepository(cfg, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Empty(t, groqRepo.apiKey)
	assert.Empty(t, groqRepo.model)
}

func TestNewGroqRepository_PartialConfig(t *testing.T) {
	cfg := config.Config{
		GroqAPIKey: "partial-key",
		ChatModel:  "",
	}

	// Create mock HTTP client
	mockClient := &mocks.MockHTTPClient{}
	repo, err := NewGroqRepository(cfg, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Equal(t, "partial-key", groqRepo.apiKey)
	assert.Empty(t, groqRepo.model)
}

func TestGroqRepository_Structure(t *testing.T) {
	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "groq-model",
	}

	// Create mock HTTP client
	mockClient := &mocks.MockHTTPClient{}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.NotEmpty(t, groqRepo.apiKey)
	assert.NotEmpty(t, groqRepo.model)
}

func TestGroqRepository_SendMessage_Success(t *testing.T) {
	// Create mock HTTP client
	mockClient := &mocks.MockHTTPClient{}

	// Setup mock response
	mockResponse := mocks.CreateMockResponse(200, `{
		"id": "test-id",
		"output": [
			{
				"type": "message",
				"id": "msg-1",
				"status": "completed",
				"content": [
					{
						"type": "output_text",
						"text": "Hello! How can I help you?"
					}
				]
			}
		]
	}`)

	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	// Create repository with mock client
	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	// Test the actual SendMessage method
	response, err := repo.SendMessage("Hello world")
	assert.NoError(t, err)
	assert.Equal(t, "Hello! How can I help you?", response)

	mockClient.AssertExpectations(t)
}

func TestGroqRepository_SendMessage_HTTPError(t *testing.T) {
	// Create mock HTTP client that returns HTTP error
	mockClient := &mocks.MockHTTPClient{}

	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return((*http.Response)(nil), errors.New("network error"))

	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	response, err := repo.SendMessage("Hello world")
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "network error")

	mockClient.AssertExpectations(t)
}

func TestGroqRepository_SendMessage_InvalidJSON(t *testing.T) {
	// Create mock HTTP client that returns invalid JSON
	mockClient := &mocks.MockHTTPClient{}

	mockResponse := mocks.CreateMockResponse(200, `invalid json`)
	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	response, err := repo.SendMessage("Hello world")
	assert.Error(t, err)
	assert.Empty(t, response)

	mockClient.AssertExpectations(t)
}

func TestGroqRepository_SendMessage_EmptyOutput(t *testing.T) {
	// Create mock HTTP client that returns empty output
	mockClient := &mocks.MockHTTPClient{}

	mockResponse := mocks.CreateMockResponse(200, `{
		"id": "test-id",
		"output": []
	}`)
	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	response, err := repo.SendMessage("Hello world")
	assert.NoError(t, err)
	assert.Empty(t, response)

	mockClient.AssertExpectations(t)
}

func TestGroqRepository_SendMessage_NonMessageType(t *testing.T) {
	// Create mock HTTP client that returns non-message type
	mockClient := &mocks.MockHTTPClient{}

	mockResponse := mocks.CreateMockResponse(200, `{
		"id": "test-id",
		"output": [
			{
				"type": "notification",
				"id": "notif-1",
				"status": "completed"
			}
		]
	}`)
	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	response, err := repo.SendMessage("Hello world")
	assert.NoError(t, err)
	assert.Empty(t, response)

	mockClient.AssertExpectations(t)
}

func TestGroqRepository_SendMessage_NonTextContent(t *testing.T) {
	// Create mock HTTP client that returns non-text content
	mockClient := &mocks.MockHTTPClient{}

	mockResponse := mocks.CreateMockResponse(200, `{
		"id": "test-id",
		"output": [
			{
				"type": "message",
				"id": "msg-1",
				"status": "completed",
				"content": [
					{
						"type": "image",
						"url": "http://example.com/image.jpg"
					}
				]
			}
		]
	}`)
	mockClient.On("Post", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
	}
	repo, err := NewGroqRepository(cfg, mockClient)
	assert.NoError(t, err)

	response, err := repo.SendMessage("Hello world")
	assert.NoError(t, err)
	assert.Empty(t, response)

	mockClient.AssertExpectations(t)
}


