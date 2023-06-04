package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	GoroutineID  string `json:"gid"`
	RequestTime  string `json:"rqts"`
	ResponseTime string `json:"rpts"`
	MessageCount int    `json:"mc"`
}

type RateLimiter struct {
	mu              sync.Mutex
	lastAccessTimes map[string]time.Time
}

func (rl *RateLimiter) AllowRequest(goroutineID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	currentTime := time.Now()
	lastAccessTime, lastAccessExists := rl.lastAccessTimes[goroutineID]

	if lastAccessExists && currentTime.Sub(lastAccessTime) < time.Second {
		return false
	}

	rl.lastAccessTimes[goroutineID] = currentTime
	return true
}

func main() {
	router := gin.Default()

	// Define the route and attach the handler function
	router.POST("/", handleRequest)

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

	// Echo back the received data and append server response time
	responseData := RequestData{
		GoroutineID:  requestData.GoroutineID,
		RequestTime:  requestData.RequestTime,
		ResponseTime: time.Now().UTC().Format(time.RFC3339),
		MessageCount: requestData.MessageCount,
	}

	c.JSON(http.StatusOK, responseData)
}

func startClients(numClients int) {
	var wg sync.WaitGroup
	client := http.DefaultClient

	limiter := RateLimiter{
		lastAccessTimes: make(map[string]time.Time),
	}

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(goroutineID string) {
			defer wg.Done()
			for {
				requestData := RequestData{
					GoroutineID:  goroutineID,
					RequestTime:  time.Now().UTC().Format(time.RFC3339),
					MessageCount: 0, // Set the initial message count to 0
				}
				sendRequest(client, &limiter, requestData)
				requestData.MessageCount++
				time.Sleep(time.Second) // Control the timing between requests
			}
		}(fmt.Sprintf("goroutine-%d", i+1))
	}

	wg.Wait()
}

func sendRequest(client *http.Client, limiter *RateLimiter, requestData RequestData) {
	if !limiter.AllowRequest(requestData.GoroutineID) {
		fmt.Printf("Request blocked for goroutine '%s'. Sleeping for 1 minute.\n", requestData.GoroutineID)
		time.Sleep(time.Minute)
		return
	}

	// Convert request data to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Send POST request to the server
	resp, err := client.Post("http://localhost:8080", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(body))
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
