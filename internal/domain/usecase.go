package domain

import "context"

// ChatUseCase defines the interface for the chat use case
type ChatUseCase interface {
	ProcessChat(ctx context.Context, prompt PromptRequest) (*ChatResponse, error)
}
