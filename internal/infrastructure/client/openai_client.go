package client

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient interface for dependency injection
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// DefaultOpenAIClient wraps the standard OpenAI client
type DefaultOpenAIClient struct {
	client *openai.Client
}

func NewDefaultOpenAIClient(apiKey string) OpenAIClient {
	return &DefaultOpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

func (c *DefaultOpenAIClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return c.client.CreateChatCompletion(ctx, request)
}
