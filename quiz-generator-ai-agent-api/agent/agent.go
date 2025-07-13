package agent

import (
	"encoding/json"
	"fmt"
	"quiz-generator-ai-agent-api/models"
)

func AgentHandler(user_request models.UserRequest) ([]models.Question, error) {

	question_model := ` Each question must be a JSON object with these exact fields:{
		"QuestionId": [integer],
		"Ques": "[clear, specific question]",
		"OptionA": "[first option]",
		"OptionB": "[second option]",
		"OptionC": "[third option]",
		"OptionD": "[fourth option]",
		"Answer": "[complete correct answer text matching one of the options exactly]"
	}
	**Your response should be a JSON array, with perfectly formatted in triple back ticks, of the questions.**
	For example, 
	[{
		"QuestionId": 1,
		"Ques": "What is a qubit?",
		"OptionA": "A bit that can be either 0 or 1.",
		"OptionB": "A bit that can exist in multiple states simultaneously.",
		"OptionC": "A physical wire used to transmit quantum data.",
		"OptionD": "A type of classical computer.",
		"Answer": "A bit that can exist in multiple states simultaneously."
	},
	{
        "QuestionId": 2,
        "Ques": "Where is the Indus Valley Civilization located?",
        "OptionA": "In ancient Egypt",
        "OptionB": "In ancient Mesopotamia",
        "OptionC": "In modern-day Pakistan and India",
        "OptionD": "In South America",
        "Answer": "In modern-day Pakistan and India"
	},.... ]`

	var user_message models.Message
	user_message.Role = "user"
	user_message.Content = "Generate " + user_request.NoQ + " questions about this topic: " + user_request.Topic + " with " + user_request.Difficulty + " difficulty. " + question_model
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
