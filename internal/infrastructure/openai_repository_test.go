package infrastructure

import (
	"anyprompt/internal/infrastructure/mocks"
	"errors"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewOpenAIRepository_Success(t *testing.T) {
	// Set test API key
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	repo, err := NewOpenAIRepository()

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.IsType(t, &OpenAIRepository{}, repo)
}

func TestNewOpenAIRepository_MissingAPIKey(t *testing.T) {
	// Ensure API key is not set
	os.Unsetenv("OPENAI_API_KEY")

	repo, err := NewOpenAIRepository()

	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "OPENAI_API_KEY environment variable is required")
}

func TestNewOpenAIRepository_EmptyAPIKey(t *testing.T) {
	// Set empty API key
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Unsetenv("OPENAI_API_KEY")

	repo, err := NewOpenAIRepository()

	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "OPENAI_API_KEY environment variable is required")
}

func TestOpenAIRepository_Structure(t *testing.T) {
	// Set test API key
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	repo, err := NewOpenAIRepository()
	assert.NoError(t, err)

	openAIRepo, ok := repo.(*OpenAIRepository)
	assert.True(t, ok)
	assert.NotNil(t, openAIRepo.client)
}

func TestOpenAIRepository_SendMessage_Success(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &mocks.MockOpenAIClient{}
	
	// Setup mock response
	mockResponse := mocks.CreateMockOpenAIResponse("Hello! How can I assist you today?")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo := NewOpenAIRepositoryWithClient(mockClient)

	// Test the actual SendMessage method
	response, err := repo.SendMessage("Hello world")
	assert.NoError(t, err)
	assert.Equal(t, "Hello! How can I assist you today?", response)

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_APIError(t *testing.T) {
	// Create mock OpenAI client that returns error
	mockClient := &mocks.MockOpenAIClient{}
	
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(openai.ChatCompletionResponse{}, errors.New("API connection failed"))

	// Create repository with mock client
	repo := NewOpenAIRepositoryWithClient(mockClient)

	response, err := repo.SendMessage("Hello world")
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "error calling OpenAI API")

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_EmptyResponse(t *testing.T) {
	// Create mock OpenAI client that returns empty choices
	mockClient := &mocks.MockOpenAIClient{}
	
	mockResponse := mocks.CreateEmptyOpenAIResponse()
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo := NewOpenAIRepositoryWithClient(mockClient)

	response, err := repo.SendMessage("Hello world")
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "no response from OpenAI API")

	mockClient.AssertExpectations(t)
}

