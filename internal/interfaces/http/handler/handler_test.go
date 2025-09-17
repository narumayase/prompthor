package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"prompthor/internal/domain"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUseCase is a mock implementation of ChatUseCase
type MockChatUseCase struct {
	mock.Mock
}

func (m *MockChatUseCase) ProcessChat(ctx context.Context, prompt domain.PromptRequest) (*domain.ChatResponse, error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).(*domain.ChatResponse), args.Error(1)
}

func TestNewChatHandler(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	assert.NotNil(t, handler)
	assert.IsType(t, &ChatHandler{}, handler)
	assert.Equal(t, mockUseCase, handler.usecase)
}

func TestChatHandler_HandleChat_Success(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Prepare test data
	request := domain.PromptRequest{
		Prompt: "Hello, how are you?",
	}
	expectedResponse := &domain.ChatResponse{
		Response: "I'm doing well, thank you!",
	}

	mockUseCase.On("ProcessChat", context.Background(), request).Return(expectedResponse, nil)

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
	assert.Equal(t, expectedResponse.Response, response.Response)

	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_UseCaseError(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	request := domain.PromptRequest{
		Prompt: "Test prompt",
	}
	expectedError := errors.New("API connection failed")

	mockUseCase.On("ProcessChat", context.Background(), request).Return((*domain.ChatResponse)(nil), expectedError)

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

func TestChatHandler_HandleChat_MissingContentType(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	request := domain.PromptRequest{Prompt: "Test prompt"}
	expectedResponse := &domain.ChatResponse{
		Response: "Response to test prompt",
	}
	mockUseCase.On("ProcessChat", context.Background(), request).Return(expectedResponse, nil)

	// Send request without proper JSON content type
	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	// Don't set Content-Type header - Gin will still process valid JSON

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin actually processes valid JSON even without Content-Type header
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Response, response.Response)

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

	request := domain.PromptRequest{
		Prompt: longPrompt,
	}
	expectedResponse := &domain.ChatResponse{
		Response: "Response to long prompt",
	}

	mockUseCase.On("ProcessChat", context.Background(), request).Return(expectedResponse, nil)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Response, response.Response)

	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_WithHeaders(t *testing.T) {
	mockUseCase := &MockChatUseCase{}
	handler := NewChatHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	request := domain.PromptRequest{
		Prompt: "Test with headers",
	}
	expectedResponse := &domain.ChatResponse{
		Response: "Response with headers",
	}

	mockUseCase.On("ProcessChat", context.Background(), request).Return(expectedResponse, nil)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", "test-correlation-id")
	req.Header.Set("X-Routing-ID", "test-routing-id")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Response, response.Response)

	mockUseCase.AssertExpectations(t)
}
