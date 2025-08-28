package main

import (
	"anyompt/cmd/server"
	"anyompt/config"
	"anyompt/internal/application"
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"anyompt/internal/infrastructure/repository"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repository based on configuration
	chatRepo, eventRepo := initializeRepositories(cfg)
	if eventRepo != nil {
		defer eventRepo.Close()
	}

	// Create use case
	chatUseCase := application.NewChatUseCase(chatRepo, eventRepo)

	server.Run(cfg, chatUseCase)
}

// initializeRepositories creates and returns the appropriate chat repository based on configuration
func initializeRepositories(config config.Config) (domain.LLMRepository, domain.ProducerRepository) {
	var chatRepo domain.LLMRepository
	if config.ChatModel == "OpenAI" && config.OpenAIKey != "" {
		chatRepo = initializeOpenAIRepository(config)
	} else if config.GroqAPIKey != "" {
		chatRepo = initializeGroqRepository(config)
	}

	eventRepo, err := repository.NewKafkaRepository(config)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create Kafka repository: %v", err)
	}

	return chatRepo, eventRepo
}

// initializeGroqRepository creates and configures a Groq repository instance
func initializeGroqRepository(config config.Config) domain.LLMRepository {
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
func initializeOpenAIRepository(config config.Config) domain.LLMRepository {
	client := client.NewOpenAIClient(config.OpenAIKey)

	log.Info().Msg("🚀 Starting with OpenAI API")
	chatRepo, err := repository.NewOpenAIRepository(client)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create OpenAI repository: %v", err)
		log.Fatal()
	}
	return chatRepo
}
