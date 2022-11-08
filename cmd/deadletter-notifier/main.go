package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MessageContent struct {
	To      string `json:"name"`
	Message string `json:"bucket"`
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"messageId"`
	} `json:"message"`
}

func main() {
	r := gin.Default()

	r.POST("/", processGsNotification)

	r.Run()
}

func processGsNotification(c *gin.Context) {
	var m PubSubMessage
	var content MessageContent

	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := json.Unmarshal(m.Message.Data, &content); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Sendgrid Go Send Email
	from := mail.NewEmail("PubSub Deadletter", "test@twilio.com")
	subject := "PubSub Deadletter Notifier"
	to := mail.NewEmail("DevTools", "test@twilio.com")
	// plainTextContent WITH CONTENT
	msgContent:= fmt.Sprintf("PubSub Deadletter Notifier: %s, %s", content.To, content.Message)
	htmlContent := fmt.Sprintf("<strong>%s</strong>", msgContent)
	message := mail.NewSingleEmail(from, subject, to, msgContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("E-mail sent: %s", response.Body)
}
