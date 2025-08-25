package domain

// ChatRepository defines the interface for the chat repository
type ChatRepository interface {
	SendMessage(prompt string) (string, error)
}
