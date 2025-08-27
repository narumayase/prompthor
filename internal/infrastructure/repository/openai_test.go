package repository

import (
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"anyompt/internal/infrastructure/client/mocks"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewOpenAIRepository_Success(t *testing.T) {
	// Create mock client
	mockClient := &mocks.MockOpenAIClient{}

	repo, err := NewOpenAIRepository(mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.IsType(t, &OpenAIRepository{}, repo)
}

func TestNewOpenAIRepository_WithNilClient(t *testing.T) {
	// Test with nil client
	repo, err := NewOpenAIRepository(nil)

	assert.NoError(t, err) // Constructor doesn't validate nil client
	assert.NotNil(t, repo)
}

func TestNewOpenAIClientCreation(t *testing.T) {
	// Test creating a real OpenAI client
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	client := client.NewOpenAIClient("test-api-key")
	assert.NotNil(t, client)

	repo, err := NewOpenAIRepository(client)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestOpenAIRepository_Structure(t *testing.T) {
	// Create mock client
	mockClient := &mocks.MockOpenAIClient{}

	repo, err := NewOpenAIRepository(mockClient)
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
	repo, _ := NewOpenAIRepository(mockClient)

	// Test the actual Send method with domain.PromptRequest
	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Hello! How can I assist you today?", response)

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_APIError(t *testing.T) {
	// Create mock OpenAI client that returns error
	mockClient := &mocks.MockOpenAIClient{}

	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(openai.ChatCompletionResponse{}, errors.New("API connection failed"))

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(promptRequest)
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
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(promptRequest)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "no response from OpenAI API")

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_EmptyPrompt(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &mocks.MockOpenAIClient{}

	// Setup mock response for empty prompt
	mockResponse := mocks.CreateMockOpenAIResponse("Please provide a prompt.")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	// Test with empty prompt
	promptRequest := domain.PromptRequest{Prompt: ""}
	response, err := repo.Send(promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Please provide a prompt.", response)

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_LongPrompt(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &mocks.MockOpenAIClient{}

	longPrompt := strings.Repeat("This is a very long prompt. ", 100)
	mockResponse := mocks.CreateMockOpenAIResponse("Response to long prompt")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: longPrompt}
	response, err := repo.Send(promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Response to long prompt", response)

	mockClient.AssertExpectations(t)
}
