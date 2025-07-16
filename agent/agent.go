package agent

import (
	"errors"
	"fmt"
	"quiz-generator-ai-agent-api/models"
	"quiz-generator-ai-agent-api/tools"
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

	validation_result, questions := tools.Validator_tool(extracted_content)

	if !validation_result.IsValid {
		fmt.Println(response_message.Content)
		fmt.Println(validation_result.Remark)
		return questions, errors.New("not valid questions")
	}

	fmt.Println(response_message.Content)
	return questions, err
}
