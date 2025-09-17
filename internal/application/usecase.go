package application

import (
	"context"
	"github.com/rs/zerolog/log"
	"prompthor/internal/domain"
)

// ChatUseCaseImpl implements ChatUseCase
type ChatUseCaseImpl struct {
	chatRepository domain.LLMRepository
}

// NewChatUseCase creates a new instance of the chat use case
func NewChatUseCase(chatRepository domain.LLMRepository) domain.ChatUseCase {
	return &ChatUseCaseImpl{
		chatRepository: chatRepository,
	}
}

// ProcessChat processes the chat request
func (uc *ChatUseCaseImpl) ProcessChat(ctx context.Context, prompt domain.PromptRequest) (*domain.ChatResponse, error) {
	messageResponse, err := uc.chatRepository.Send(ctx, prompt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to send message")
		return nil, err
	}
	response := domain.ChatResponse{
		Response: messageResponse,
	}
	return &response, nil
}
