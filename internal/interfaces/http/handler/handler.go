package handler

import (
	"anyompt/internal/domain"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ChatHandler handles HTTP requests related to chat
type ChatHandler struct {
	usecase domain.ChatUseCase
}

// NewChatHandler creates a new instance of the chat controller
func NewChatHandler(chatUseCase domain.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		usecase: chatUseCase,
	}
}

// HandleChat processes the POST chat request
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var request domain.PromptRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	// TODO ver de hacer un middleware que inyecte los headers automáticamente?
	ctx := context.WithValue(c.Request.Context(), "correlation_id", c.GetHeader("X-Correlation-ID"))
	ctx = context.WithValue(ctx, "routing_id", c.GetHeader("X-Routing-ID"))
	// TODO y otro middleware con el request id.. y el logging... ver

	log.Info().Msgf("X-Correlation-ID received: %v", c.GetHeader("X-Correlation-ID"))
	log.Info().Msgf("X-Routing-ID received: %v", c.GetHeader("X-Routing-ID"))

	response, err := h.usecase.ProcessChat(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing chat: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
