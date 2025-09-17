package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"prompthor/internal/domain"
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

	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	response, err := h.usecase.ProcessChat(ctx, request)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error process chat")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing chat: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
