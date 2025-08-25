package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"anyprompt/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUseCase is a mock implementation of ChatUseCase
type MockChatUseCase struct {
	mock.Mock
}

func (m *MockChatUseCase) ProcessChat(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func TestNewChatHandler(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	assert.NotNil(t, handler)
	assert.IsType(t, &ChatHandler{}, handler)
	assert.Equal(t, mockUseCase, handler.chatUseCase)
}

func TestChatHandler_HandleChat_Success(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Prepare test data
	request := domain.ChatRequest{
		Prompt: "Hello, how are you?",
	}
	expectedResponse := "I'm doing well, thank you!"

	mockUseCase.On("ProcessChat", request.Prompt).Return(expectedResponse, nil)

	// Create request body
	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response.Response)
	assert.Empty(t, response.Error)

	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_UseCaseError(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	request := domain.ChatRequest{
		Prompt: "Test prompt",
	}
	expectedError := errors.New("API connection failed")

	mockUseCase.On("ProcessChat", request.Prompt).Return("", expectedError)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Error processing chat")
	assert.Contains(t, response["error"], expectedError.Error())

	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_InvalidJSON(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request format")
}

func TestChatHandler_HandleChat_EmptyPrompt(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	request := domain.ChatRequest{
		Prompt: "",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// This should fail validation due to the "required" tag
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request format")
}

func TestChatHandler_HandleChat_MissingContentType(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Mock the use case since Gin will still bind valid JSON even without Content-Type
	expectedResponse := "Response to test prompt"
	mockUseCase.On("ProcessChat", "Test prompt").Return(expectedResponse, nil)

	// Send request without proper JSON content type
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer([]byte(`{"prompt":"Test prompt"}`)))
	// Don't set Content-Type header - Gin will still process valid JSON

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin actually processes valid JSON even without Content-Type header
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response.Response)

	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_LongPrompt(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Create a long prompt
	longPrompt := ""
	for i := 0; i < 1000; i++ {
		longPrompt += "This is a very long prompt. "
	}

	request := domain.ChatRequest{
		Prompt: longPrompt,
	}
	expectedResponse := "Response to long prompt"

	mockUseCase.On("ProcessChat", longPrompt).Return(expectedResponse, nil)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response.Response)

	mockUseCase.AssertExpectations(t)
}
