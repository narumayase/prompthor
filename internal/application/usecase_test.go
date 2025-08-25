package application

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatRepository is a mock implementation of ChatRepository
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) SendMessage(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func TestNewChatUseCase(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := NewChatUseCase(mockRepo)

	assert.NotNil(t, useCase)
	assert.IsType(t, &ChatUseCaseImpl{}, useCase)
}

func TestChatUseCaseImpl_ProcessChat_Success(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockRepo}

	prompt := "Hello, how are you?"
	expectedResponse := "I'm doing well, thank you!"

	mockRepo.On("SendMessage", prompt).Return(expectedResponse, nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_Error(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockRepo}

	prompt := "Hello, how are you?"
	expectedError := errors.New("API connection failed")

	mockRepo.On("SendMessage", prompt).Return("", expectedError)

	result, err := useCase.ProcessChat(prompt)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_EmptyPrompt(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockRepo}

	prompt := ""
	expectedResponse := "Please provide a valid prompt"

	mockRepo.On("SendMessage", prompt).Return(expectedResponse, nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_LongPrompt(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockRepo}

	// Create a long prompt
	longPrompt := ""
	for i := 0; i < 1000; i++ {
		longPrompt += "This is a very long prompt. "
	}
	expectedResponse := "Response to long prompt"

	mockRepo.On("SendMessage", longPrompt).Return(expectedResponse, nil)

	result, err := useCase.ProcessChat(longPrompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_SpecialCharacters(t *testing.T) {
	mockRepo := &MockChatRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockRepo}

	prompt := "Hello! @#$%^&*()_+ 你好 🚀"
	expectedResponse := "Response with special characters handled"

	mockRepo.On("SendMessage", prompt).Return(expectedResponse, nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockRepo.AssertExpectations(t)
}
