package repository

import (
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

// OpenAIRepository implements ChatRepository using OpenAI API
type OpenAIRepository struct {
	client client.OpenAIClient
}

// NewOpenAIRepository creates a new instance of the OpenAI repository
func NewOpenAIRepository(client client.OpenAIClient) (domain.ChatRepository, error) {
	return &OpenAIRepository{
		client: client,
	}, nil
}

// SendMessage sends a message to ChatGPT and returns the response
func (r *OpenAIRepository) SendMessage(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
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
