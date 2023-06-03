package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define the route and attach the handler function
	router.GET("/", handleRequest)

	// Start the server
	router.Run(":8080")
}

func handleRequest(c *gin.Context) {
	var requestData RequestData

	// Parse the request JSON into requestData struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the received request data
	response := fmt.Sprintf("Received Request:\nGoroutine-ID: %s\nRequest-Time: %s\nMessage-Count: %d",
		requestData.GoroutineID, requestData.RequestTime, requestData.MessageCount)

	c.String(http.StatusOK, response)
}
