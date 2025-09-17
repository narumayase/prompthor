package server

import (
	"errors"
	"log"
	"net/http"
	"prompthor/config"

	"prompthor/internal/domain"
	httphandler "prompthor/internal/interfaces/http"
)

func Run(config config.Config, usecase domain.ChatUseCase) {
	// Configure router
	router := httphandler.SetupRouter(usecase)

	// Start server
	serverAddr := ":" + config.Port
	if err := router.Run(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
