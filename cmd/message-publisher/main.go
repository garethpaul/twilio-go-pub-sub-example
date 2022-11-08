package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
)

type MessageContent struct {
	ToNumber   string `json:"to"`
	Message    string `json:"message"`
}

func main() {
	r := gin.Default()

	r.POST("/new-review", publishReview)

	r.Run()
}

func publishReview(c *gin.Context) {
	var wg sync.WaitGroup

	projectId := os.Getenv("PROJECT_ID")
	topic := os.Getenv("TOPIC")

	client, err := pubsub.NewClient(c, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	defer client.Close()

	t := client.Topic(topic)

	//Skipping error handling here :-))
	var content MessageContent
	// get the message content from the request
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	msgJson, _ := json.Marshal(content)

	result := t.Publish(c, &pubsub.Message{
		Data: msgJson,
	})

	wg.Add(1)

	go func(res *pubsub.PublishResult) {
		defer wg.Done()

		_, err := res.Get(c)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Failed to publish: %v", err)
			return
		}
	}(result)
	
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"message": "Messages generated successfully",
	})
}