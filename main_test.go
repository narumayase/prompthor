package main

import (
	"anyompt/config"
	"anyompt/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializeRepositories(t *testing.T) {
	t.Run("should return OpenAI repository when configured", func(t *testing.T) {
		cfg := config.Config{
			OpenAIKey:      "test-key",
			ChatModel:      "OpenAI",
			GatewayEnabled: true,
		}
		llmRepo, eventRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.NotNil(t, eventRepo)
		assert.IsType(t, &repository.OpenAIRepository{}, llmRepo)
		assert.IsType(t, &repository.AnywayRepository{}, eventRepo)
	})

	t.Run("should return Groq repository when configured", func(t *testing.T) {
		cfg := config.Config{
			GroqAPIKey:     "test-key",
			GatewayEnabled: true,
		}
		llmRepo, eventRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.NotNil(t, eventRepo)
		assert.IsType(t, &repository.GroqRepository{}, llmRepo)
		assert.IsType(t, &repository.AnywayRepository{}, eventRepo)
	})
}

func TestInitializeGroqRepository(t *testing.T) {
	t.Run("should return a new Groq repository", func(t *testing.T) {
		cfg := config.Config{
			GroqAPIKey: "test-key",
		}
		repo := initializeGroqRepository(cfg)
		assert.NotNil(t, repo)
		assert.IsType(t, &repository.GroqRepository{}, repo)
	})
}

func TestInitializeOpenAIRepository(t *testing.T) {
	t.Run("should return a new OpenAI repository", func(t *testing.T) {
		cfg := config.Config{
			OpenAIKey: "test-key",
		}
		repo := initializeOpenAIRepository(cfg)
		assert.NotNil(t, repo)
		assert.IsType(t, &repository.OpenAIRepository{}, repo)
	})
}
