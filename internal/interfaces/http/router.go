package http

import (
	"anyprompt/internal/interfaces/http/handler"
	"anyprompt/internal/interfaces/http/middleware"
	"anyprompt/pkg/domain"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(chatUseCase domain.ChatUseCase) *gin.Engine {
	router := gin.Default()

	// Add middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	// Create the controller
	chatHandler := handler.NewChatHandler(chatUseCase)

	// API routes group
	api := router.Group("/api/v1")
	api.POST("/chat/ask", chatHandler.HandleChat)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "AnyPrompt API is running",
		})
	})
	return router
}
