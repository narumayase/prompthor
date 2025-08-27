package main

import (
	"anyompt/internal/config"
	"anyompt/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializeRepositories(t *testing.T) {
	t.Run("should return OpenAI repository when configured", func(t *testing.T) {
		cfg := config.Config{
			ChatModel: "OpenAI",
			OpenAIKey: "test-key",
		}
		llmRepo, eventRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.Nil(t, eventRepo)
		assert.IsType(t, &repository.OpenAIRepository{}, llmRepo)
	})

	t.Run("should return Groq repository when configured", func(t *testing.T) {
		cfg := config.Config{
			GroqAPIKey: "test-key",
		}
		llmRepo, eventRepo := initializeRepositories(cfg)
		assert.NotNil(t, llmRepo)
		assert.Nil(t, eventRepo)
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
