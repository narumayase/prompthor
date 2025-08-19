package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"anyprompt/internal/config"
	infra "anyprompt/internal/infrastructure"
	httphandler "anyprompt/internal/interfaces/http"
	"anyprompt/pkg/application"
	"anyprompt/pkg/domain"
)

func Run() {
	// Load configuration
	cfg := config.Load()

	// Create repository based on configuration
	chatRepo := initializeRepositories(*cfg)

	// Create use case
	chatUseCase := application.NewChatUseCase(chatRepo)

	// Configure router
	router := httphandler.SetupRouter(chatUseCase)

	// Start server
	serverAddr := ":" + cfg.Port
	if err := router.Run(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeRepositories(config config.Config) domain.ChatRepository {
	if config.ChatModel == "OpenAI" && config.OpenAIKey != "" {
		fmt.Println("🚀 Starting with OpenAI API")
		chatRepo, err := infra.NewOpenAIRepository()
		if err != nil {
			log.Fatalf("Failed to create OpenAI repository: %v", err)
		}
		return chatRepo
	}
	if config.GroqAPIKey != "" {
		fmt.Println("🚀 Starting with Groq API")
		chatRepo, err := infra.NewGroqRepository(config)
		if err != nil {
			log.Fatalf("Failed to create Groq repository: %v", err)
		}
		return chatRepo
	}
	return nil
}
