package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

// test deadletter-notifier
//
// go test -v -cover ./cmd/deadletter-notifier/...
// go test -v -cover ./cmd/deadletter-notifier/... -run TestProcessGsNotification

// write test for deadletter-notifier in main.go
func TestProcessGsNotification(t *testing.T) {
	// Create a fake message
	var m PubSubMessage
	m.Message.Data = []byte(`{"to": "+1555555555", "message": "hello-testing"}`)
	m.Message.ID = "1234567890"

	// Convert the message to JSON
	json, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	// Set the request body
	req.Body = ioutil.NopCloser(bytes.NewBuffer(json))

	// Call the function
	processMessage(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

}

