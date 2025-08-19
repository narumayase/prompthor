package response

type GroqResponse struct {
	ID     string  `json:"id"`
	Output []Entry `json:"output"`
}

type Entry struct {
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Content []Content `json:"content,omitempty"`
	Summary []string  `json:"summary,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
