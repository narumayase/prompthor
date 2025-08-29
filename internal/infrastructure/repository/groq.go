package repository

import (
	"anyompt/config"
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/response"
	"anyompt/internal/interfaces/http"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

// GroqRepository implements LLMRepository using Groq API
type GroqRepository struct {
	apiKey     string
	model      string
	httpClient http.HTTPClient
	baseURL    string
}

// NewGroqRepository creates a new instance of the Groq repository
func NewGroqRepository(config config.Config, httpClient http.HTTPClient) (domain.LLMRepository, error) {
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
		log.Error().Err(err).Msg("failed to marshal payload")
		return "", err
	}
	// send to Groq
	resp, err := r.httpClient.Post(r.baseURL, "application/json", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON to struct
	var result response.GroqResponse
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
					log.Debug().Msgf("output prompt: %s", outputPrompt)
				}
			}
		}
	}
	log.Info().Msgf("Groq API response status: %s", resp.Status)
	log.Debug().Msgf("Groq API response: %s", string(respBody))

	return outputPrompt, nil
}
