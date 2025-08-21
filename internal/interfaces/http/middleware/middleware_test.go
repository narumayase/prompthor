package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("cors middleware setup", func(t *testing.T) {
		corsMiddleware := CORS()
		assert.NotNil(t, corsMiddleware)
	})

	t.Run("cors headers are set", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("preflight options request", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})
}

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("logger middleware setup", func(t *testing.T) {
		loggerMiddleware := Logger()
		assert.NotNil(t, loggerMiddleware)
	})

	t.Run("logger middleware processes request", func(t *testing.T) {
		router := gin.New()
		router.Use(Logger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("error handler middleware setup", func(t *testing.T) {
		errorMiddleware := ErrorHandler()
		assert.NotNil(t, errorMiddleware)
	})

	t.Run("handles string panic", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic message")
		})

		req, _ := http.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		// Check response contains error structure
		assert.Contains(t, w.Body.String(), "internal_server_error")
		assert.Contains(t, w.Body.String(), "test panic message")
		assert.Contains(t, w.Body.String(), "500")
	})

	t.Run("handles non-string panic", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())
		router.GET("/panic", func(c *gin.Context) {
			panic(123) // Non-string panic
		})

		req, _ := http.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		// Check response contains generic error message
		assert.Contains(t, w.Body.String(), "internal_server_error")
		assert.Contains(t, w.Body.String(), "An unexpected error occurred")
		assert.Contains(t, w.Body.String(), "500")
	})

	t.Run("normal request passes through", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())
		router.GET("/normal", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/normal", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

func TestMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("all middlewares work together", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS())
		router.Use(Logger())
		router.Use(ErrorHandler())
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "all middlewares working"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Body.String(), "all middlewares working")
	})

	t.Run("error handler catches panic with other middlewares", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS())
		router.Use(Logger())
		router.Use(ErrorHandler())
		
		router.GET("/panic", func(c *gin.Context) {
			panic("integration test panic")
		})

		req, _ := http.NewRequest("GET", "/panic", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Body.String(), "integration test panic")
	})
}
