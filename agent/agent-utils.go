package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"quiz-generator-ai-agent-api/models"
	"strings"

	"github.com/joho/godotenv"
)

func ExtractJSONBlock(input string) string {
	start := strings.Index(input, "[")
	end := strings.LastIndex(input, "]")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return input[start : end+1]
}

func LLMcall(messages []models.Message, model string) (map[string]any, error) {

	llm_endpoint_url := "https://api.groq.com/openai/v1/chat/completions"
	var result map[string]any

	llm_request_body := models.LLMRequestBody{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	json_request_body, err := json.Marshal(llm_request_body)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return result, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", llm_endpoint_url, bytes.NewReader(json_request_body))
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return result, err
	}

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found!")
	}
	api_key := os.Getenv("LLM_API_KEY")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_key)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return result, err
	}
	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&result)
	return result, err
}

func ConvertLLMResult(result map[string]any) (models.Message, error) {
	var received_message models.Message

	if response, ok := result["choices"].([]any); ok && len(response) > 0 {
		full_message := response[0].(map[string]any)
		message := full_message["message"].(map[string]any)
		received_message.Role = message["role"].(string)
		received_message.Content = message["content"].(string)
		return received_message, nil
	} else {
		fmt.Println("ERROR.............................")
		error := result["error"].(map[string]any)
		fmt.Println(error)
		received_message.Role = "assistant"
		received_message.Content = "Some error occured"
		return received_message, errors.New(error["code"].(string))
	}

}
