package domain

import "context"

// LLMRepository defines the interface for the llm repository
type LLMRepository interface {
	Send(ctx context.Context, prompt PromptRequest) (string, error)
}
