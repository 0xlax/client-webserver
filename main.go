package main

import (
	"time"

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
	// Extract required information from the request headers
	goroutineID := c.GetHeader("Goroutine-ID")
	requestTime := c.GetHeader("Request-Time")
	messageCount := c.GetHeader("Message-Count")

	// Echo back the received data in the response
	response := gin.H{
		"Goroutine-ID":  goroutineID,
		"Request-Time":  requestTime,
		"Message-Count": messageCount,
		"Server-Time":   time.Now().UTC(),
	}

	c.JSON(200, response)
}
