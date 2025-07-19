package contextfiles

import (
	"encoding/json"
	"fmt"
	"os"
	"quiz-generator-ai-agent-api/models"
)

func Read_context_file(filename string) ([]models.Message, error) {

	previous_messages := []models.Message{}

	previous_messages_bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error")
		return previous_messages, err
	}

	err = json.Unmarshal(previous_messages_bytes, &previous_messages)
	if err != nil {
		fmt.Println("Unmarshalling error")
		return previous_messages, err
	}

	return previous_messages, err
}

func Write_context_file(filename string, messages []models.Message) error {

	byte_messages, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return err
	}

	err = os.WriteFile(filename, byte_messages, 0644)
	if err != nil {
		fmt.Println("File writing error")
		return err
	}
	return nil
}
