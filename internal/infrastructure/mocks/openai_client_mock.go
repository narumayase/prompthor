package mocks

import (
	"context"

	"github.com/sashabaranov/go-openai"
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
