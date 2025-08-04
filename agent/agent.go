package agent

import (
	"errors"
	"fmt"
	"quiz-generator-ai-agent-api/contextfiles"
	"quiz-generator-ai-agent-api/models"
)

func AgentHandler(user_request models.UserRequest) ([]models.Question, error) {

	var user_message models.Message
	user_message.Role = "user"
	user_message.Content = AnalyserPrompt(user_request)

	analyser_context, err := contextfiles.Read_context_file("contextfiles/analyzer_context.json")
	if err != nil {
		return []models.Question{}, errors.New("context file reading error")
	}
	analyser_context = append(analyser_context, user_message)

	//ANALYSER CALL
	result, err := LLMcall(analyser_context, "gemma2-9b-it")
	if err != nil {
		return []models.Question{}, err
	}
	response_message, err := ConvertLLMResult(result)
	if err != nil {
		return []models.Question{}, err
	}
	fmt.Println(response_message.Content)

	questions := []models.Question{}
	var validation_result models.Validation_result
	validation_result.IsValid = false
	validation_result.Remark = ""

	for !validation_result.IsValid { //correction loop

		req_from := "request_from_analyser"
		if validation_result.Remark != "" {
			req_from = "request_from_validator"
		}

		//GENERATOR CALL
		var analyser_message models.Message
		analyser_message.Role = "system"
		analyser_message.Content = GeneratorPrompt(req_from, user_request, response_message.Content, validation_result, questions)

		generator_context, err := contextfiles.Read_context_file("contextfiles/generator_context.json")
		if err != nil {
			return []models.Question{}, errors.New("context file reading error")
		}
		generator_context = append(generator_context, analyser_message)

		gen_result, err := LLMcall(generator_context, "llama3-8b-8192")
		if err != nil {
			return []models.Question{}, err
		}
		gen_response_message, err := ConvertLLMResult(gen_result)
		if err != nil {
			return []models.Question{}, err
		}
		fmt.Println(gen_response_message.Content)

		extracted_content := ExtractJSONBlock(gen_response_message.Content)

		//VALIDATOR CALL
		validation_result, questions = Validator_tool(extracted_content)
		fmt.Println(validation_result.Remark)
	}

	//fmt.Println(response_message.Content)
	return questions, err
}
