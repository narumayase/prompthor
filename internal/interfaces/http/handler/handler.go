package handler

import (
	"net/http"

	"anyprompt/pkg/domain"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles HTTP requests related to chat
type ChatHandler struct {
	chatUseCase domain.ChatUseCase
}

// NewChatHandler creates a new instance of the chat controller
func NewChatHandler(chatUseCase domain.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUseCase,
	}
}

// HandleChat processes the POST chat request
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var request domain.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.chatUseCase.ProcessChat(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing chat: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.ChatResponse{
		Response: response,
	})
}
