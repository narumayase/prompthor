package application

import (
	"anyompt/internal/domain"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

// ChatUseCaseImpl implements ChatUseCase
type ChatUseCaseImpl struct {
	chatRepository     domain.LLMRepository
	producerRepository domain.ProducerRepository
}

// NewChatUseCase creates a new instance of the chat use case
func NewChatUseCase(chatRepository domain.LLMRepository, producerRepository domain.ProducerRepository) domain.ChatUseCase {
	return &ChatUseCaseImpl{
		chatRepository:     chatRepository,
		producerRepository: producerRepository,
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
	if err := uc.produceMessage(ctx, response); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to produce message")
		return nil, err
	}
	return &response, nil
}

// produceMessage
func (uc *ChatUseCaseImpl) produceMessage(ctx context.Context, response domain.ChatResponse) error {
	if uc.producerRepository == nil {
		log.Ctx(ctx).Debug().Msg("No producer repository defined")
		return nil
	}
	message, err := json.Marshal(response)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal message")
		return err
	}
	return uc.producerRepository.Send(ctx, message)
}
