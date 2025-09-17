package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"prompthor/config"
	"prompthor/internal/domain"
	"testing"
)

// MockHTTPClient is a mock implementation of the HTTPClient for testing purposes.
type MockHTTPClient struct {
	PostFunc func(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error)
}

// Post delegates the call to the PostFunc field.
func (m *MockHTTPClient) Post(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
	return m.PostFunc(ctx, payload)
}

func TestNewGroqRepository(t *testing.T) {
	cfg := config.Config{
		GroqAPIKey: "test_api_key",
		ChatModel:  "test_model",
		GroqUrl:    "http://localhost",
	}
	client := &MockHTTPClient{}
	repo, err := NewGroqRepository(cfg, client)

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	groqRepo, ok := repo.(*GroqRepository)
	assert.True(t, ok)
	assert.Equal(t, "test_api_key", groqRepo.apiKey)
	assert.Equal(t, "test_model", groqRepo.model)
	assert.Equal(t, "http://localhost", groqRepo.baseURL)
	assert.Equal(t, client, groqRepo.httpClient)
}

func TestGroqRepository_Send(t *testing.T) {
	ctx := context.WithValue(context.Background(), "X-Request-Id", "test-request-id")
	prompt := domain.PromptRequest{Prompt: "Hello"}

	t.Run("successful response", func(t *testing.T) {
		mockResponse := GroqResponse{
			Output: []Entry{
				{
					Type: "message",
					Content: []Content{
						{
							Type: "output_text",
							Text: "World",
						},
					},
				},
			},
		}
		mockBody, _ := json.Marshal(mockResponse)
		mockHTTPClient := &MockHTTPClient{
			PostFunc: func(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(mockBody)),
					Status:     "200 OK",
				}, nil
			},
		}

		repo := &GroqRepository{
			apiKey:     "test_api_key",
			model:      "test_model",
			httpClient: mockHTTPClient,
			baseURL:    "http://localhost",
		}

		response, err := repo.Send(ctx, prompt)

		assert.NoError(t, err)
		assert.Equal(t, "World", response)
	})

	t.Run("http client error", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			PostFunc: func(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
				return nil, errors.New("http client error")
			},
		}

		repo := &GroqRepository{
			apiKey:     "test_api_key",
			model:      "test_model",
			httpClient: mockHTTPClient,
			baseURL:    "http://localhost",
		}

		response, err := repo.Send(ctx, prompt)

		assert.Error(t, err)
		assert.Equal(t, "", response)
		assert.Equal(t, "http client error", err.Error())
	})

	t.Run("invalid json response", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			PostFunc: func(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("invalid json"))),
					Status:     "200 OK",
				}, nil
			},
		}

		repo := &GroqRepository{
			apiKey:     "test_api_key",
			model:      "test_model",
			httpClient: mockHTTPClient,
			baseURL:    "http://localhost",
		}

		response, err := repo.Send(ctx, prompt)

		assert.Error(t, err)
		assert.Equal(t, "", response)
	})
}
