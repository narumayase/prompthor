package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "empty environment variable",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			os.Unsetenv(tt.key)

			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	// Clean up environment variables
	envVars := []string{"PORT", "OPENAI_API_KEY", "GROQ_API_KEY", "CHAT_MODEL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Note: We can't easily test Load() function due to godotenv.Load()
	// which would cause log.Fatal. Instead, we test the structure and getEnv function.

	t.Run("config structure", func(t *testing.T) {
		config := Config{
			Port:       getEnv("PORT", "8080"),
			OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
			GroqAPIKey: getEnv("GROQ_API_KEY", ""),
			ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
		}

		assert.Equal(t, "8080", config.Port)
		assert.Empty(t, config.OpenAIKey)
		assert.Empty(t, config.GroqAPIKey)
		assert.Equal(t, "openai/gpt-oss-20b", config.ChatModel)
	})
}

func TestConfig_WithEnvironmentVariables(t *testing.T) {
	// Set test environment variables
	testEnvVars := map[string]string{
		"PORT":           "3000",
		"OPENAI_API_KEY": "test-openai-key",
		"GROQ_API_KEY":   "test-groq-key",
		"CHAT_MODEL":     "test-model",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	t.Run("config with environment variables", func(t *testing.T) {
		config := Config{
			Port:       getEnv("PORT", "8080"),
			OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
			GroqAPIKey: getEnv("GROQ_API_KEY", ""),
			ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
		}

		assert.Equal(t, "3000", config.Port)
		assert.Equal(t, "test-openai-key", config.OpenAIKey)
		assert.Equal(t, "test-groq-key", config.GroqAPIKey)
		assert.Equal(t, "test-model", config.ChatModel)
	})
}

func TestConfig_PartialEnvironmentVariables(t *testing.T) {
	// Clean up all environment variables first
	envVars := []string{"PORT", "OPENAI_API_KEY", "GROQ_API_KEY", "CHAT_MODEL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Set only some environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("OPENAI_API_KEY", "partial-test-key")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("OPENAI_API_KEY")

	t.Run("config with partial environment variables", func(t *testing.T) {
		config := Config{
			Port:       getEnv("PORT", "8080"),
			OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
			GroqAPIKey: getEnv("GROQ_API_KEY", ""),
			ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
		}

		assert.Equal(t, "9000", config.Port)
		assert.Equal(t, "partial-test-key", config.OpenAIKey)
		assert.Empty(t, config.GroqAPIKey)
		assert.Equal(t, "openai/gpt-oss-20b", config.ChatModel)
	})
}

func TestConfig_GroqUrl_Default(t *testing.T) {
	// Test that GroqUrl has the correct default value
	os.Unsetenv("GROQ_URL")

	groqUrl := getEnv("GROQ_URL", "https://api.groq.com/openai/v1/responses")
	assert.Equal(t, "https://api.groq.com/openai/v1/responses", groqUrl)
}

func TestConfig_GroqUrl_Custom(t *testing.T) {
	// Test custom GroqUrl
	customUrl := "https://custom.groq.api.com/v2/responses"
	os.Setenv("GROQ_URL", customUrl)
	defer os.Unsetenv("GROQ_URL")

	groqUrl := getEnv("GROQ_URL", "https://api.groq.com/openai/v1/responses")
	assert.Equal(t, customUrl, groqUrl)
}

func TestLoad_WithEnvFile(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := `PORT=3000
OPENAI_API_KEY=test-openai-key
GROQ_API_KEY=test-groq-key
GROQ_URL=https://test.groq.com/api
CHAT_MODEL=test-model
LOG_LEVEL=debug`

	// Write temporary .env file
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env") // Clean up

	// Test Load function
	config := Load()

	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, "test-openai-key", config.OpenAIKey)
	assert.Equal(t, "test-groq-key", config.GroqAPIKey)
	assert.Equal(t, "https://test.groq.com/api", config.GroqUrl)
	assert.Equal(t, "test-model", config.ChatModel)
}

func TestLoad_WithoutEnvFile(t *testing.T) {
	// Save original working directory
	originalWd, _ := os.Getwd()

	// Create a temporary directory for this test
	tempDir, err := os.MkdirTemp("", "config_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)
	defer os.Chdir(originalWd)

	// Clean environment variables
	envVars := []string{"PORT", "OPENAI_API_KEY", "GROQ_API_KEY", "GROQ_URL", "CHAT_MODEL", "LOG_LEVEL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Test Load function with defaults
	config := Load()

	assert.Equal(t, "8080", config.Port)
	assert.Empty(t, config.OpenAIKey)
	assert.Empty(t, config.GroqAPIKey)
	assert.Equal(t, "https://api.groq.com/openai/v1/responses", config.GroqUrl)
	assert.Equal(t, "openai/gpt-oss-20b", config.ChatModel)
}
