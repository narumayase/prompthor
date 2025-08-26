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

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Produce(event []byte) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepository) Close() {
	m.Called()
}

func TestNewChatUseCase(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := NewChatUseCase(mockChatRepo, mockEventRepo)

	assert.NotNil(t, useCase)
	assert.IsType(t, &ChatUseCaseImpl{}, useCase)
}

func TestChatUseCaseImpl_ProcessChat_Success(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockChatRepo, producerRepo: mockEventRepo}

	prompt := "Hello, how are you?"
	expectedResponse := "I'm doing well, thank you!"

	mockChatRepo.On("SendMessage", prompt).Return(expectedResponse, nil)
	mockEventRepo.On("Produce", mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockChatRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_Error(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockChatRepo, producerRepo: mockEventRepo}

	prompt := "Hello, how are you?"
	expectedError := errors.New("API connection failed")

	mockChatRepo.On("SendMessage", prompt).Return("", expectedError)

	result, err := useCase.ProcessChat(prompt)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, expectedError, err)
	mockChatRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_EmptyPrompt(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockChatRepo, producerRepo: mockEventRepo}

	prompt := ""
	expectedResponse := "Please provide a valid prompt"

	mockChatRepo.On("SendMessage", prompt).Return(expectedResponse, nil)
	mockEventRepo.On("Produce", mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockChatRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_LongPrompt(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockChatRepo, producerRepo: mockEventRepo}

	// Create a long prompt
	longPrompt := ""
	for i := 0; i < 1000; i++ {
		longPrompt += "This is a very long prompt. "
	}
	expectedResponse := "Response to long prompt"

	mockChatRepo.On("SendMessage", longPrompt).Return(expectedResponse, nil)
	mockEventRepo.On("Produce", mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(longPrompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockChatRepo.AssertExpectations(t)
}

func TestChatUseCaseImpl_ProcessChat_SpecialCharacters(t *testing.T) {
	mockChatRepo := &MockChatRepository{}
	mockEventRepo := &MockEventRepository{}
	useCase := &ChatUseCaseImpl{chatRepo: mockChatRepo, producerRepo: mockEventRepo}

	prompt := "Hello! @#$%^&*()_+ 你好 🚀"
	expectedResponse := "Response with special characters handled"

	mockChatRepo.On("SendMessage", prompt).Return(expectedResponse, nil)
	mockEventRepo.On("Produce", mock.Anything).Return(nil)

	result, err := useCase.ProcessChat(prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockChatRepo.AssertExpectations(t)
}
