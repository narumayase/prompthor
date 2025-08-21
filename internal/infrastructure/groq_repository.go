package infrastructure

import (
	"anyprompt/internal/config"
	"anyprompt/internal/infrastructure/response"
	"anyprompt/pkg/domain"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

const (
	url = "https://api.groq.com/openai/v1/responses"
	// TODO: make this configurable!
)

// GroqRepository implements ChatRepository using Groq API
type GroqRepository struct {
	apiKey     string
	model      string
	httpClient HttpClient
	baseURL    string
}

// NewGroqRepository creates a new instance of the Groq repository
func NewGroqRepository(config config.Config, httpClient HttpClient) (domain.ChatRepository, error) {
	return &GroqRepository{
		apiKey:     config.GroqAPIKey,
		model:      config.ChatModel,
		httpClient: httpClient,
		baseURL:    config.GroqUrl,
	}, nil
}

// SendMessage sends a message to ChatGPT and returns the response
func (r *GroqRepository) SendMessage(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": r.model,
		"input": prompt,
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
