package domain

// ChatRepository defines the interface for the chat repository
type ChatRepository interface {
	SendMessage(prompt string) (string, error)
}

// ProducerRepository defines the interface for the producer repository for queue messages
type ProducerRepository interface {
	Produce(message []byte) error
	Close()
}
