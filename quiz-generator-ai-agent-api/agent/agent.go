package agent

import (
	"fmt"
	"quiz-generator-ai-agent-api/models"
)

func AgentHandler(user_message models.Message) (any, error) {
	messages := []models.Message{user_message}

	result, err := LLMcall(messages)
	if err != nil {
		return models.Message{Role: "system", Content: "something went wrong (LLM)"}, err
	}

	response_message, err := convertLLMResult(result)
	if err != nil {
		return models.Message{Role: "system", Content: "something went wrong (Decoding)"}, err
	}

	fmt.Println(response_message)
	return response_message, err
}
