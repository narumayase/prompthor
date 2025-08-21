package infrastructure

import (
	"anyprompt/internal/infrastructure/client"
	"anyprompt/pkg/domain"
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIRepository implements ChatRepository using OpenAI API
type OpenAIRepository struct {
	client client.OpenAIClient
}

// NewOpenAIRepository creates a new instance of the OpenAI repository
func NewOpenAIRepository() (domain.ChatRepository, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return &OpenAIRepository{
		client: client.NewDefaultOpenAIClient(apiKey),
	}, nil
}

// NewOpenAIRepositoryWithClient creates a repository with custom OpenAI client (for testing)
func NewOpenAIRepositoryWithClient(client client.OpenAIClient) domain.ChatRepository {
	return &OpenAIRepository{
		client: client,
	}
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
