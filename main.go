package main

import (
	"anyompt/cmd/server"
	"anyompt/config"
	"anyompt/internal/application"
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"anyompt/internal/infrastructure/repository"
	"fmt"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repository based on configuration
	chatRepository := initializeRepositories(cfg)

	// Create use case
	chatUseCase := application.NewChatUseCase(chatRepository)

	server.Run(cfg, chatUseCase)
}

// initializeRepositories creates and returns the appropriate chat repository based on configuration
func initializeRepositories(config config.Config) domain.LLMRepository {
	var chatRepository domain.LLMRepository
	switch {
	case config.ChatModel == "OpenAI" && config.OpenAIKey != "":
		// initialize OpenAI repository
		chatRepository = initializeOpenAIRepository(config)
	case config.GroqAPIKey != "":
		// initialize Groq repository
		chatRepository = initializeGroqRepository(config)
	default:
		log.Panic().Err(fmt.Errorf("no valid LLM repository configuration found"))
	}
	return chatRepository
}

// initializeGroqRepository creates and configures a Groq repository instance
func initializeGroqRepository(config config.Config) domain.LLMRepository {
	// Create a new HTTP client
	httpClient := anysherhttp.NewClient(&http.Client{})

	log.Info().Msg("🚀 Starting with Groq API")
	chatRepo, err := repository.NewGroqRepository(config, httpClient)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create Groq repository: %v", err)
		log.Fatal()
	}
	return chatRepo
}

// initializeOpenAIRepository creates and configures an OpenAI repository instance
func initializeOpenAIRepository(config config.Config) domain.LLMRepository {
	openaiClient := client.NewOpenAIClient(config.OpenAIKey)

	log.Info().Msg("🚀 Starting with OpenAI API")
	chatRepo, err := repository.NewOpenAIRepository(openaiClient)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create OpenAI repository: %v", err)
		log.Fatal()
	}
	return chatRepo
}
