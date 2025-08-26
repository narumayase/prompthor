package application

import (
	"anyompt/internal/domain"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

// ChatUseCaseImpl implements ChatUseCase
type ChatUseCaseImpl struct {
	chatRepository     domain.ChatRepository
	producerRepository domain.ProducerRepository
}

// NewChatUseCase creates a new instance of the chat use case
func NewChatUseCase(chatRepository domain.ChatRepository, producerRepository domain.ProducerRepository) domain.ChatUseCase {
	return &ChatUseCaseImpl{
		chatRepository:     chatRepository,
		producerRepository: producerRepository,
	}
}

// ProcessChat processes the chat request
func (uc *ChatUseCaseImpl) ProcessChat(prompt string) (*domain.ChatResponse, error) {
	responsePrompt, err := uc.chatRepository.SendMessage(prompt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return nil, err
	}
	response := &domain.ChatResponse{
		Response: responsePrompt,
	}
	if uc.producerRepository != nil {
		message, err := json.Marshal(*response)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal message")
			return nil, err
		}
		if err := uc.producerRepository.Produce(message); err != nil {
			log.Error().Err(err).Msg("Failed to send message to Kafka")
			return nil, err
		}
	}
	return response, nil
}
