package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define the route and attach the handler function
	router.GET("/", handleRequest)

	// Start the server
	go func() {
		if err := router.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	// Start the clients
	numClients := 3
	startClients(numClients)
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

func startClients(numClients int) {
	var wg sync.WaitGroup
	client := http.DefaultClient

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go sendRequest(client, &wg, RequestData{
			GoroutineID:  fmt.Sprintf("Goroutine-%d", i),
			RequestTime:  "2023-06-03T12:00:00Z",
			MessageCount: strconv.Itoa(i + 1),
		})
	}

	wg.Wait()
}
