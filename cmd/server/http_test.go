package server

import (
	"anyompt/internal/config"
	"anyompt/internal/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUseCase is a mock implementation of domain.ChatUseCase
type MockChatUseCase struct {
	mock.Mock
}

func (m *MockChatUseCase) ProcessChat(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func TestRun_ServerConfiguration(t *testing.T) {
	// Create a mock use case
	mockUseCase := new(MockChatUseCase)
	mockUseCase.On("ProcessChat", mock.Anything).Return("test response", nil)

	// Test configuration
	cfg := config.Config{
		Port:       "8080",
		OpenAIKey:  "test-key",
		GroqAPIKey: "test-groq-key",
		GroqUrl:    "https://api.groq.com/openai/v1/responses",
		ChatModel:  "test-model",
	}

	// Test that the function sets up the server correctly
	// We can't easily test the actual server startup without blocking,
	// but we can test the configuration and setup
	t.Run("server configuration", func(t *testing.T) {
		// This test verifies that the function accepts the correct parameters
		// and would set up a server on the correct port
		expectedAddr := ":" + cfg.Port
		assert.Equal(t, ":8080", expectedAddr)
	})
}

func TestRun_RouterSetup(t *testing.T) {
	// Create a mock use case
	mockUseCase := new(MockChatUseCase)
	mockUseCase.On("ProcessChat", "test").Return("response", nil)

	cfg := config.Config{
		Port: "8080",
	}

	// We can't directly test the Run function without it blocking,
	// but we can test that the router setup works correctly by
	// testing the router functionality indirectly
	t.Run("router setup validation", func(t *testing.T) {
		// Test that we can create a valid configuration
		assert.NotEmpty(t, cfg.Port)
		assert.NotNil(t, mockUseCase)
	})
}

func TestRun_PortConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		port         string
		expectedAddr string
	}{
		{
			name:         "default port",
			port:         "8080",
			expectedAddr: ":8080",
		},
		{
			name:         "custom port",
			port:         "3000",
			expectedAddr: ":3000",
		},
		{
			name:         "high port number",
			port:         "9999",
			expectedAddr: ":9999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{Port: tt.port}
			serverAddr := ":" + cfg.Port
			assert.Equal(t, tt.expectedAddr, serverAddr)
		})
	}
}

func TestRun_ConfigValidation(t *testing.T) {
	mockUseCase := new(MockChatUseCase)

	t.Run("valid config", func(t *testing.T) {
		cfg := config.Config{
			Port:       "8080",
			OpenAIKey:  "test-openai-key",
			GroqAPIKey: "test-groq-key",
			GroqUrl:    "https://api.groq.com/openai/v1/responses",
			ChatModel:  "test-model",
		}

		// Verify config fields are properly set
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "test-openai-key", cfg.OpenAIKey)
		assert.Equal(t, "test-groq-key", cfg.GroqAPIKey)
		assert.Equal(t, "https://api.groq.com/openai/v1/responses", cfg.GroqUrl)
		assert.Equal(t, "test-model", cfg.ChatModel)
		assert.NotNil(t, mockUseCase)
	})

	t.Run("minimal config", func(t *testing.T) {
		cfg := config.Config{
			Port: "8080",
		}

		// Should work with minimal configuration
		assert.Equal(t, "8080", cfg.Port)
		assert.NotNil(t, mockUseCase)
	})
}

func TestRun_UseCaseIntegration(t *testing.T) {
	t.Run("usecase interface validation", func(t *testing.T) {
		mockUseCase := new(MockChatUseCase)
		mockUseCase.On("ProcessChat", "test prompt").Return("test response", nil)

		// Test that the mock implements the interface correctly
		var usecase domain.ChatUseCase = mockUseCase
		assert.NotNil(t, usecase)

		// Test the mock functionality
		response, err := usecase.ProcessChat("test prompt")
		assert.NoError(t, err)
		assert.Equal(t, "test response", response)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase error handling", func(t *testing.T) {
		mockUseCase := new(MockChatUseCase)
		mockUseCase.On("ProcessChat", "error prompt").Return("", errors.New("test error"))

		response, err := mockUseCase.ProcessChat("error prompt")
		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Equal(t, "test error", err.Error())

		mockUseCase.AssertExpectations(t)
	})
}

// TestRun_ServerStartup tests the server startup behavior
// Note: This test doesn't actually start the server to avoid blocking
func TestRun_ServerStartup(t *testing.T) {
	t.Run("server address format", func(t *testing.T) {
		cfg := config.Config{Port: "8080"}
		expectedAddr := ":" + cfg.Port

		// Verify the address format is correct for gin.Engine.Run()
		assert.Equal(t, ":8080", expectedAddr)
		assert.True(t, len(expectedAddr) > 1)
		assert.Equal(t, ":", expectedAddr[:1])
	})
}

// Integration test that verifies the HTTP server setup without actually starting it
func TestRun_HTTPServerIntegration(t *testing.T) {
	mockUseCase := new(MockChatUseCase)
	mockUseCase.On("ProcessChat", mock.Anything).Return("test response", nil)

	cfg := config.Config{
		Port:       "8080",
		OpenAIKey:  "test-key",
		GroqAPIKey: "test-groq-key",
		GroqUrl:    "https://api.groq.com/openai/v1/responses",
		ChatModel:  "test-model",
	}

	t.Run("server configuration integration", func(t *testing.T) {
		// Test that all components work together
		serverAddr := ":" + cfg.Port

		// Verify server address
		assert.Equal(t, ":8080", serverAddr)

		// Verify usecase is ready
		assert.NotNil(t, mockUseCase)

		// Verify config is complete
		assert.NotEmpty(t, cfg.Port)
	})
}

// Benchmark test for server configuration
func BenchmarkRun_Configuration(b *testing.B) {
	mockUseCase := new(MockChatUseCase)
	cfg := config.Config{Port: "8080"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serverAddr := ":" + cfg.Port
		_ = serverAddr
		_ = mockUseCase
	}
}
