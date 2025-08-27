package repository

import (
	"anyompt/internal/config"
	"anyompt/internal/domain"
	"anyompt/internal/infrastructure/client"
	"anyompt/internal/infrastructure/response"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

const (
	url = "https://api.groq.com/openai/v1/responses"
	// TODO: make this configurable!
)

// GroqRepository implements LLMRepository using Groq API
type GroqRepository struct {
	apiKey     string
	model      string
	httpClient client.HttpClient
	baseURL    string
}

// NewGroqRepository creates a new instance of the Groq repository
func NewGroqRepository(config config.Config, httpClient client.HttpClient) (domain.LLMRepository, error) {
	return &GroqRepository{
		apiKey:     config.GroqAPIKey,
		model:      config.ChatModel,
		httpClient: httpClient,
		baseURL:    config.GroqUrl,
	}, nil
}

// Send sends a message to Groq and returns the response
func (r *GroqRepository) Send(prompt domain.PromptRequest) (string, error) {
	payload := map[string]interface{}{
		"model": r.model,
		"input": prompt.Prompt,
	}
	resp, err := r.httpClient.Post(payload, r.baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON to struct
	var result response.GroqResponse
	if err := json.Unmarshal(body, &result); err != nil {
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
	log.Debug().Msgf("Groq API response: %s", string(body))

	return outputPrompt, nil
}
