package repository

import (
	"anyompt/config"
	"anyompt/internal/domain"
	"context"
	"encoding/json"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

// HTTPClient is an interface for an HTTP client.
type HTTPClient interface {
	Post(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error)
}

// GroqResponse is the response from the Groq API
type GroqResponse struct {
	ID     string  `json:"id"`
	Output []Entry `json:"output"`
}

// Entry is a single entry in the Groq response
type Entry struct {
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Content []Content `json:"content,omitempty"`
	Summary []string  `json:"summary,omitempty"`
}

// Content is the content of a Groq response entry
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// GroqRepository implements LLMRepository using Groq API
type GroqRepository struct {
	apiKey     string
	model      string
	httpClient HTTPClient
	baseURL    string
}

// NewGroqRepository creates a new instance of the Groq repository
func NewGroqRepository(config config.Config, httpClient HTTPClient) (domain.LLMRepository, error) {
	return &GroqRepository{
		apiKey:     config.GroqAPIKey,
		model:      config.ChatModel,
		httpClient: httpClient,
		baseURL:    config.GroqUrl,
	}, nil
}

// Send sends a message to Groq and returns the response
func (r *GroqRepository) Send(ctx context.Context, prompt domain.PromptRequest) (string, error) {
	payload := map[string]interface{}{
		"model": r.model,
		"input": prompt.Prompt,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to marshal payload")
		return "", err
	}
	// send to anyway
	resp, err := r.httpClient.Post(ctx, anysherhttp.Payload{
		URL:   r.baseURL,
		Token: r.apiKey,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-Id": ctx.Value("X-Request-Id").(string),
		},
		Content: body,
	})

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Parse JSON to struct
	var result GroqResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	var outputPrompt string
	for _, entry := range result.Output {
		// TODO ver de hacer más bonito esto
		if entry.Type == "message" {
			for _, content := range entry.Content {
				if content.Type == "output_text" {
					outputPrompt = content.Text
					log.Ctx(ctx).Debug().Msgf("output prompt: %s", outputPrompt)
				}
			}
		}
	}
	log.Ctx(ctx).Info().Msgf("Groq API response status: %s", resp.Status)
	log.Ctx(ctx).Debug().Msgf("Groq API response: %s", string(respBody))

	return outputPrompt, nil
}
