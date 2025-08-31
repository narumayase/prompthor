package domain

// PromptRequest represents the chat request
type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// ChatResponse represents the chat response
type ChatResponse struct {
	Response string `json:"response"`
}
