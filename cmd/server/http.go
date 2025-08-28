package server

import (
	"anyompt/config"
	"errors"
	"log"
	"net/http"

	"anyompt/internal/domain"
	httphandler "anyompt/internal/interfaces/http"
)

func Run(cfg config.Config, usecase domain.ChatUseCase) {
	// Configure router
	router := httphandler.SetupRouter(usecase)

	// Start server
	serverAddr := ":" + cfg.Port
	if err := router.Run(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
