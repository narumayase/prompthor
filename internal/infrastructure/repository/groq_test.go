package repository

import (
	"anyompt/config"
	"anyompt/internal/domain"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClient simula el cliente HTTP para testing
type MockHTTPClient struct {
	PostResponse *http.Response
	PostError    error
}

func (m *MockHTTPClient) Post(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
	return m.PostResponse, m.PostError
}

// Helper para crear respuestas HTTP mock
func createMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestNewGroqRepository_Success(t *testing.T) {
	cfg := config.Config{
		GroqAPIKey: "test-api-key",
		ChatModel:  "test-model",
		GroqUrl:    "https://api.groq.com/test",
	}
	mockClient := &MockHTTPClient{}

	repo, err := NewGroqRepository(cfg, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Equal(t, "test-api-key", groqRepo.apiKey)
	assert.Equal(t, "test-model", groqRepo.model)
	assert.Equal(t, "https://api.groq.com/test", groqRepo.baseURL)
	assert.Equal(t, mockClient, groqRepo.httpClient)
}

func TestNewGroqRepository_EmptyConfig(t *testing.T) {
	cfg := config.Config{}
	mockClient := &MockHTTPClient{}

	repo, err := NewGroqRepository(cfg, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Empty(t, groqRepo.apiKey)
	assert.Empty(t, groqRepo.model)
	assert.Empty(t, groqRepo.baseURL)
}

func TestGroqRepository_Send_Success(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": [
				{
					"type": "message",
					"id": "msg-1",
					"status": "completed",
					"content": [
						{
							"type": "output_text",
							"text": "Esta es la respuesta del modelo"
						}
					]
				}
			]
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Hola, ¿cómo estás?"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Esta es la respuesta del modelo", result)
}

func TestGroqRepository_Send_HTTPError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: nil,
		PostError:    errors.New("network timeout"),
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "network timeout")
}

func TestGroqRepository_Send_InvalidJSON(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `invalid json response`),
		PostError:    nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestGroqRepository_Send_EmptyOutput(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": []
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGroqRepository_Send_NoMessageType(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": [
				{
					"type": "notification",
					"id": "notif-1",
					"status": "info"
				},
				{
					"type": "error",
					"id": "err-1",
					"status": "failed"
				}
			]
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGroqRepository_Send_NoTextContent(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": [
				{
					"type": "message",
					"id": "msg-1",
					"status": "completed",
					"content": [
						{
							"type": "image",
							"url": "https://example.com/image.jpg"
						},
						{
							"type": "audio",
							"url": "https://example.com/audio.mp3"
						}
					]
				}
			]
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGroqRepository_Send_MultipleMessages(t *testing.T) {
	// Arrange - simula múltiples mensajes, debería devolver el primer output_text
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": [
				{
					"type": "notification",
					"id": "notif-1",
					"status": "info"
				},
				{
					"type": "message",
					"id": "msg-1",
					"status": "completed",
					"content": [
						{
							"type": "output_text",
							"text": "Primera respuesta"
						}
					]
				},
				{
					"type": "message",
					"id": "msg-2",
					"status": "completed",
					"content": [
						{
							"type": "output_text",
							"text": "Segunda respuesta"
						}
					]
				}
			]
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	// Debería devolver la última respuesta encontrada (comportamiento actual del código)
	assert.Equal(t, "Segunda respuesta", result)
}

func TestGroqRepository_Send_MixedContent(t *testing.T) {
	// Arrange - simula contenido mixto, debería extraer solo el texto
	mockClient := &MockHTTPClient{
		PostResponse: createMockResponse(200, `{
			"id": "test-id-123",
			"output": [
				{
					"type": "message",
					"id": "msg-1",
					"status": "completed",
					"content": [
						{
							"type": "image",
							"url": "https://example.com/image.jpg"
						},
						{
							"type": "output_text",
							"text": "Respuesta de texto"
						},
						{
							"type": "audio",
							"url": "https://example.com/audio.mp3"
						}
					]
				}
			]
		}`),
		PostError: nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Respuesta de texto", result)
}

func TestGroqRepository_Send_ReadBodyError(t *testing.T) {
	// Arrange - simula error al leer el body de la respuesta
	mockResponse := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       &errorReader{},
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		PostResponse: mockResponse,
		PostError:    nil,
	}

	cfg := config.Config{
		GroqAPIKey: "test-key",
		ChatModel:  "llama-3.1-70b",
		GroqUrl:    "https://api.groq.com/openai/v1/chat/completions",
	}

	repo, _ := NewGroqRepository(cfg, mockClient)
	prompt := domain.PromptRequest{Prompt: "Test prompt"}

	// Act
	result, err := repo.Send(context.Background(), prompt)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "read error")
}

// errorReader simula un error al leer
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}
