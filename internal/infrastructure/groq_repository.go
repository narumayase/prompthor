package infrastructure

import (
	"anyprompt/internal/config"
	"anyprompt/internal/infrastructure/client"
	"anyprompt/internal/infrastructure/response"
	"anyprompt/pkg/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	url = "https://api.groq.com/openai/v1/responses"
	//TODO hacerla parametrizable!
)

// GroqRepository implements ChatRepository using Groq API
type GroqRepository struct {
	apiKey     string
	model      string
	httpClient client.HTTPClient
	baseURL    string
}

// NewGroqRepository creates a new instance of the Groq repository
func NewGroqRepository(config config.Config) (domain.ChatRepository, error) {
	return &GroqRepository{
		apiKey:     config.GroqAPIKey,
		model:      config.ChatModel,
		httpClient: client.NewDefaultHTTPClient(),
		baseURL:    url,
	}, nil
}

// NewGroqRepositoryWithClient creates a repository with custom HTTP client (for testing)
// TODO ver de cambiar esto
func NewGroqRepositoryWithClient(config config.Config, httpClient client.HTTPClient, baseURL string) (domain.ChatRepository, error) {
	return &GroqRepository{
		apiKey:     config.GroqAPIKey,
		model:      config.ChatModel,
		httpClient: httpClient,
		baseURL:    baseURL,
	}, nil
}

// SendMessage sends a message to ChatGPT and returns the response
func (r *GroqRepository) SendMessage(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": r.model,
		"input": prompt,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", r.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parsear JSON a struct
	var response response.GroqResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	var outputPrompt string
	for _, entry := range response.Output {
		if entry.Type == "message" {
			for _, content := range entry.Content {
				if content.Type == "output_text" {
					outputPrompt = content.Text
					fmt.Println("OutputPrompt:", outputPrompt)
				}
			}
		}
	}
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))

	return outputPrompt, nil
}
