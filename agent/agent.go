package agent

import (
	"encoding/json"
	"fmt"
	"quiz-generator-ai-agent-api/models"
)

func AgentHandler(user_request models.UserRequest) ([]models.Question, error) {

	var user_message models.Message
	user_message.Role = "user"
	user_message.Content = getPrompt(user_request)
	messages := []models.Message{user_message}

	result, err := LLMcall(messages)
	if err != nil {
		return []models.Question{}, err
	}

	response_message, err := convertLLMResult(result)
	if err != nil {
		return []models.Question{}, err
	}

	extracted_content := ExtractJSONBlock(response_message.Content)

	var questions []models.Question
	err = json.Unmarshal([]byte(extracted_content), &questions)
	if err != nil {
		fmt.Println("Unmarshalling error")
		return []models.Question{}, err
	}

	fmt.Println(response_message.Content)
	return questions, err
}
