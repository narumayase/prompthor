package repository

import (
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

// OpenAIRepository implements LLMRepository using OpenAI API
type OpenAIRepository struct {
	client client.OpenAIClient
}

// NewOpenAIRepository creates a new instance of the OpenAI repository
func NewOpenAIRepository(client client.OpenAIClient) (domain.LLMRepository, error) {
	return &OpenAIRepository{
		client: client,
	}, nil
}

// Send sends a message to ChatGPT and returns the response
func (r *OpenAIRepository) Send(ctx context.Context, prompt domain.PromptRequest) (string, error) {
	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt.Prompt,
				},
			},
		},
	)
	response := ""
	if err != nil {
		return response, fmt.Errorf("error calling OpenAI API: %w", err)
	}
	if len(resp.Choices) == 0 {
		return response, fmt.Errorf("no response from OpenAI API")
	}
	return resp.Choices[0].Message.Content, nil
}
