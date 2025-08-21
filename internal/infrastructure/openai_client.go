package infrastructure

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient interface for dependency injection
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// OpenAIClientImpl wraps the standard OpenAI client
type OpenAIClientImpl struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) OpenAIClient {
	return &OpenAIClientImpl{
		client: openai.NewClient(apiKey),
	}
}

func (c *OpenAIClientImpl) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return c.client.CreateChatCompletion(ctx, request)
}
