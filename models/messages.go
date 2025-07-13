package models

type LLMRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UserRequest struct {
	Topic      string `json:"topic"`
	NoQ        string `json:"noq"`
	Difficulty string `json:"difficulty"`
}
