package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	// add custom user agent
	// See https://github.com/twilio/twilio-go/blob/main/rest/api/v2010/accounts_messages.go
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_FROM_NUMBER"))
	// MILLIONS OF MESSAGES?? 
	// Lots of Messages Recommend Using a Messaging Service
	// See https://www.twilio.com/docs/messaging/services for details
	//params.SetMessagingServiceSid(os.GetEnv("TWILIO_SUBSCRIPTION_SID"))
	//params.ScheduleType("fixed")
	//params.SendAt(time.Now().UTC().Format("2022-12-25T00:04:05-0700"))
	//params.MaxPrice(float32(0.05))
	params.SetTo(content.To)
	params.SetBody(content.Message)

	// TOO MANY REQUESTS ??
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
			// Error Parsing
			// TwilioRestError provides information about an unsuccessful request.
			// https://pkg.go.dev/github.com/twilio/twilio-go/client#TwilioRestError
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

type MyExcellentClient struct {
	twilioClient.Client
}

func (c *MyExcellentClient) SendRequest(method string, rawURL string, data url.Values, headers map[string]interface{}) (*http.Response, error) {
	/* 
		Do something with the request before sending it
		// https://github.com/twilio/twilio-go#using-a-custom-client
	*/
	// Custom code to pre-process request here
	resp, err := c.Client.SendRequest(method, rawURL, data, headers)
	// Custom code to pre-process response here
	fmt.Println(resp.StatusCode)
	return resp, err
}