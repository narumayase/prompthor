package main

import (
	"anyprompt/cmd/server"
	"anyprompt/internal/application"
	"anyprompt/internal/config"
	"anyprompt/internal/domain"
	"anyprompt/internal/infrastructure/client"
	"anyprompt/internal/infrastructure/repository"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repository based on configuration
	chatRepo := initializeRepositories(cfg)

	// Create use case
	chatUseCase := application.NewChatUseCase(chatRepo)

	server.Run(cfg, chatUseCase)
}

// initializeRepositories creates and returns the appropriate chat repository based on configuration
func initializeRepositories(config config.Config) domain.ChatRepository {
	if config.ChatModel == "OpenAI" && config.OpenAIKey != "" {
		return initializeOpenAIRepository(config)
	}
	if config.GroqAPIKey != "" {
		return initializeGroqRepository(config)
	}
	return nil
}

// initializeGroqRepository creates and configures a Groq repository instance
func initializeGroqRepository(config config.Config) domain.ChatRepository {
	groqClient := &http.Client{}
	groqHttpClient := client.NewHttpClient(groqClient, config.GroqAPIKey)

	log.Info().Msg("🚀 Starting with Groq API")
	chatRepo, err := repository.NewGroqRepository(config, groqHttpClient)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create Groq repository: %v", err)
		log.Fatal()
	}
	return chatRepo
}

// initializeOpenAIRepository creates and configures an OpenAI repository instance
func initializeOpenAIRepository(config config.Config) domain.ChatRepository {
	client := client.NewOpenAIClient(config.OpenAIKey)

	log.Info().Msg("🚀 Starting with OpenAI API")
	chatRepo, err := repository.NewOpenAIRepository(client)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create OpenAI repository: %v", err)
		log.Fatal()
	}
	return chatRepo
}
