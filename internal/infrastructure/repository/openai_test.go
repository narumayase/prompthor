package repository

import (
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAIClient is a mock implementation of OpenAIClient for testing
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

// Helper function to create mock OpenAI responses
func CreateMockOpenAIResponse(content string) openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		ID:      "chatcmpl-test",
		Object:  "chat.completion",
		Created: 1677652288,
		Model:   "gpt-3.5-turbo",
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: content,
				},
				FinishReason: "stop",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}
}

// Helper function to create empty OpenAI response (no choices)
func CreateEmptyOpenAIResponse() openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		ID:      "chatcmpl-test",
		Object:  "chat.completion",
		Created: 1677652288,
		Model:   "gpt-3.5-turbo",
		Choices: []openai.ChatCompletionChoice{}, // Empty choices
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 0,
			TotalTokens:      10,
		},
	}
}

func TestNewOpenAIRepository_Success(t *testing.T) {
	// Create mock client
	mockClient := &MockOpenAIClient{}

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
	mockClient := &MockOpenAIClient{}

	repo, err := NewOpenAIRepository(mockClient)
	assert.NoError(t, err)

	openAIRepo, ok := repo.(*OpenAIRepository)
	assert.True(t, ok)
	assert.NotNil(t, openAIRepo.client)
}

func TestOpenAIRepository_SendMessage_Success(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &MockOpenAIClient{}

	// Setup mock response
	mockResponse := CreateMockOpenAIResponse("Hello! How can I assist you today?")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	// Test the actual Send method with domain.PromptRequest
	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(context.Background(), promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Hello! How can I assist you today?", response)

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_APIError(t *testing.T) {
	// Create mock OpenAI client that returns error
	mockClient := &MockOpenAIClient{}

	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(openai.ChatCompletionResponse{}, errors.New("API connection failed"))

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(context.Background(), promptRequest)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "error calling OpenAI API")

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_EmptyResponse(t *testing.T) {
	// Create mock OpenAI client that returns empty choices
	mockClient := &MockOpenAIClient{}

	mockResponse := CreateEmptyOpenAIResponse()
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: "Hello world"}
	response, err := repo.Send(context.Background(), promptRequest)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "no response from OpenAI API")

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_EmptyPrompt(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &MockOpenAIClient{}

	// Setup mock response for empty prompt
	mockResponse := CreateMockOpenAIResponse("Please provide a prompt.")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	// Test with empty prompt
	promptRequest := domain.PromptRequest{Prompt: ""}
	response, err := repo.Send(context.Background(), promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Please provide a prompt.", response)

	mockClient.AssertExpectations(t)
}

func TestOpenAIRepository_SendMessage_LongPrompt(t *testing.T) {
	// Create mock OpenAI client
	mockClient := &MockOpenAIClient{}

	longPrompt := strings.Repeat("This is a very long prompt. ", 100)
	mockResponse := CreateMockOpenAIResponse("Response to long prompt")
	mockClient.On("CreateChatCompletion", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("openai.ChatCompletionRequest")).Return(mockResponse, nil)

	// Create repository with mock client
	repo, _ := NewOpenAIRepository(mockClient)

	promptRequest := domain.PromptRequest{Prompt: longPrompt}
	response, err := repo.Send(context.Background(), promptRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Response to long prompt", response)

	mockClient.AssertExpectations(t)
}
