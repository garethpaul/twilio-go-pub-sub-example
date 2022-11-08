package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	analytics "github.com/segmentio/analytics-go/v3"
	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageContent struct {
	To  string `json:"to"`
	Message string `json:"message"`
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"messageId"`
	} `json:"message"`
}

func main() {
	r := gin.Default()

	r.POST("/", processMessage)

	r.Run()
}

func processMessage(c *gin.Context) {
	var m PubSubMessage
	var content MessageContent

	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		// print out the error
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(m.Message.Data, &content); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	client := twilio.NewRestClient()
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_FROM_NUMBER"))
	// Lots of Messages Recommend Using a Messaging Service
	// See https://www.twilio.com/docs/messaging/services for details
	// params.SetMessagingServiceSid("MG9752274e9e519418a7406176694466fa")
	params.SetTo(content.To)
	params.SetBody(content.Message)

	// Add Exponential Backoff to retry sending the message twice if there is a failure
	// See https://www.twilio.com/docs/api/errors for details
	// See https://pkg.go.dev/github.com/cenkalti/backoff/v4#section-readme
	// for details on how to implement exponential backoff
	retries := 3
	for retries > 0 {
		resp, err := client.Api.CreateMessage(params)
		// if error is a rate limit error, retry
		// if the error.code is 21610, retry
		// TODO: FIX THIS
		if err != nil {
			twilioError := err.(*twilioClient.TwilioRestError)
			if twilioError.Code == 21610 {
				retries--
				time.Sleep(5 * time.Second)
				continue
			}
		}
		// if there is no error, break out of the loop
		if err == nil {
			// print out the response from resp
			fmt.Println(resp)
			// Log the user in Segment
			// See https://segment.com/docs/connections/sources/catalog/libraries/server/http-api/#identify
			// for details on how to log the user in Segment
			client := analytics.New(os.Getenv("SEGMENT_WRITE_KEY"))
			defer client.Close()

			client.Enqueue(analytics.Track{
				UserId: content.To,
				Event:  "Message Sent",
				Properties: analytics.NewProperties().
					Set("to", content.To).
					Set("message", content.Message),
			})
			break
		}
	}

	if os.Getenv("SUBSCRIPTION_NAME") == "ordinary" {
		c.JSON(http.StatusOK, "ordinary")
		return
	}
}
