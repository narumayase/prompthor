package domain

import "context"

// LLMRepository defines the interface for the llm repository
type LLMRepository interface {
	Send(ctx context.Context, prompt PromptRequest) (string, error)
}

// ProducerRepository defines the interface for the producer repository for queue messages
type ProducerRepository interface {
	Produce(ctx context.Context, message []byte) error
}
