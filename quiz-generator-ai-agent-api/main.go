package main

import (
	"fmt"
	"net/http"
	"quiz-generator-ai-agent-api/agent"
	"quiz-generator-ai-agent-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.POST("/query", handleAgentQuery)

	err := server.Run("localhost:8080")
	if err != nil {
		fmt.Println("error running server: " + err.Error())
	}

}

func handleAgentQuery(context *gin.Context) {

	var message_from_user models.Message
	context.ShouldBindJSON(&message_from_user)
	result, err := agent.AgentHandler(message_from_user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": result})

}
