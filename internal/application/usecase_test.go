package application

import (
	"anyompt/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLLMRepository is a mock implementation of LLMRepository
type MockLLMRepository struct {
	mock.Mock
}

func (m *MockLLMRepository) Send(prompt domain.PromptRequest) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

// MockProducerRepository is a mock implementation of ProducerRepository
type MockProducerRepository struct {
	mock.Mock
}

func (m *MockProducerRepository) Close() {
}

func (m *MockProducerRepository) Produce(ctx context.Context, message []byte) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func TestNewChatUseCase(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := NewChatUseCase(mockChatRepo, mockProducerRepo)

	assert.NotNil(t, useCase)
	assert.IsType(t, &ChatUseCaseImpl{}, useCase)
}

func TestChatUseCaseImpl_ProcessChat_Success(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: mockProducerRepo,
	}

	promptRequest := domain.PromptRequest{Prompt: "Hello, how are you?"}
	expectedResponse := "I'm doing well, thank you!"

	mockChatRepo.On("Send", promptRequest).Return(expectedResponse, nil)
	// Fix: Use mock.Anything for context instead of specific type
	mockProducerRepo.On("Produce", mock.Anything, mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result.MessageResponse)
	mockChatRepo.AssertExpectations(t)
	mockProducerRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_Error(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: mockProducerRepo,
	}

	promptRequest := domain.PromptRequest{Prompt: "Hello, how are you?"}
	expectedError := errors.New("API connection failed")

	mockChatRepo.On("Send", promptRequest).Return("", expectedError)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockChatRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_EmptyPrompt(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: mockProducerRepo,
	}

	promptRequest := domain.PromptRequest{Prompt: ""}
	expectedResponse := "Please provide a valid prompt"

	mockChatRepo.On("Send", promptRequest).Return(expectedResponse, nil)
	mockProducerRepo.On("Produce", mock.Anything, mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result.MessageResponse)
	mockChatRepo.AssertExpectations(t)
	mockProducerRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_LongPrompt(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: mockProducerRepo,
	}

	// Create a long prompt
	longPrompt := ""
	for i := 0; i < 1000; i++ {
		longPrompt += "This is a very long prompt. "
	}
	promptRequest := domain.PromptRequest{Prompt: longPrompt}
	expectedResponse := "MessageResponse to long prompt"

	mockChatRepo.On("Send", promptRequest).Return(expectedResponse, nil)
	mockProducerRepo.On("Produce", mock.Anything, mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result.MessageResponse)
	mockChatRepo.AssertExpectations(t)
	mockProducerRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_SpecialCharacters(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	mockProducerRepo := &MockProducerRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: mockProducerRepo,
	}

	promptRequest := domain.PromptRequest{Prompt: "Hello! @#$%^&*()_+ Ã¤Â½ Ã¥Â¥Â½ Ã°Å¸Å¡â‚¬"}
	expectedResponse := "MessageResponse with special characters handled"

	mockChatRepo.On("Send", promptRequest).Return(expectedResponse, nil)
	mockProducerRepo.On("Produce", mock.Anything, mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result.MessageResponse)
	mockChatRepo.AssertExpectations(t)
	mockProducerRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_WithNilProducer(t *testing.T) {
	mockChatRepo := &MockLLMRepository{}
	useCase := &ChatUseCaseImpl{
		chatRepository:     mockChatRepo,
		producerRepository: nil, // nil producer
	}

	promptRequest := domain.PromptRequest{Prompt: "Test with nil producer"}
	expectedResponse := "Test response"

	mockChatRepo.On("Send", promptRequest).Return(expectedResponse, nil)

	result, err := useCase.ProcessChat(context.Background(), promptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result.MessageResponse)
	mockChatRepo.AssertExpectations(t)
}
