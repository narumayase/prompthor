package application

import (
	"anyprompt/internal/domain"
	"github.com/rs/zerolog/log"
)

// ChatUseCaseImpl implements ChatUseCase
type ChatUseCaseImpl struct {
	chatRepo domain.ChatRepository
}

// NewChatUseCase creates a new instance of the chat use case
func NewChatUseCase(chatRepo domain.ChatRepository) domain.ChatUseCase {
	return &ChatUseCaseImpl{
		chatRepo: chatRepo,
	}
}

// ProcessChat processes the chat request
func (uc *ChatUseCaseImpl) ProcessChat(prompt string) (string, error) {
	response, err := uc.chatRepo.SendMessage(prompt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return "", err
	}
	return response, nil
}
