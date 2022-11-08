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
	MediaLink string `json:"mediaLink"`
	Name      string `json:"name"`
	Bucket    string `json:"bucket"`
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
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail("Example User", "test@example.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := content.MediaLink //"<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
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
