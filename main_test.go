package main

import (
	"github.com/stretchr/testify/assert"
	"prompthor/config"
	"prompthor/internal/infrastructure/repository"
	"testing"
)

func TestInitializeRepositories(t *testing.T) {
	t.Run("should return OpenAI repository when configured", func(t *testing.T) {
		cfg := config.Config{
			OpenAIKey: "test-key",
			ChatModel: "OpenAI",
		}
		llmRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.IsType(t, &repository.OpenAIRepository{}, llmRepo)
	})

	t.Run("should return Groq repository when configured", func(t *testing.T) {
		cfg := config.Config{
			GroqAPIKey: "test-key",
		}
		llmRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.IsType(t, &repository.GroqRepository{}, llmRepo)
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
