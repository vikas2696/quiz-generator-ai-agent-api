package main

import (
	"fmt"
	"net/http"
	"os"
	"quiz-generator-ai-agent-api/agent"
	"quiz-generator-ai-agent-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.POST("/query", handleAgentQuery)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // for local dev
	}

	err := server.Run(":" + port)
	if err != nil {
		fmt.Println("error running server: " + err.Error())
	}

}

func handleAgentQuery(context *gin.Context) {

	var user_request models.UserRequest
	context.ShouldBindJSON(&user_request)
	result, err := agent.AgentHandler(user_request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, just try again!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"questions": result})

}
