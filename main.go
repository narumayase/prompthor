package main

import (
	"anyompt/cmd/server"
	"anyompt/config"
	"anyompt/internal/application"
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"anyompt/internal/infrastructure/repository"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repository based on configuration
	chatRepository, producerRepository := initializeRepositories(cfg)

	// Create use case
	chatUseCase := application.NewChatUseCase(chatRepository, producerRepository)

	server.Run(cfg, chatUseCase)
}

// initializeRepositories creates and returns the appropriate chat repository based on configuration
func initializeRepositories(config config.Config) (domain.LLMRepository, domain.ProducerRepository) {
	var chatRepository domain.LLMRepository
	if config.ChatModel == "OpenAI" && config.OpenAIKey != "" {
		chatRepository = initializeOpenAIRepository(config)
	} else if config.GroqAPIKey != "" {
		chatRepository = initializeGroqRepository(config)
	}
	producerRepository := initializeProducerRepository(config)
	return chatRepository, producerRepository
}

// initializeGroqRepository creates and configures a Groq repository instance
func initializeGroqRepository(config config.Config) domain.LLMRepository {
	cfg := anysherhttp.NewConfiguration(config.LogLevel)

	// Create a new HTTP client
	httpClient := anysherhttp.NewClient(&http.Client{}, cfg)

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

// initializeProducerRepository creates and configures a producer repository instance
func initializeProducerRepository(config config.Config) domain.ProducerRepository {
	if config.GatewayEnabled {
		cfg := anysherhttp.NewConfiguration(config.LogLevel)

		// Create a new HTTP client
		httpClient := anysherhttp.NewClient(&http.Client{}, cfg)

		return repository.NewAnywayRepository(config, httpClient)

	}
	return nil
}
