package agent

import (
	"fmt"
	"quiz-generator-ai-agent-api/models"
)

func AgentHandler(user_request models.UserRequest) (any, error) {
	var user_message models.Message
	user_message.Role = "user"
	user_message.Content = "Generate " + user_request.NoQ + " questions about this topic: " + user_request.Topic + " with " + user_request.Difficulty + " difficulty."
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
