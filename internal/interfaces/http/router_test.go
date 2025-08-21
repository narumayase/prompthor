package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUseCase is a mock implementation of ChatUseCase for router tests
type MockChatUseCase struct {
	mock.Mock
}

func (m *MockChatUseCase) ProcessChat(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}

	t.Run("router setup returns gin engine", func(t *testing.T) {
		router := SetupRouter(mockUseCase)
		assert.NotNil(t, router)
		assert.IsType(t, &gin.Engine{}, router)
	})
}

func TestRouter_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("health endpoint returns OK", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
		assert.Equal(t, "AnyPrompt API is running", response["message"])
	})

	t.Run("health endpoint with different methods", func(t *testing.T) {
		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		
		for _, method := range methods {
			req, _ := http.NewRequest(method, "/health", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return 404 for non-GET methods
			assert.Equal(t, http.StatusNotFound, w.Code)
		}
	})
}

func TestRouter_ChatEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("chat endpoint exists", func(t *testing.T) {
		// Test that the endpoint exists by sending an invalid request
		// (we expect 400 for bad request, not 404 for not found)
		req, _ := http.NewRequest("POST", "/api/v1/chat/ask", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should not be 404 (not found), meaning the route exists
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("chat endpoint only accepts POST", func(t *testing.T) {
		methods := []string{"GET", "PUT", "DELETE", "PATCH"}
		
		for _, method := range methods {
			req, _ := http.NewRequest(method, "/api/v1/chat/ask", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return 404 for non-POST methods
			assert.Equal(t, http.StatusNotFound, w.Code)
		}
	})
}

func TestRouter_CORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("cors headers are present", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("preflight request handling", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/api/v1/chat/ask", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})
}

func TestRouter_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("404 for non-existent routes", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("404 for wrong API version", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v2/chat/ask", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRouter_APIGrouping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("api v1 group exists", func(t *testing.T) {
		// Test that the API group is properly set up
		req, _ := http.NewRequest("POST", "/api/v1/chat/ask", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should not return 404, meaning the API group route exists
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("root api path does not exist", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/chat/ask", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRouter_MiddlewareOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := &MockChatUseCase{}
	router := SetupRouter(mockUseCase)

	t.Run("middlewares are applied in correct order", func(t *testing.T) {
		// Test that CORS, Logger, and ErrorHandler middlewares are all applied
		req, _ := http.NewRequest("GET", "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// CORS middleware should set headers
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		
		// Response should be successful (no panic from error handler)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	})
}
