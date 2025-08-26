package domain

// ChatUseCase defines the interface for the chat use case
type ChatUseCase interface {
	ProcessChat(prompt string) (*ChatResponse, error)
}
