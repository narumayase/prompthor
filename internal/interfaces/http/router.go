package http

import (
	"anyompt/internal/domain"
	"anyompt/internal/interfaces/http/handler"
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/middleware"
	"github.com/narumayase/anysher/middleware/gateway"
)

// SetupRouter configures the API routes
func SetupRouter(chatUseCase domain.ChatUseCase) *gin.Engine {
	router := gin.Default()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.HeadersToContext())
	router.Use(middleware.RequestIDToLogger())
	router.Use(gateway.Sender())
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
			"message": "anyompt API is running",
		})
	})
	return router
}
