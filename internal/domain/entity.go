package domain

// PromptRequest represents the chat request
type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// ChatResponse represents the chat response
type ChatResponse struct {
	MessageResponse string `json:"message_response"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
