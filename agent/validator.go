package agent

import (
	"encoding/json"
	"fmt"
	"quiz-generator-ai-agent-api/models"
	"strconv"
)

func validate_format(ques_string string) (models.Validation_result, []models.Question) {
	fmt.Println("**********VALIDATING FORMAT**********")
	var validation_result models.Validation_result
	validation_result.IsValid = true
	validation_result.Remark = "Format validation passed."

	var questions []models.Question
	err := json.Unmarshal([]byte(ques_string), &questions)
	if err != nil {
		fmt.Println("Unmarshalling error in validate_format: " + err.Error())
		validation_result.IsValid = false
		validation_result.Remark = fmt.Sprintf(`Format validation failed.
		Error in unmarshalling the json: %s. Fallback and regenerate in the correct format without error.`, err.Error())
		return validation_result, []models.Question{}
	}
	return validation_result, questions
}

func validate_content(questions []models.Question) models.Validation_result {
	fmt.Println("**********VALIDATING CONTENT**********")
	var validation_result models.Validation_result
	validation_result.IsValid = true
	validation_result.Remark = "Content validation passed."

	for i, q := range questions {
		if q.Ques == "" || q.OptionA == "" || q.OptionB == "" || q.OptionC == "" || q.OptionD == "" || q.Answer == "" {
			fmt.Println("Empty field(s) in generated questions")
			validation_result.IsValid = false
			validation_result.Remark = "Content validation failed. Question with QuestionId: " + strconv.Itoa(i+1) + " have an empty field."
			return validation_result
		}
		if !((q.Answer == q.OptionA) || (q.Answer == q.OptionB) || (q.Answer == q.OptionC) || (q.Answer == q.OptionD)) {
			fmt.Println("No option matches the answer")
			validation_result.IsValid = false
			validation_result.Remark = "Content validation failed. The Answer string in Question with QuestionId: " + strconv.Itoa(i+1) + " does not match any option."
			return validation_result
		}
	}

	return validation_result
}

func validate_quality(questions []models.Question) models.Validation_result {
	fmt.Println("**********VALIDATING QUALITY**********")
	var validation_result models.Validation_result
	validation_result.IsValid = true
	validation_result.Remark = "Quality validation passed."

	var generator_message models.Message
	generator_message.Role = "system"
	generator_message.Content = QualityValidatorPrompt(questions)

	result, err := LLMcall([]models.Message{generator_message}, "gemma2-9b-it")
	if err != nil {
		validation_result.Remark = "llm error in quality validation: " + err.Error()
		return validation_result
	}
	validator_response, err := ConvertLLMResult(result)
	if err != nil {
		validation_result.Remark = "decoding llm error in quality validation: " + err.Error()
		return validation_result
	}

	decoded_validator_response := ExtractJSONBlock(validator_response.Content)
	validation_result.Remark = decoded_validator_response
	fmt.Println("Quality ------------- " + decoded_validator_response)
	return validation_result
}

func Validator_tool(ques_string string) (models.Validation_result, []models.Question) {
	var validation_result models.Validation_result

	validation_result, questions := validate_format(ques_string)
	if validation_result.IsValid {
		validation_result = validate_content(questions)
		if validation_result.IsValid {
			validation_result = validate_quality(questions)
			if validation_result.IsValid {
				return validation_result, questions
			}
		}
	}

	return validation_result, questions

}
