package server

import (
	"errors"
	"log"
	"net/http"

	"anyprompt/internal/config"
	httphandler "anyprompt/internal/interfaces/http"
	"anyprompt/pkg/domain"
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
